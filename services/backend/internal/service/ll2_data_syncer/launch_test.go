package ll2datasyncer

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestLaunchSyncer_Sync(t *testing.T) {
	clearCollections(t, "ll2_launch", "launch")

	rl := &mockRateLimiter{ch: make(chan struct{})}
	syncer := NewLaunchSyncer(rl, ll2Service, coreService, nil)

	syncer.Start()

	// Allow one request
	go func() {
		rl.ch <- struct{}{}
	}()

	select {
	case <-syncer.Done():
		// Success
	case <-time.After(10 * time.Second): // Generous timeout
		t.Fatal("Timeout waiting for syncer to finish")
	}

	ctx := context.Background()

	// Check LL2 Launch Collection
	count, err := testDB.Database.Collection("ll2_launch").CountDocuments(ctx, bson.M{})
	assert.NoError(t, err)
	// We expect some documents. Not asserting exact count 10 to be robust against testdata changes, but > 0
	assert.Greater(t, count, int64(0), "Should have saved launches to ll2_launch")

	// Check Core Launch Collection
	count, err = testDB.Database.Collection("launch").CountDocuments(ctx, bson.M{})
	assert.NoError(t, err)
	assert.Greater(t, count, int64(0), "Should have generated core launches")
}
