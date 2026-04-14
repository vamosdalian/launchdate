package ll2

import (
	"context"
	"time"

	"github.com/vamosdalian/launchdate-backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const SyncTaskCollection = "sync_task"

func (s *LL2Service) GetSyncTask(syncType string) (*models.SyncTask, error) {
	var task models.SyncTask
	err := s.mongoClient.Collection(SyncTaskCollection).
		FindOne(context.Background(), bson.M{"_id": syncType}).
		Decode(&task)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (s *LL2Service) UpsertSyncTask(task *models.SyncTask) error {
	now := time.Now().UTC()
	if task.CreatedAt.IsZero() {
		task.CreatedAt = now
	}
	task.UpdatedAt = now

	update := bson.M{
		"$set": bson.M{
			"type":                      task.Type,
			"status":                    task.Status,
			"progress":                  task.Progress,
			"started_at":                task.StartedAt,
			"finished_at":               task.FinishedAt,
			"last_run":                  task.LastRun,
			"last_success_at":           task.LastSuccessAt,
			"next_run_at":               task.NextRunAt,
			"watermark_last_updated":    task.WatermarkLastUpdated,
			"current_window_start":      task.CurrentWindowStart,
			"current_window_end":        task.CurrentWindowEnd,
			"current_offset":            task.CurrentOffset,
			"current_total":             task.CurrentTotal,
			"max_observed_last_updated": task.MaxObservedLastUpdated,
			"overlap_seconds":           task.OverlapSeconds,
			"last_error":                task.LastError,
			"updated_at":                task.UpdatedAt,
		},
		"$setOnInsert": bson.M{
			"created_at": task.CreatedAt,
		},
	}

	_, err := s.mongoClient.Collection(SyncTaskCollection).UpdateOne(
		context.Background(),
		bson.M{"_id": task.ID},
		update,
		options.Update().SetUpsert(true),
	)
	return err
}

func (s *LL2Service) GetActiveSyncTask() (*models.SyncTask, error) {
	statuses := []models.SyncTaskStatus{models.SyncTaskStatusRunning, models.SyncTaskStatusPaused}
	return s.findOneSyncTask(
		bson.M{"status": bson.M{"$in": statuses}},
		options.FindOne().SetSort(bson.D{{Key: "updated_at", Value: -1}}),
	)
}

func (s *LL2Service) GetLatestVisibleSyncTask() (*models.SyncTask, error) {
	return s.findOneSyncTask(
		bson.M{"status": bson.M{"$ne": models.SyncTaskStatusIdle}},
		options.FindOne().SetSort(bson.D{{Key: "updated_at", Value: -1}}),
	)
}

func (s *LL2Service) ListSyncTasksByStatuses(statuses ...models.SyncTaskStatus) ([]*models.SyncTask, error) {
	filter := bson.M{}
	if len(statuses) > 0 {
		filter["status"] = bson.M{"$in": statuses}
	}

	cursor, err := s.mongoClient.Collection(SyncTaskCollection).Find(
		context.Background(),
		filter,
		options.Find().SetSort(bson.D{{Key: "updated_at", Value: -1}}),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var tasks []*models.SyncTask
	for cursor.Next(context.Background()) {
		var task models.SyncTask
		if err := cursor.Decode(&task); err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *LL2Service) ListRecentSyncTasks(limit int) ([]*models.SyncTask, error) {
	if limit <= 0 {
		limit = 10
	}

	cursor, err := s.mongoClient.Collection(SyncTaskCollection).Find(
		context.Background(),
		bson.M{"status": bson.M{"$ne": models.SyncTaskStatusIdle}},
		options.Find().
			SetSort(bson.D{{Key: "updated_at", Value: -1}}).
			SetLimit(int64(limit)),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var tasks []*models.SyncTask
	for cursor.Next(context.Background()) {
		var task models.SyncTask
		if err := cursor.Decode(&task); err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *LL2Service) findOneSyncTask(filter interface{}, opts *options.FindOneOptions) (*models.SyncTask, error) {
	var task models.SyncTask
	err := s.mongoClient.Collection(SyncTaskCollection).
		FindOne(context.Background(), filter, opts).
		Decode(&task)
	if err != nil {
		return nil, err
	}
	return &task, nil
}
