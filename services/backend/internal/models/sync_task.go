package models

import (
	"time"
)

type SyncTaskStatus string

const (
	SyncTaskStatusIdle    SyncTaskStatus = "idle"
	SyncTaskStatusRunning SyncTaskStatus = "running"
	SyncTaskStatusPaused  SyncTaskStatus = "paused"
	SyncTaskStatusFailed  SyncTaskStatus = "failed"
)

type SyncTask struct {
	ID        string                 `bson:"_id" json:"id"` // Unique identifier for the task type (e.g. "launch_sync")
	Status    SyncTaskStatus         `bson:"status" json:"status"`
	Progress  map[string]interface{} `bson:"progress" json:"progress"` // Store arbitrary progress data
	LastRun   time.Time              `bson:"last_run" json:"last_run"`
	CreatedAt time.Time              `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time              `bson:"updated_at" json:"updated_at"`
}
