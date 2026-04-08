package ll2datasyncer

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestLauncherSyncer_Sync(t *testing.T) {
	clearCollections(t, "ll2_launcher", "rocket")

	rl := &mockRateLimiter{ch: make(chan struct{})}
	syncer := NewLauncherSyncer(rl, ll2Service, coreService)

	syncer.Start()

	go func() {
		rl.ch <- struct{}{}
	}()

	select {
	case <-syncer.Done():
		// Success
	case <-time.After(10 * time.Second):
		t.Fatal("Timeout waiting for syncer to finish")
	}

	ctx := context.Background()

	// Check LL2 Launcher Collection
	count, err := testDB.Database.Collection("ll2_launcher").CountDocuments(ctx, bson.M{})
	assert.NoError(t, err)
	assert.Greater(t, count, int64(0), "Should have saved launchers to ll2_launcher")

	// Check Core Rockets Collection
	count, err = testDB.Database.Collection("rocket").CountDocuments(ctx, bson.M{})
	assert.NoError(t, err)
	assert.Greater(t, count, int64(0), "Should have generated core rockets")
}
