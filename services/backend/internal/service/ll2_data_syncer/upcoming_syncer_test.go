package ll2datasyncer

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestUpcomingSyncer_Sync(t *testing.T) {
	clearCollections(t, "ll2_launch")

	// Insert a dummy "latest" launch to test the "GetLatestLaunchFromDB" logic
	// We insert one with a future net date
	_, err := testDB.Collection("ll2_launch").InsertOne(context.Background(), bson.M{
		"id":   "test-latest-launch-id",
		"name": "Test Latest Launch",
		"net":  "2099-01-01T00:00:00Z",
	})
	assert.NoError(t, err)

	rl := &mockRateLimiter{ch: make(chan struct{})}
	syncer := NewUpcomingSyncer(rl, ll2Service, coreService)

	syncer.Start()

	// Allow one request (Fetch upcoming launches)
	go func() {
		rl.ch <- struct{}{}
	}()

	// Wait for processing
	// We can't rely on Syncer.Done() unless we stop it.
	// The syncer runs forever. We cancel it after a short delay.
	time.Sleep(500 * time.Millisecond)
	syncer.Cancel()
	<-syncer.Done()

	ctx := context.Background()

	// Check LL2 Launch Collection
	// existing 1 + newly fetched
	count, err := testDB.Database.Collection("ll2_launch").CountDocuments(ctx, bson.M{})
	assert.NoError(t, err)
	assert.Greater(t, count, int64(1), "Should have saved fetched upcoming launches")

	// Check that we logged the latest launch (This would require capturing logs, skipping for now, relying on code execution path)

	// Check that a launch from launch.json exists
	var launch bson.M
	// launch.json usually has some known IDs.
	// Assuming launch.json content is loaded.
	err = testDB.Database.Collection("ll2_launch").FindOne(ctx, bson.M{"id": bson.M{"$ne": "test-latest-launch-id"}}).Decode(&launch)
	assert.NoError(t, err, "Should find at least one synced launch")
}
