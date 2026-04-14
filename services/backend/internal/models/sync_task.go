package models

import (
	"time"
)

type SyncTaskStatus string

const (
	SyncTaskStatusIdle      SyncTaskStatus = "idle"
	SyncTaskStatusRunning   SyncTaskStatus = "running"
	SyncTaskStatusPaused    SyncTaskStatus = "paused"
	SyncTaskStatusCompleted SyncTaskStatus = "completed"
	SyncTaskStatusCanceled  SyncTaskStatus = "canceled"
	SyncTaskStatusFailed    SyncTaskStatus = "failed"
)

type SyncTask struct {
	ID                     string                 `bson:"_id" json:"id"`
	Type                   string                 `bson:"type" json:"type"`
	Status                 SyncTaskStatus         `bson:"status" json:"status"`
	Progress               map[string]interface{} `bson:"progress" json:"progress"`
	StartedAt              time.Time              `bson:"started_at" json:"started_at"`
	FinishedAt             time.Time              `bson:"finished_at" json:"finished_at"`
	LastRun                time.Time              `bson:"last_run" json:"last_run"`
	LastSuccessAt          time.Time              `bson:"last_success_at" json:"last_success_at"`
	NextRunAt              time.Time              `bson:"next_run_at" json:"next_run_at"`
	WatermarkLastUpdated   time.Time              `bson:"watermark_last_updated" json:"watermark_last_updated"`
	CurrentWindowStart     time.Time              `bson:"current_window_start" json:"current_window_start"`
	CurrentWindowEnd       time.Time              `bson:"current_window_end" json:"current_window_end"`
	CurrentOffset          int                    `bson:"current_offset" json:"current_offset"`
	CurrentTotal           int                    `bson:"current_total" json:"current_total"`
	MaxObservedLastUpdated time.Time              `bson:"max_observed_last_updated" json:"max_observed_last_updated"`
	OverlapSeconds         int                    `bson:"overlap_seconds" json:"overlap_seconds"`
	LastError              string                 `bson:"last_error" json:"last_error"`
	CreatedAt              time.Time              `bson:"created_at" json:"created_at"`
	UpdatedAt              time.Time              `bson:"updated_at" json:"updated_at"`
}

func (s SyncTaskStatus) IsActive() bool {
	return s == SyncTaskStatusRunning || s == SyncTaskStatusPaused
}

func (s SyncTaskStatus) IsTerminal() bool {
	return s == SyncTaskStatusCompleted || s == SyncTaskStatusCanceled || s == SyncTaskStatusFailed
}
