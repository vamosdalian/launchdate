package ll2datasyncer

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vamosdalian/launchdate-backend/internal/models"
)

func TestTaskInfoFromSyncTask_UsesStartedAt(t *testing.T) {
	startedAt := time.Date(2026, 4, 13, 7, 0, 0, 0, time.UTC)
	updatedAt := startedAt.Add(10 * time.Minute)

	info := taskInfoFromSyncTask(&models.SyncTask{
		Type:      SyncTypeUpdate,
		Status:    models.SyncTaskStatusRunning,
		StartedAt: startedAt,
		UpdatedAt: updatedAt,
	})

	assert.Equal(t, startedAt, info.StartedAt)
	assert.Equal(t, updatedAt, info.UpdatedAt)
}

func TestTaskInfoFromSyncTask_MapsTypedProgress(t *testing.T) {
	now := time.Date(2026, 4, 13, 7, 0, 0, 0, time.UTC)
	nextRunAt := now.Add(5 * time.Minute)
	lastSuccessAt := now.Add(-time.Minute)
	finishedAt := now.Add(10 * time.Minute)

	info := taskInfoFromSyncTask(&models.SyncTask{
		Type:          SyncTypeUpdate,
		Status:        models.SyncTaskStatusCompleted,
		StartedAt:     now,
		UpdatedAt:     now,
		FinishedAt:    finishedAt,
		LastSuccessAt: lastSuccessAt,
		NextRunAt:     nextRunAt,
		CurrentOffset: 200,
		CurrentTotal:  300,
		Progress:      buildCountProgress(200, 300),
	})

	if assert.NotNil(t, info.Progress) {
		assert.Equal(t, 200, info.Progress.CurrentCount)
		assert.Equal(t, 300, info.Progress.TotalCount)
		assert.NotNil(t, info.Progress.NextRunAt)
		assert.Equal(t, nextRunAt, *info.Progress.NextRunAt)
		assert.NotNil(t, info.Progress.LastSuccessAt)
		assert.Equal(t, lastSuccessAt, *info.Progress.LastSuccessAt)
	}
	if assert.NotNil(t, info.FinishedAt) {
		assert.Equal(t, finishedAt, *info.FinishedAt)
	}
}
