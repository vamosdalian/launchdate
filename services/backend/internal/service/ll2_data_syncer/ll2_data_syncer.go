package ll2datasyncer

import (
	"errors"
	"sync"
	"time"

	"github.com/vamosdalian/launchdate-backend/internal/config"
	"github.com/vamosdalian/launchdate-backend/internal/models"
	"github.com/vamosdalian/launchdate-backend/internal/service/core"
	"github.com/vamosdalian/launchdate-backend/internal/service/ll2"
	"github.com/vamosdalian/launchdate-backend/internal/util"
)

const (
	SyncTypeLaunch         = "launch"
	SyncTypeAgency         = "agency"
	SyncTypeLauncher       = "launcher"
	SyncTypeLauncherFamily = "launcher_family"
	SyncTypePad            = "pad"
	SyncTypeLocation       = "location"
	SyncTypeUpcoming       = "upcoming"
)

type LL2DataSyncer struct {
	ll2Service    *ll2.LL2Service
	core          *core.MainService
	rl            util.RateLimiter
	runMu         *sync.Mutex
	stateMu       sync.Mutex
	currentSyncer Syncer
	currentTask   *TaskInfo
}

func NewLL2DataSyncer(cfg *config.Config, ll2Service *ll2.LL2Service, core *core.MainService) *LL2DataSyncer {
	return &LL2DataSyncer{
		ll2Service: ll2Service,
		core:       core,
		rl:         util.NewRateLimit(time.Duration(cfg.LL2RequestInterval) * time.Second),
		runMu:      &sync.Mutex{},
	}
}

func (ds *LL2DataSyncer) InitSync(syncType string) error {
	if !isValidSyncType(syncType) {
		return errors.New("unknown sync type")
	}

	if !ds.runMu.TryLock() {
		return errors.New("another syncer is already running")
	}

	now := time.Now()
	ds.stateMu.Lock()
	defer ds.stateMu.Unlock()

	switch syncType {
	case SyncTypeLaunch:
		ds.currentSyncer = NewLaunchSyncer(ds.rl, ds.ll2Service, ds.core)
	case SyncTypeAgency:
		ds.currentSyncer = NewAgencySyncer(ds.rl, ds.ll2Service, ds.core)
	case SyncTypeLauncher:
		ds.currentSyncer = NewLauncherSyncer(ds.rl, ds.ll2Service, ds.core)
	case SyncTypeLauncherFamily:
		ds.currentSyncer = NewLauncherFamilySyncer(ds.rl, ds.ll2Service)
	case SyncTypePad:
		ds.currentSyncer = NewPadSyncer(ds.rl, ds.ll2Service)
	case SyncTypeLocation:
		ds.currentSyncer = NewLocationSyncer(ds.rl, ds.ll2Service, ds.core)
	case SyncTypeUpcoming:
		ds.currentSyncer = NewUpcomingSyncer(ds.rl, ds.ll2Service, ds.core)
	}
	ds.currentTask = &TaskInfo{
		Type:      syncType,
		Status:    models.SyncTaskStatusRunning,
		StartedAt: now,
		UpdatedAt: now,
	}

	ds.currentSyncer.Start()

	go func() {
		<-ds.currentSyncer.Done()
		ds.stateMu.Lock()
		ds.currentSyncer = nil
		ds.currentTask = nil
		ds.stateMu.Unlock()
		ds.runMu.Unlock()
	}()

	return nil
}

func (ds *LL2DataSyncer) PauseSync() error {
	ds.stateMu.Lock()
	defer ds.stateMu.Unlock()

	if ds.currentSyncer == nil {
		return errors.New("no syncer is running")
	}
	if ds.currentTask == nil || ds.currentTask.Status != models.SyncTaskStatusRunning {
		return errors.New("task is not running")
	}

	ds.currentSyncer.Pause()
	ds.currentTask.Status = models.SyncTaskStatusPaused
	ds.currentTask.UpdatedAt = time.Now()
	return nil
}

func (ds *LL2DataSyncer) ResumeSync() error {
	ds.stateMu.Lock()
	defer ds.stateMu.Unlock()

	if ds.currentSyncer == nil {
		return errors.New("no syncer is running")
	}
	if ds.currentTask == nil || ds.currentTask.Status != models.SyncTaskStatusPaused {
		return errors.New("task is not paused")
	}

	ds.currentSyncer.Resume()
	ds.currentTask.Status = models.SyncTaskStatusRunning
	ds.currentTask.UpdatedAt = time.Now()
	return nil
}

func (ds *LL2DataSyncer) CancelSync() error {
	ds.stateMu.Lock()
	defer ds.stateMu.Unlock()

	if ds.currentSyncer == nil {
		return errors.New("no syncer is running")
	}
	if ds.currentTask != nil &&
		ds.currentTask.Status != models.SyncTaskStatusRunning &&
		ds.currentTask.Status != models.SyncTaskStatusPaused {
		return errors.New("task can only be canceled from running or paused state")
	}

	ds.currentSyncer.Cancel()
	if ds.currentTask != nil {
		ds.currentTask.UpdatedAt = time.Now()
	}
	return nil
}

func (ds *LL2DataSyncer) GetCurrentTask() *TaskInfo {
	ds.stateMu.Lock()
	defer ds.stateMu.Unlock()

	if ds.currentTask == nil {
		return nil
	}

	taskCopy := *ds.currentTask
	if taskCopy.Progress != nil {
		progressCopy := make(map[string]interface{}, len(taskCopy.Progress))
		for k, v := range taskCopy.Progress {
			progressCopy[k] = v
		}
		taskCopy.Progress = progressCopy
	}
	return &taskCopy
}
