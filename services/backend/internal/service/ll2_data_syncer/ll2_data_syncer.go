package ll2datasyncer

import (
	"errors"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/vamosdalian/launchdate-backend/internal/config"
	"github.com/vamosdalian/launchdate-backend/internal/models"
	"github.com/vamosdalian/launchdate-backend/internal/service/core"
	"github.com/vamosdalian/launchdate-backend/internal/service/ll2"
	"github.com/vamosdalian/launchdate-backend/internal/util"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	SyncTypeLaunch         = "launch"
	SyncTypeAgency         = "agency"
	SyncTypeLauncher       = "launcher"
	SyncTypeLauncherFamily = "launcher_family"
	SyncTypePad            = "pad"
	SyncTypeLocation       = "location"
	SyncTypeUpdate         = "update"

	interruptedTaskMessage = "task interrupted by service restart"
	unexpectedStopMessage  = "task stopped unexpectedly"
)

type LL2DataSyncer struct {
	ll2Service      *ll2.LL2Service
	core            *core.MainService
	rl              util.RateLimiter
	runMu           *sync.Mutex
	stateMu         sync.Mutex
	currentSyncer   Syncer
	currentTaskType string
}

func NewLL2DataSyncer(cfg *config.Config, ll2Service *ll2.LL2Service, core *core.MainService) *LL2DataSyncer {
	return &LL2DataSyncer{
		ll2Service: ll2Service,
		core:       core,
		rl:         util.NewRateLimit(time.Duration(cfg.LL2RequestInterval) * time.Second),
		runMu:      &sync.Mutex{},
	}
}

func (ds *LL2DataSyncer) InitSync(syncType string) (*TaskInfo, error) {
	if !isValidSyncType(syncType) {
		return nil, errors.New("unknown sync type")
	}
	if !ds.runMu.TryLock() {
		return nil, errors.New("another syncer is already running")
	}

	now := time.Now().UTC()
	task, err := ds.loadOrCreateTask(syncType)
	if err != nil {
		ds.runMu.Unlock()
		return nil, err
	}

	ds.prepareTaskForStart(task, now)
	if err := ds.saveTask(task); err != nil {
		ds.runMu.Unlock()
		return nil, err
	}

	syncer, err := ds.newSyncer(syncType)
	if err != nil {
		ds.runMu.Unlock()
		return nil, err
	}

	ds.stateMu.Lock()
	ds.currentSyncer = syncer
	ds.currentTaskType = syncType
	ds.stateMu.Unlock()

	syncer.Start()
	go ds.waitForSyncDone(syncType, syncer)

	return taskInfoFromSyncTask(task), nil
}

func (ds *LL2DataSyncer) PauseSync() (*TaskInfo, error) {
	task, err := ds.getActiveTask()
	if err != nil {
		return nil, err
	}
	if task.Status != models.SyncTaskStatusRunning {
		return nil, errors.New("task is not running")
	}

	syncer, currentType := ds.snapshotCurrentSyncer()
	if syncer == nil {
		if task.Type != SyncTypeUpdate {
			return nil, errors.New("no syncer is running")
		}
	} else if currentType != task.Type {
		return nil, errors.New("another task is active")
	} else {
		syncer.Pause()
	}

	task.Status = models.SyncTaskStatusPaused
	task.LastError = ""
	task.FinishedAt = time.Time{}
	if err := ds.saveTask(task); err != nil {
		return nil, err
	}

	return taskInfoFromSyncTask(task), nil
}

func (ds *LL2DataSyncer) ResumeSync() (*TaskInfo, error) {
	task, err := ds.getActiveTask()
	if err != nil {
		return nil, err
	}
	if task.Status != models.SyncTaskStatusPaused {
		return nil, errors.New("task is not paused")
	}

	syncer, currentType := ds.snapshotCurrentSyncer()
	if syncer != nil {
		if currentType != task.Type {
			return nil, errors.New("another task is active")
		}
		syncer.Resume()
	} else {
		if task.Type != SyncTypeUpdate {
			return nil, errors.New("no syncer is running")
		}
		if !ds.runMu.TryLock() {
			return nil, errors.New("another syncer is already running")
		}

		restoredSyncer, newErr := ds.newSyncer(task.Type)
		if newErr != nil {
			ds.runMu.Unlock()
			return nil, newErr
		}

		ds.stateMu.Lock()
		ds.currentSyncer = restoredSyncer
		ds.currentTaskType = task.Type
		ds.stateMu.Unlock()

		restoredSyncer.Start()
		go ds.waitForSyncDone(task.Type, restoredSyncer)
	}

	now := time.Now().UTC()
	task.Status = models.SyncTaskStatusRunning
	task.LastError = ""
	task.FinishedAt = time.Time{}
	if task.Type == SyncTypeUpdate {
		task.NextRunAt = now
	}
	if err := ds.saveTask(task); err != nil {
		return nil, err
	}

	return taskInfoFromSyncTask(task), nil
}

func (ds *LL2DataSyncer) CancelSync() (*TaskInfo, error) {
	task, err := ds.getActiveTask()
	if err != nil {
		return nil, err
	}
	if !task.Status.IsActive() {
		return nil, errors.New("task can only be canceled from running or paused state")
	}

	syncer, currentType := ds.snapshotCurrentSyncer()
	if syncer != nil {
		if currentType != task.Type {
			return nil, errors.New("another task is active")
		}
		syncer.Cancel()
	} else if task.Type != SyncTypeUpdate {
		return nil, errors.New("no syncer is running")
	}

	now := time.Now().UTC()
	task.Status = models.SyncTaskStatusCanceled
	task.LastError = ""
	task.FinishedAt = now
	if task.Type == SyncTypeUpdate {
		task.NextRunAt = time.Time{}
		task.CurrentWindowStart = time.Time{}
		task.CurrentWindowEnd = time.Time{}
	}
	if err := ds.saveTask(task); err != nil {
		return nil, err
	}

	return taskInfoFromSyncTask(task), nil
}

func (ds *LL2DataSyncer) GetCurrentTask() *TaskInfo {
	task, err := ds.ll2Service.GetActiveSyncTask()
	if err == nil {
		return taskInfoFromSyncTask(task)
	}
	if err != mongo.ErrNoDocuments {
		logrus.Errorf("get active sync task failed: %v", err)
	}

	return nil
}

func (ds *LL2DataSyncer) GetTaskHistory(limit int) ([]*TaskInfo, error) {
	tasks, err := ds.ll2Service.ListRecentSyncTasks(limit)
	if err != nil {
		return nil, err
	}

	history := make([]*TaskInfo, 0, len(tasks))
	for _, task := range tasks {
		history = append(history, taskInfoFromSyncTask(task))
	}

	return history, nil
}

func (ds *LL2DataSyncer) RestoreTasks() error {
	tasks, err := ds.ll2Service.ListSyncTasksByStatuses(models.SyncTaskStatusRunning, models.SyncTaskStatusPaused)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	var updateTask *models.SyncTask
	for _, task := range tasks {
		if task.Type == SyncTypeUpdate {
			if task.Status == models.SyncTaskStatusRunning && updateTask == nil {
				updateTask = task
			}
			continue
		}

		task.Status = models.SyncTaskStatusFailed
		task.LastError = interruptedTaskMessage
		task.FinishedAt = now
		if err := ds.saveTask(task); err != nil {
			return err
		}
	}

	if updateTask == nil {
		return nil
	}
	if !ds.runMu.TryLock() {
		return errors.New("another syncer is already running")
	}

	syncer, err := ds.newSyncer(SyncTypeUpdate)
	if err != nil {
		ds.runMu.Unlock()
		return err
	}

	ds.stateMu.Lock()
	ds.currentSyncer = syncer
	ds.currentTaskType = SyncTypeUpdate
	ds.stateMu.Unlock()

	syncer.Start()
	go ds.waitForSyncDone(SyncTypeUpdate, syncer)
	return nil
}

func (ds *LL2DataSyncer) waitForSyncDone(syncType string, syncer Syncer) {
	<-syncer.Done()

	ds.stateMu.Lock()
	if ds.currentSyncer == syncer {
		ds.currentSyncer = nil
		ds.currentTaskType = ""
	}
	ds.stateMu.Unlock()

	task, err := ds.loadOrCreateTask(syncType)
	if err != nil {
		logrus.Errorf("load sync task %s after completion failed: %v", syncType, err)
		ds.runMu.Unlock()
		return
	}

	now := time.Now().UTC()
	switch task.Status {
	case models.SyncTaskStatusRunning:
		if syncType == SyncTypeUpdate {
			task.Status = models.SyncTaskStatusFailed
			if task.LastError == "" {
				task.LastError = unexpectedStopMessage
			}
		} else {
			task.Status = models.SyncTaskStatusCompleted
			task.LastError = ""
		}
		task.FinishedAt = now
	case models.SyncTaskStatusCanceled, models.SyncTaskStatusFailed, models.SyncTaskStatusCompleted:
		if task.FinishedAt.IsZero() {
			task.FinishedAt = now
		}
	}

	if err := ds.saveTask(task); err != nil {
		logrus.Errorf("save sync task %s after completion failed: %v", syncType, err)
	}

	ds.runMu.Unlock()
}

func (ds *LL2DataSyncer) newSyncer(syncType string) (Syncer, error) {
	switch syncType {
	case SyncTypeLaunch:
		return NewLaunchSyncer(ds.rl, ds.ll2Service, ds.core, ds.reportCurrentTaskProgress), nil
	case SyncTypeAgency:
		return NewAgencySyncer(ds.rl, ds.ll2Service, ds.core, ds.reportCurrentTaskProgress), nil
	case SyncTypeLauncher:
		return NewLauncherSyncer(ds.rl, ds.ll2Service, ds.core, ds.reportCurrentTaskProgress), nil
	case SyncTypeLauncherFamily:
		return NewLauncherFamilySyncer(ds.rl, ds.ll2Service, ds.reportCurrentTaskProgress), nil
	case SyncTypePad:
		return NewPadSyncer(ds.rl, ds.ll2Service, ds.reportCurrentTaskProgress), nil
	case SyncTypeLocation:
		return NewLocationSyncer(ds.rl, ds.ll2Service, ds.core, ds.reportCurrentTaskProgress), nil
	case SyncTypeUpdate:
		return NewUpdateSyncer(ds.rl, ds.ll2Service, ds.core), nil
	default:
		return nil, errors.New("unknown sync type")
	}
}

func (ds *LL2DataSyncer) prepareTaskForStart(task *models.SyncTask, now time.Time) {
	task.Type = task.ID
	task.Status = models.SyncTaskStatusRunning
	task.StartedAt = now
	task.FinishedAt = time.Time{}
	task.LastRun = now
	task.LastError = ""
	task.CurrentOffset = 0
	task.CurrentTotal = 0

	if task.Type == SyncTypeUpdate {
		if task.OverlapSeconds == 0 {
			task.OverlapSeconds = int(defaultUpdateOverlap / time.Second)
		}
		task.NextRunAt = now
		task.CurrentWindowStart = time.Time{}
		task.CurrentWindowEnd = time.Time{}
		task.Progress = buildSyncProgress(task)
		return
	}

	task.NextRunAt = time.Time{}
	task.Progress = buildCountProgress(0, 0)
}

func (ds *LL2DataSyncer) saveTask(task *models.SyncTask) error {
	if task.Type == SyncTypeUpdate {
		task.Progress = buildSyncProgress(task)
	}
	return ds.ll2Service.UpsertSyncTask(task)
}

func (ds *LL2DataSyncer) reportCurrentTaskProgress(progress map[string]interface{}) {
	ds.stateMu.Lock()
	syncType := ds.currentTaskType
	ds.stateMu.Unlock()
	if syncType == "" {
		return
	}

	task, err := ds.loadOrCreateTask(syncType)
	if err != nil {
		logrus.Errorf("load sync task %s for progress update failed: %v", syncType, err)
		return
	}

	progressCopy := make(map[string]interface{}, len(progress))
	for key, value := range progress {
		progressCopy[key] = value
	}
	task.Progress = progressCopy
	task.LastError = ""
	if current, ok := intValue(progressCopy, TaskProgressCurrentCount); ok {
		task.CurrentOffset = current
	}
	if total, ok := intValue(progressCopy, TaskProgressTotalCount); ok {
		task.CurrentTotal = total
	}
	if err := ds.saveTask(task); err != nil {
		logrus.Errorf("persist sync task %s progress failed: %v", syncType, err)
	}
}

func (ds *LL2DataSyncer) loadOrCreateTask(syncType string) (*models.SyncTask, error) {
	task, err := ds.ll2Service.GetSyncTask(syncType)
	if err == nil {
		if task.Type == "" {
			task.Type = syncType
		}
		if syncType == SyncTypeUpdate && task.OverlapSeconds == 0 {
			task.OverlapSeconds = int(defaultUpdateOverlap / time.Second)
		}
		return task, nil
	}
	if err != mongo.ErrNoDocuments {
		return nil, err
	}

	task = &models.SyncTask{
		ID:     syncType,
		Type:   syncType,
		Status: models.SyncTaskStatusIdle,
	}
	if syncType == SyncTypeUpdate {
		task.OverlapSeconds = int(defaultUpdateOverlap / time.Second)
	}
	return task, nil
}

func (ds *LL2DataSyncer) getActiveTask() (*models.SyncTask, error) {
	task, err := ds.ll2Service.GetActiveSyncTask()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("no syncer is running")
		}
		return nil, err
	}
	return task, nil
}

func (ds *LL2DataSyncer) snapshotCurrentSyncer() (Syncer, string) {
	ds.stateMu.Lock()
	defer ds.stateMu.Unlock()
	return ds.currentSyncer, ds.currentTaskType
}
