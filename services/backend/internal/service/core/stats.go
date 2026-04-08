package core

import (
	"context"
	"time"

	"github.com/vamosdalian/launchdate-backend/internal/models"
)

func (m *MainService) GetStats() (models.Stats, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var rocketCount, launchCount, agencyCount, launchBaseCount int64

	rocketCount, _ = m.mc.Collection(COLLECTION_ROCKET).EstimatedDocumentCount(ctx)
	launchCount, _ = m.mc.Collection(COLLECTION_LAUNCH).EstimatedDocumentCount(ctx)
	agencyCount, _ = m.mc.Collection(COLLECTION_AGENCY).EstimatedDocumentCount(ctx)
	launchBaseCount, _ = m.mc.Collection(COLLECTION_LAUNCH_BASE).EstimatedDocumentCount(ctx)

	return models.Stats{
		Rockets:     int(rocketCount),
		Launches:    int(launchCount),
		Agencies:    int(agencyCount),
		LaunchBases: int(launchBaseCount),
	}, nil
}
