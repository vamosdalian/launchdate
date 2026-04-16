package ll2datasyncer

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/vamosdalian/launchdate-backend/internal/models"
	"github.com/vamosdalian/launchdate-backend/internal/service/core"
	"github.com/vamosdalian/launchdate-backend/internal/service/ll2"
	"github.com/vamosdalian/launchdate-backend/internal/util"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	defaultUpdateOverlap         = 15 * time.Minute
	defaultUpdateInitialLookback = 24 * time.Hour
	updateHotPastWindow          = 1 * time.Hour
	updateHotFutureWindow        = 6 * time.Hour
	updateNearWindow             = 24 * time.Hour
	updateWeeklyWindow           = 7 * 24 * time.Hour
	updateIntervalHot            = 5 * time.Minute
	updateIntervalNear           = 15 * time.Minute
	updateIntervalWeekly         = 30 * time.Minute
	updateIntervalIdle           = 60 * time.Minute
)

type UpdateSyncer struct {
	*BaseSyncer
	ll2Service *ll2.LL2Service
	core       *core.MainService
}

func NewUpdateSyncer(rl util.RateLimiter, ll2Service *ll2.LL2Service, core *core.MainService) *UpdateSyncer {
	s := &UpdateSyncer{
		ll2Service: ll2Service,
		core:       core,
	}
	s.BaseSyncer = NewBaseSyncer(rl, s.sync)
	return s
}

func (s *UpdateSyncer) sync() {
	task, err := s.loadTask()
	if err != nil {
		logrus.Errorf("load update sync task failed: %v", err)
		return
	}

	if task.Status != models.SyncTaskStatusRunning {
		return
	}

	now := time.Now().UTC()
	if task.CurrentWindowStart.IsZero() {
		if !task.NextRunAt.IsZero() && now.Before(task.NextRunAt) {
			return
		}

		task.LastRun = now
		task.LastError = ""
		task.CurrentOffset = 0
		task.CurrentWindowEnd = now

		since := task.WatermarkLastUpdated
		if since.IsZero() {
			since = now.Add(-defaultUpdateInitialLookback)
		}
		task.CurrentWindowStart = since.Add(-time.Duration(task.OverlapSeconds) * time.Second)

		if err := s.saveTask(task); err != nil {
			logrus.Errorf("persist update sync start failed: %v", err)
			return
		}
	}

	resp, err := s.ll2Service.GetLaunchesUpdatedFromAPI(task.CurrentWindowStart, task.CurrentWindowEnd, PageSize, task.CurrentOffset)
	if err != nil {
		task.LastError = err.Error()
		if saveErr := s.saveTask(task); saveErr != nil {
			logrus.Errorf("persist update sync error failed: %v", saveErr)
		}
		logrus.Errorf("get updated launches from ll2 api failed: %v", err)
		return
	}
	task.CurrentTotal = resp.Count

	for _, launch := range resp.Results {
		if err := s.ll2Service.SaveLaunchesToDB([]*models.LL2LaunchDetailed{launch}); err != nil {
			logrus.Errorf("save updated launch %s failed: %v", launch.ID, err)
			continue
		}

		if err := s.core.GenerateLaunchFromLL2([]string{launch.ID}); err != nil {
			logrus.Errorf("generate core launch %s failed: %v", launch.ID, err)
		}
	}

	task.LastError = ""
	task.CurrentOffset += len(resp.Results)

	if len(resp.Results) < PageSize {
		task.WatermarkLastUpdated = task.CurrentWindowEnd
		task.LastSuccessAt = now
		task.NextRunAt = now.Add(s.nextRunInterval(now))
		task.CurrentWindowStart = time.Time{}
		task.CurrentWindowEnd = time.Time{}
		task.CurrentOffset = 0
		task.CurrentTotal = 0
		task.MaxObservedLastUpdated = time.Time{}
		s.core.InvalidatePublicCacheForSync(SyncTypeUpdate)
	}

	if err := s.saveTask(task); err != nil {
		logrus.Errorf("persist update sync progress failed: %v", err)
	}
}

func (s *UpdateSyncer) loadTask() (*models.SyncTask, error) {
	task, err := s.ll2Service.GetSyncTask(SyncTypeUpdate)
	if err == nil {
		if task.Type == "" {
			task.Type = SyncTypeUpdate
		}
		if task.OverlapSeconds == 0 {
			task.OverlapSeconds = int(defaultUpdateOverlap / time.Second)
		}
		return task, nil
	}
	if err != mongo.ErrNoDocuments {
		return nil, err
	}

	return &models.SyncTask{
		ID:             SyncTypeUpdate,
		Type:           SyncTypeUpdate,
		Status:         models.SyncTaskStatusRunning,
		OverlapSeconds: int(defaultUpdateOverlap / time.Second),
	}, nil
}

func (s *UpdateSyncer) saveTask(task *models.SyncTask) error {
	task.Progress = buildSyncProgress(task)
	return s.ll2Service.UpsertSyncTask(task)
}

func (s *UpdateSyncer) nextRunInterval(now time.Time) time.Duration {
	hasHot, err := s.ll2Service.HasLaunchesInNetWindow(now.Add(-updateHotPastWindow), now.Add(updateHotFutureWindow))
	if err != nil {
		logrus.Errorf("check hot launch window failed: %v", err)
		return updateIntervalIdle
	}
	if hasHot {
		return updateIntervalHot
	}

	hasNear, err := s.ll2Service.HasLaunchesInNetWindow(now.Add(-updateNearWindow), now.Add(updateNearWindow))
	if err != nil {
		logrus.Errorf("check near launch window failed: %v", err)
		return updateIntervalIdle
	}
	if hasNear {
		return updateIntervalNear
	}

	hasWeekly, err := s.ll2Service.HasLaunchesInNetWindow(now, now.Add(updateWeeklyWindow))
	if err != nil {
		logrus.Errorf("check weekly launch window failed: %v", err)
		return updateIntervalIdle
	}
	if hasWeekly {
		return updateIntervalWeekly
	}

	return updateIntervalIdle
}

func buildSyncProgress(task *models.SyncTask) map[string]interface{} {
	progress := map[string]interface{}{
		"overlap_seconds": task.OverlapSeconds,
	}

	if !task.WatermarkLastUpdated.IsZero() {
		progress["watermark_last_updated"] = task.WatermarkLastUpdated.UTC().Format(time.RFC3339)
	}
	if !task.CurrentWindowStart.IsZero() {
		progress["current_window_start"] = task.CurrentWindowStart.UTC().Format(time.RFC3339)
	}
	if !task.CurrentWindowEnd.IsZero() {
		progress["current_window_end"] = task.CurrentWindowEnd.UTC().Format(time.RFC3339)
	}
	if !task.NextRunAt.IsZero() {
		progress["next_run_at"] = task.NextRunAt.UTC().Format(time.RFC3339)
	}
	if task.CurrentOffset > 0 {
		progress[TaskProgressCurrentCount] = task.CurrentOffset
	}
	if task.CurrentTotal > 0 {
		progress[TaskProgressTotalCount] = task.CurrentTotal
	}
	if task.CurrentOffset > 0 {
		progress["current_offset"] = task.CurrentOffset
	}
	if !task.LastSuccessAt.IsZero() {
		progress["last_success_at"] = task.LastSuccessAt.UTC().Format(time.RFC3339)
	}

	return progress
}
