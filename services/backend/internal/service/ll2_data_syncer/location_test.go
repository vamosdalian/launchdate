package ll2datasyncer

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestLocationSyncer_Sync(t *testing.T) {
	clearCollections(t, "ll2_location", "launch_base")

	rl := &mockRateLimiter{ch: make(chan struct{})}
	syncer := NewLocationSyncer(rl, ll2Service, coreService)

	syncer.Start()

	// Feed the rate limiter until done or timeout
	go func() {
		for {
			select {
			case rl.ch <- struct{}{}:
			case <-syncer.Done():
				return
			}
		}
	}()

	select {
	case <-syncer.Done():
		// Success
	case <-time.After(10 * time.Second):
		t.Fatal("Timeout waiting for syncer to finish")
	}

	ctx := context.Background()

	// Check LL2 Location Collection
	count, err := testDB.Database.Collection("ll2_location").CountDocuments(ctx, bson.M{})
	assert.NoError(t, err)
	assert.Greater(t, count, int64(0), "Should have saved locations to ll2_location")

	// Check Core LaunchBase Collection
	count, err = testDB.Database.Collection("launch_base").CountDocuments(ctx, bson.M{})
	assert.NoError(t, err)
	assert.Greater(t, count, int64(0), "Should have generated core launch bases")
}
