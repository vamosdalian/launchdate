package ll2datasyncer

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestLauncherFamilySyncer_Sync(t *testing.T) {
	clearCollections(t, "ll2_launcher_family")

	rl := &mockRateLimiter{ch: make(chan struct{})}
	syncer := NewLauncherFamilySyncer(rl, ll2Service, nil)

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

	// Check LL2 Launcher Family Collection
	count, err := testDB.Database.Collection("ll2_launcher_family").CountDocuments(ctx, bson.M{})
	assert.NoError(t, err)
	assert.Greater(t, count, int64(0), "Should have saved launcher families to ll2_launcher_family")
}
