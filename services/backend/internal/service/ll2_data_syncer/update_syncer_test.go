package ll2datasyncer

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vamosdalian/launchdate-backend/internal/config"
	"github.com/vamosdalian/launchdate-backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
)

func TestUpdateSyncer_SyncPersistsState(t *testing.T) {
	clearCollections(t, "ll2_launch", "launch", "sync_task")

	syncer := NewUpdateSyncer(&mockRateLimiter{ch: make(chan struct{})}, ll2Service, coreService)
	syncer.sync()

	ctx := context.Background()

	launchCount, err := testDB.Collection("ll2_launch").CountDocuments(ctx, bson.M{})
	require.NoError(t, err)
	assert.Greater(t, launchCount, int64(0), "should save updated launches to ll2_launch")

	coreCount, err := testDB.Collection("launch").CountDocuments(ctx, bson.M{})
	require.NoError(t, err)
	assert.Greater(t, coreCount, int64(0), "should ensure core launch documents exist")

	task, err := ll2Service.GetSyncTask(SyncTypeUpdate)
	require.NoError(t, err)
	assert.Equal(t, SyncTypeUpdate, task.Type)
	assert.Equal(t, models.SyncTaskStatusRunning, task.Status)
	assert.False(t, task.WatermarkLastUpdated.IsZero(), "watermark should be persisted after a completed cycle")
	assert.False(t, task.NextRunAt.IsZero(), "next run time should be persisted")
	assert.True(t, task.CurrentWindowStart.IsZero(), "completed cycle should clear current window start")
	assert.True(t, task.CurrentWindowEnd.IsZero(), "completed cycle should clear current window end")
	assert.Equal(t, 0, task.CurrentOffset)
	assert.NotEmpty(t, task.Progress["watermark_last_updated"])
	assert.NotEmpty(t, task.Progress["next_run_at"])
	assert.EqualValues(t, int(defaultUpdateOverlap/time.Second), task.Progress["overlap_seconds"])
}

func TestUpdateSyncer_NextRunInterval(t *testing.T) {
	syncer := NewUpdateSyncer(&mockRateLimiter{ch: make(chan struct{})}, ll2Service, coreService)
	now := time.Date(2026, 4, 13, 12, 0, 0, 0, time.UTC)

	t.Run("hot window uses five minutes", func(t *testing.T) {
		clearCollections(t, "ll2_launch")
		insertLaunchNet(t, "hot-launch", now.Add(2*time.Hour))
		assert.Equal(t, updateIntervalHot, syncer.nextRunInterval(now))
	})

	t.Run("near window uses fifteen minutes", func(t *testing.T) {
		clearCollections(t, "ll2_launch")
		insertLaunchNet(t, "near-launch", now.Add(12*time.Hour))
		assert.Equal(t, updateIntervalNear, syncer.nextRunInterval(now))
	})

	t.Run("upcoming week uses thirty minutes", func(t *testing.T) {
		clearCollections(t, "ll2_launch")
		insertLaunchNet(t, "future-launch", now.Add(3*24*time.Hour))
		assert.Equal(t, updateIntervalWeekly, syncer.nextRunInterval(now))
	})

	t.Run("no nearby launches uses sixty minutes", func(t *testing.T) {
		clearCollections(t, "ll2_launch")
		assert.Equal(t, updateIntervalIdle, syncer.nextRunInterval(now))
	})
}

func TestLL2DataSyncer_RestoreTasks(t *testing.T) {
	clearCollections(t, "sync_task")

	now := time.Now().UTC()
	require.NoError(t, ll2Service.UpsertSyncTask(&models.SyncTask{
		ID:                   SyncTypeUpdate,
		Type:                 SyncTypeUpdate,
		Status:               models.SyncTaskStatusRunning,
		NextRunAt:            now,
		WatermarkLastUpdated: now.Add(-time.Hour),
		OverlapSeconds:       int(defaultUpdateOverlap / time.Second),
	}))

	ds := NewLL2DataSyncer(&config.Config{LL2RequestInterval: 1}, ll2Service, coreService)
	require.NoError(t, ds.RestoreTasks())

	task := ds.GetCurrentTask()
	require.NotNil(t, task)
	assert.Equal(t, SyncTypeUpdate, task.Type)
	assert.Equal(t, models.SyncTaskStatusRunning, task.Status)

	_, err := ds.CancelSync()
	require.NoError(t, err)
	select {
	case <-ds.currentSyncer.Done():
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for restored update sync to stop")
	}
}

func TestLL2DataSyncer_CancelPersistedUpdateWithoutSyncer(t *testing.T) {
	clearCollections(t, "sync_task")

	now := time.Now().UTC()
	require.NoError(t, ll2Service.UpsertSyncTask(&models.SyncTask{
		ID:             SyncTypeUpdate,
		Type:           SyncTypeUpdate,
		Status:         models.SyncTaskStatusRunning,
		StartedAt:      now.Add(-5 * time.Minute),
		NextRunAt:      now.Add(30 * time.Minute),
		CurrentOffset:  200,
		CurrentTotal:   300,
		Progress:       buildCountProgress(200, 300),
		OverlapSeconds: int(defaultUpdateOverlap / time.Second),
	}))

	ds := NewLL2DataSyncer(&config.Config{LL2RequestInterval: 1}, ll2Service, coreService)
	_, err := ds.CancelSync()
	require.NoError(t, err)

	task, err := ll2Service.GetSyncTask(SyncTypeUpdate)
	require.NoError(t, err)
	assert.Equal(t, models.SyncTaskStatusCanceled, task.Status)
	assert.True(t, task.NextRunAt.IsZero())
	assert.Equal(t, 200, task.CurrentOffset)
	assert.Equal(t, 300, task.CurrentTotal)
	assert.NotNil(t, task.Progress)
	assert.Nil(t, ds.GetCurrentTask())
}

func TestLL2DataSyncer_RestoreTasksMarksOneOffTasksFailed(t *testing.T) {
	clearCollections(t, "sync_task")

	now := time.Now().UTC()
	require.NoError(t, ll2Service.UpsertSyncTask(&models.SyncTask{
		ID:        SyncTypeLaunch,
		Type:      SyncTypeLaunch,
		Status:    models.SyncTaskStatusRunning,
		StartedAt: now.Add(-time.Minute),
	}))

	ds := NewLL2DataSyncer(&config.Config{LL2RequestInterval: 1}, ll2Service, coreService)
	require.NoError(t, ds.RestoreTasks())

	task, err := ll2Service.GetSyncTask(SyncTypeLaunch)
	require.NoError(t, err)
	assert.Equal(t, models.SyncTaskStatusFailed, task.Status)
	assert.Equal(t, interruptedTaskMessage, task.LastError)
	assert.False(t, task.FinishedAt.IsZero())
	assert.Nil(t, ds.GetCurrentTask())
}

func insertLaunchNet(t *testing.T, id string, net time.Time) {
	t.Helper()
	_, err := testDB.Collection("ll2_launch").InsertOne(context.Background(), bson.M{
		"id":   id,
		"name": id,
		"net":  net.UTC().Format(time.RFC3339),
	})
	require.NoError(t, err)
}
