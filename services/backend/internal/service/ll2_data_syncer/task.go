package ll2datasyncer

import (
	"fmt"
	"strconv"
	"time"

	"github.com/vamosdalian/launchdate-backend/internal/models"
)

const (
	TaskProgressCurrentCount = "current_count"
	TaskProgressTotalCount   = "total_count"
	TaskProgressNextRunAt    = "next_run_at"
)

type TaskProgress struct {
	CurrentCount         int        `json:"current_count"`
	TotalCount           int        `json:"total_count"`
	NextRunAt            *time.Time `json:"next_run_at,omitempty"`
	LastSuccessAt        *time.Time `json:"last_success_at,omitempty"`
	WatermarkLastUpdated *time.Time `json:"watermark_last_updated,omitempty"`
	CurrentWindowStart   *time.Time `json:"current_window_start,omitempty"`
	CurrentWindowEnd     *time.Time `json:"current_window_end,omitempty"`
	OverlapSeconds       int        `json:"overlap_seconds,omitempty"`
}

type TaskInfo struct {
	Type       string                `json:"type"`
	Status     models.SyncTaskStatus `json:"status"`
	Progress   *TaskProgress         `json:"progress,omitempty"`
	LastError  string                `json:"last_error,omitempty"`
	StartedAt  time.Time             `json:"started_at"`
	UpdatedAt  time.Time             `json:"updated_at"`
	FinishedAt *time.Time            `json:"finished_at,omitempty"`
}

func buildCountProgress(current, total int) map[string]interface{} {
	progress := map[string]interface{}{
		TaskProgressCurrentCount: current,
	}
	if total > 0 {
		progress[TaskProgressTotalCount] = total
	}
	return progress
}

func taskInfoFromSyncTask(task *models.SyncTask) *TaskInfo {
	if task == nil {
		return nil
	}

	info := &TaskInfo{
		Type:      task.Type,
		Status:    task.Status,
		LastError: task.LastError,
		StartedAt: task.StartedAt,
		UpdatedAt: task.UpdatedAt,
	}
	if !task.FinishedAt.IsZero() {
		finishedAt := task.FinishedAt.UTC()
		info.FinishedAt = &finishedAt
	}
	if progress := taskProgressFromSyncTask(task); progress != nil {
		info.Progress = progress
	}
	return info
}

func taskProgressFromSyncTask(task *models.SyncTask) *TaskProgress {
	if task == nil {
		return nil
	}

	progress := &TaskProgress{}
	if current, ok := intValue(task.Progress, TaskProgressCurrentCount); ok {
		progress.CurrentCount = current
	}
	if total, ok := intValue(task.Progress, TaskProgressTotalCount); ok {
		progress.TotalCount = total
	}
	if progress.CurrentCount == 0 && task.CurrentOffset > 0 {
		progress.CurrentCount = task.CurrentOffset
	}
	if progress.TotalCount == 0 && task.CurrentTotal > 0 {
		progress.TotalCount = task.CurrentTotal
	}
	if !task.NextRunAt.IsZero() {
		nextRunAt := task.NextRunAt.UTC()
		progress.NextRunAt = &nextRunAt
	}
	if !task.LastSuccessAt.IsZero() {
		lastSuccessAt := task.LastSuccessAt.UTC()
		progress.LastSuccessAt = &lastSuccessAt
	}
	if !task.WatermarkLastUpdated.IsZero() {
		watermark := task.WatermarkLastUpdated.UTC()
		progress.WatermarkLastUpdated = &watermark
	}
	if !task.CurrentWindowStart.IsZero() {
		windowStart := task.CurrentWindowStart.UTC()
		progress.CurrentWindowStart = &windowStart
	}
	if !task.CurrentWindowEnd.IsZero() {
		windowEnd := task.CurrentWindowEnd.UTC()
		progress.CurrentWindowEnd = &windowEnd
	}
	if task.OverlapSeconds > 0 {
		progress.OverlapSeconds = task.OverlapSeconds
	}

	if progress.CurrentCount == 0 && progress.TotalCount == 0 &&
		progress.NextRunAt == nil && progress.LastSuccessAt == nil &&
		progress.WatermarkLastUpdated == nil && progress.CurrentWindowStart == nil &&
		progress.CurrentWindowEnd == nil && progress.OverlapSeconds == 0 {
		return nil
	}

	return progress
}

func intValue(values map[string]interface{}, key string) (int, bool) {
	if len(values) == 0 {
		return 0, false
	}

	raw, ok := values[key]
	if !ok {
		return 0, false
	}

	switch value := raw.(type) {
	case int:
		return value, true
	case int32:
		return int(value), true
	case int64:
		return int(value), true
	case float64:
		return int(value), true
	case float32:
		return int(value), true
	case string:
		parsed, err := strconv.Atoi(value)
		if err != nil {
			return 0, false
		}
		return parsed, true
	default:
		parsed, err := strconv.Atoi(fmt.Sprint(value))
		if err != nil {
			return 0, false
		}
		return parsed, true
	}
}

func isValidSyncType(syncType string) bool {
	switch syncType {
	case SyncTypeLaunch,
		SyncTypeAgency,
		SyncTypeLauncher,
		SyncTypeLauncherFamily,
		SyncTypePad,
		SyncTypeLocation,
		SyncTypeUpdate:
		return true
	default:
		return false
	}
}
