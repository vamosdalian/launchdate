package ll2datasyncer

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestPadSyncer_Sync(t *testing.T) {
	clearCollections(t, "ll2_pad")

	rl := &mockRateLimiter{ch: make(chan struct{})}
	syncer := NewPadSyncer(rl, ll2Service, nil)

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

	// Check LL2 Pad Collection
	count, err := testDB.Database.Collection("ll2_pad").CountDocuments(ctx, bson.M{})
	assert.NoError(t, err)
	assert.Greater(t, count, int64(0), "Should have saved pads to ll2_pad")
}
