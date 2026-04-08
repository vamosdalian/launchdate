package core

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"github.com/vamosdalian/launchdate-backend/internal/db"
	"github.com/vamosdalian/launchdate-backend/internal/models"
)

func TestMainService_GetStats(t *testing.T) {
	ctx := context.Background()

	// 1. Setup Mongo Container
	mongoContainer, err := mongodb.Run(ctx, "mongo:6")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := mongoContainer.Terminate(ctx); err != nil {
			t.Fatal(err)
		}
	}()

	uri, err := mongoContainer.ConnectionString(ctx)
	if err != nil {
		t.Fatal(err)
	}

	testDB, cleanup, err := db.NewMongoDB(uri, "test_core_db")
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	// 2. Initialize Service
	service := NewMainService(testDB)

	// 3. Clear existing collections (if any)
	collections := []string{COLLECTION_ROCKET, COLLECTION_LAUNCH, COLLECTION_AGENCY, COLLECTION_LAUNCH_BASE}
	for _, c := range collections {
		testDB.Database.Collection(c).Drop(ctx)
	}

	// 4. Test with empty collections
	stats, err := service.GetStats()
	assert.NoError(t, err)
	assert.Equal(t, 0, stats.Rockets)
	assert.Equal(t, 0, stats.Launches)
	assert.Equal(t, 0, stats.Agencies)
	assert.Equal(t, 0, stats.LaunchBases)

	// 5. Insert Data
	// Rocket
	_, err = testDB.Collection(COLLECTION_ROCKET).InsertOne(ctx, models.Rocket{ID: 1})
	assert.NoError(t, err)

	// Launch
	_, err = testDB.Collection(COLLECTION_LAUNCH).InsertOne(ctx, models.Launch{ID: 1})
	assert.NoError(t, err)
	_, err = testDB.Collection(COLLECTION_LAUNCH).InsertOne(ctx, models.Launch{ID: 2})
	assert.NoError(t, err)

	// Agency
	_, err = testDB.Collection(COLLECTION_AGENCY).InsertOne(ctx, models.Agency{ID: 1})
	assert.NoError(t, err)

	// LaunchBase
	// (Simulate LaunchBase insertion)
	_, err = testDB.Collection(COLLECTION_LAUNCH_BASE).InsertOne(ctx, models.LaunchBase{ID: 1})
	assert.NoError(t, err)

	// 6. Test with populated collections
	stats, err = service.GetStats()
	assert.NoError(t, err)
	assert.Equal(t, 1, stats.Rockets)
	assert.Equal(t, 2, stats.Launches)
	assert.Equal(t, 1, stats.Agencies)
	assert.Equal(t, 1, stats.LaunchBases)
}
