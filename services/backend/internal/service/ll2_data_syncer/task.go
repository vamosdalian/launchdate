package ll2datasyncer

import (
	"time"

	"github.com/vamosdalian/launchdate-backend/internal/models"
)

type TaskInfo struct {
	Type      string                 `json:"type"`
	Status    models.SyncTaskStatus  `json:"status"`
	Progress  map[string]interface{} `json:"progress,omitempty"`
	LastError string                 `json:"last_error,omitempty"`
	StartedAt time.Time              `json:"started_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

func isValidSyncType(syncType string) bool {
	switch syncType {
	case SyncTypeLaunch,
		SyncTypeAgency,
		SyncTypeLauncher,
		SyncTypeLauncherFamily,
		SyncTypePad,
		SyncTypeLocation,
		SyncTypeUpcoming:
		return true
	default:
		return false
	}
}
