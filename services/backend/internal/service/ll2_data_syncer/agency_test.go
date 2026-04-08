package ll2datasyncer

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestAgencySyncer_Sync(t *testing.T) {
	clearCollections(t, "ll2_agency", "agency")

	rl := &mockRateLimiter{ch: make(chan struct{})}
	syncer := NewAgencySyncer(rl, ll2Service, coreService)

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

	// Check LL2 Agency Collection
	count, err := testDB.Database.Collection("ll2_agency").CountDocuments(ctx, bson.M{})
	assert.NoError(t, err)
	assert.Greater(t, count, int64(0), "Should have saved agencies to ll2_agency")

	// Check Core Agency Collection
	count, err = testDB.Database.Collection("agency").CountDocuments(ctx, bson.M{})
	assert.NoError(t, err)
	assert.Greater(t, count, int64(0), "Should have generated core agencies")
}
