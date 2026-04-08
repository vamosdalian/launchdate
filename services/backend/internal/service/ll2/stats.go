package ll2

import (
	"context"
	"time"

	"github.com/vamosdalian/launchdate-backend/internal/models"
)

func (s *LL2Service) GetStats() (models.LL2Stats, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var launches, agencies, launchers, launcherFamilies, locations, pads int64

	launches, _ = s.mongoClient.Collection(LL2LAUNCH).EstimatedDocumentCount(ctx)
	agencies, _ = s.mongoClient.Collection(LL2AGENCY).EstimatedDocumentCount(ctx)
	launchers, _ = s.mongoClient.Collection(LL2LAUNCHER).EstimatedDocumentCount(ctx)
	launcherFamilies, _ = s.mongoClient.Collection(LL2LAUNCHERFAMILY).EstimatedDocumentCount(ctx)
	locations, _ = s.mongoClient.Collection(LL2LOCATION).EstimatedDocumentCount(ctx)
	pads, _ = s.mongoClient.Collection(LL2PAD).EstimatedDocumentCount(ctx)

	return models.LL2Stats{
		Launches:         int(launches),
		Agencies:         int(agencies),
		Launchers:        int(launchers),
		LauncherFamilies: int(launcherFamilies),
		Locations:        int(locations),
		Pads:             int(pads),
	}, nil
}
