package ll2

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vamosdalian/launchdate-backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
)

func TestLL2Service_ListRecentSyncTasks(t *testing.T) {
	s := newTestLL2Service()
	clearCollections(t, SyncTaskCollection)

	now := time.Now().UTC()
	tasks := []models.SyncTask{
		{
			ID:        "launch",
			Type:      "launch",
			Status:    models.SyncTaskStatusCompleted,
			StartedAt: now.Add(-30 * time.Minute),
			UpdatedAt: now.Add(-20 * time.Minute),
		},
		{
			ID:        "agency",
			Type:      "agency",
			Status:    models.SyncTaskStatusFailed,
			StartedAt: now.Add(-15 * time.Minute),
			UpdatedAt: now.Add(-10 * time.Minute),
			LastError: "boom",
		},
		{
			ID:        "pad",
			Type:      "pad",
			Status:    models.SyncTaskStatusIdle,
			StartedAt: now.Add(-5 * time.Minute),
			UpdatedAt: now.Add(-5 * time.Minute),
		},
		{
			ID:        "update",
			Type:      "update",
			Status:    models.SyncTaskStatusRunning,
			StartedAt: now.Add(-2 * time.Minute),
			UpdatedAt: now.Add(-time.Minute),
		},
	}

	for _, task := range tasks {
		_, err := s.mongoClient.Collection(SyncTaskCollection).InsertOne(context.Background(), task)
		require.NoError(t, err)
	}

	recent, err := s.ListRecentSyncTasks(2)
	require.NoError(t, err)
	require.Len(t, recent, 2)
	assert.Equal(t, "update", recent[0].Type)
	assert.Equal(t, "agency", recent[1].Type)
}

func TestLL2Service_GetLatestVisibleSyncTask(t *testing.T) {
	s := newTestLL2Service()
	clearCollections(t, SyncTaskCollection)

	now := time.Now().UTC()
	_, err := s.mongoClient.Collection(SyncTaskCollection).InsertMany(context.Background(), []interface{}{
		bson.M{"_id": "launch", "type": "launch", "status": models.SyncTaskStatusCompleted, "updated_at": now.Add(-2 * time.Minute)},
		bson.M{"_id": "update", "type": "update", "status": models.SyncTaskStatusRunning, "updated_at": now},
		bson.M{"_id": "pad", "type": "pad", "status": models.SyncTaskStatusIdle, "updated_at": now.Add(time.Minute)},
	})
	require.NoError(t, err)

	task, err := s.GetLatestVisibleSyncTask()
	require.NoError(t, err)
	assert.Equal(t, "update", task.Type)
}
