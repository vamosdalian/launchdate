package ll2

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vamosdalian/launchdate-backend/internal/models"
)

func TestLL2Service_GetStats(t *testing.T) {
	s := newTestLL2Service()
	clearCollections(t, LL2LAUNCH, LL2AGENCY, LL2LAUNCHER, LL2LAUNCHERFAMILY, LL2LOCATION, LL2PAD)

	// Context for insertions
	ctx := context.Background()

	// 1. Insert Launches
	launches := []models.LL2LaunchDetailed{
		{
			LL2LaunchNormal: models.LL2LaunchNormal{
				LL2LaunchBasic: models.LL2LaunchBasic{
					ID: "launch1",
				},
			},
		},
		{
			LL2LaunchNormal: models.LL2LaunchNormal{
				LL2LaunchBasic: models.LL2LaunchBasic{
					ID: "launch2",
				},
			},
		},
	}
	for _, l := range launches {
		_, err := s.mongoClient.Collection(LL2LAUNCH).InsertOne(ctx, l)
		assert.NoError(t, err)
	}

	// 2. Insert Agencies
	agencies := []models.LL2AgencyDetailed{
		{
			LL2AgencyNormal: models.LL2AgencyNormal{
				LL2AgencyMini: models.LL2AgencyMini{
					ID: 1,
				},
			},
		},
	}
	for _, a := range agencies {
		_, err := s.mongoClient.Collection(LL2AGENCY).InsertOne(ctx, a)
		assert.NoError(t, err)
	}

	// 3. Insert Launchers
	launchers := []models.LL2LauncherConfigNormal{
		{
			LL2LauncherConfigList: models.LL2LauncherConfigList{
				ID: 10,
			},
		},
		{
			LL2LauncherConfigList: models.LL2LauncherConfigList{
				ID: 11,
			},
		},
		{
			LL2LauncherConfigList: models.LL2LauncherConfigList{
				ID: 12,
			},
		},
	}
	for _, l := range launchers {
		_, err := s.mongoClient.Collection(LL2LAUNCHER).InsertOne(ctx, l)
		assert.NoError(t, err)
	}

	// 4. Insert Launcher Families
	families := []models.LL2LauncherConfigFamilyDetailed{
		{
			LL2LauncherConfigFamilyNormal: models.LL2LauncherConfigFamilyNormal{
				LL2LauncherConfigFamilyMini: models.LL2LauncherConfigFamilyMini{
					ID: 100,
				},
			},
		},
	}
	for _, f := range families {
		_, err := s.mongoClient.Collection(LL2LAUNCHERFAMILY).InsertOne(ctx, f)
		assert.NoError(t, err)
	}

	// 5. Insert Locations
	locations := []models.LL2LocationSerializerWithPads{
		{
			LL2Location: models.LL2Location{ID: 1000},
		},
	}
	for _, l := range locations {
		_, err := s.mongoClient.Collection(LL2LOCATION).InsertOne(ctx, l)
		assert.NoError(t, err)
	}

	// 6. Insert Pads
	pads := []models.LL2Pad{
		{
			LL2PadSerializerNoLocation: models.LL2PadSerializerNoLocation{Id: 2000},
		},
		{
			LL2PadSerializerNoLocation: models.LL2PadSerializerNoLocation{Id: 2001},
		},
	}
	for _, p := range pads {
		_, err := s.mongoClient.Collection(LL2PAD).InsertOne(ctx, p)
		assert.NoError(t, err)
	}

	// Call GetStats
	stats, err := s.GetStats()
	assert.NoError(t, err)

	// Assertions
	assert.Equal(t, 2, stats.Launches, "Launches count mismatch")
	assert.Equal(t, 1, stats.Agencies, "Agencies count mismatch")
	assert.Equal(t, 3, stats.Launchers, "Launchers count mismatch")
	assert.Equal(t, 1, stats.LauncherFamilies, "Launcher Families count mismatch")
	assert.Equal(t, 1, stats.Locations, "Locations count mismatch")
	assert.Equal(t, 2, stats.Pads, "Pads count mismatch")
}
