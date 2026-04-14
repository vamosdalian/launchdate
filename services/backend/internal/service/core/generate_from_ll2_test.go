package core

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"github.com/vamosdalian/launchdate-backend/internal/db"
	"github.com/vamosdalian/launchdate-backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
)

func TestMainService_GenerateFromLL2IsIdempotent(t *testing.T) {
	ctx := context.Background()

	mongoContainer, err := mongodb.Run(ctx, "mongo:6")
	require.NoError(t, err)
	defer func() {
		require.NoError(t, mongoContainer.Terminate(ctx))
	}()

	uri, err := mongoContainer.ConnectionString(ctx)
	require.NoError(t, err)

	testDB, cleanup, err := db.NewMongoDB(uri, "test_core_generate_ll2")
	require.NoError(t, err)
	defer cleanup()

	service := NewMainService(testDB)

	t.Run("launch", func(t *testing.T) {
		dropCollection(t, ctx, testDB, COLLECTION_LAUNCH)

		require.NoError(t, service.GenerateLaunchFromLL2([]string{"launch-1", "launch-1"}))

		var first models.Launch
		require.NoError(t, testDB.Collection(COLLECTION_LAUNCH).FindOne(ctx, bson.M{"external_id": "launch-1"}).Decode(&first))

		_, err := testDB.Collection(COLLECTION_LAUNCH).UpdateOne(ctx, bson.M{"external_id": "launch-1"}, bson.M{"$set": bson.M{"background_image": "keep-me"}})
		require.NoError(t, err)

		require.NoError(t, service.GenerateLaunchFromLL2([]string{"launch-1"}))

		count, err := testDB.Collection(COLLECTION_LAUNCH).CountDocuments(ctx, bson.M{"external_id": "launch-1"})
		require.NoError(t, err)
		assert.EqualValues(t, 1, count)

		var got models.Launch
		require.NoError(t, testDB.Collection(COLLECTION_LAUNCH).FindOne(ctx, bson.M{"external_id": "launch-1"}).Decode(&got))
		assert.Equal(t, first.ID, got.ID)
		assert.Equal(t, "keep-me", got.BackgroundImage)
	})

	t.Run("agency", func(t *testing.T) {
		dropCollection(t, ctx, testDB, COLLECTION_AGENCY)

		require.NoError(t, service.GenerateAgencyFromLL2([]int64{101, 101}))

		var first models.Agency
		require.NoError(t, testDB.Collection(COLLECTION_AGENCY).FindOne(ctx, bson.M{"external_id": int64(101)}).Decode(&first))

		_, err := testDB.Collection(COLLECTION_AGENCY).UpdateOne(ctx, bson.M{"external_id": int64(101)}, bson.M{"$set": bson.M{"thumb_image": "keep-me"}})
		require.NoError(t, err)

		require.NoError(t, service.GenerateAgencyFromLL2([]int64{101}))

		count, err := testDB.Collection(COLLECTION_AGENCY).CountDocuments(ctx, bson.M{"external_id": int64(101)})
		require.NoError(t, err)
		assert.EqualValues(t, 1, count)

		var got models.Agency
		require.NoError(t, testDB.Collection(COLLECTION_AGENCY).FindOne(ctx, bson.M{"external_id": int64(101)}).Decode(&got))
		assert.Equal(t, first.ID, got.ID)
		assert.Equal(t, "keep-me", got.ThumbImage)
	})

	t.Run("rocket", func(t *testing.T) {
		dropCollection(t, ctx, testDB, COLLECTION_ROCKET)

		require.NoError(t, service.GenerateRocketsFromLL2([]int64{202, 202}))

		var first models.Rocket
		require.NoError(t, testDB.Collection(COLLECTION_ROCKET).FindOne(ctx, bson.M{"external_id": int64(202)}).Decode(&first))

		_, err := testDB.Collection(COLLECTION_ROCKET).UpdateOne(ctx, bson.M{"external_id": int64(202)}, bson.M{"$set": bson.M{"thumb_image": "keep-me"}})
		require.NoError(t, err)

		require.NoError(t, service.GenerateRocketsFromLL2([]int64{202}))

		count, err := testDB.Collection(COLLECTION_ROCKET).CountDocuments(ctx, bson.M{"external_id": int64(202)})
		require.NoError(t, err)
		assert.EqualValues(t, 1, count)

		var got models.Rocket
		require.NoError(t, testDB.Collection(COLLECTION_ROCKET).FindOne(ctx, bson.M{"external_id": int64(202)}).Decode(&got))
		assert.Equal(t, first.ID, got.ID)
		assert.Equal(t, "keep-me", got.ThumbImage)
	})

	t.Run("launch_base", func(t *testing.T) {
		dropCollection(t, ctx, testDB, COLLECTION_LAUNCH_BASE)

		require.NoError(t, service.GenerateLaunchBaseFromLL2([]int64{303, 303}))

		var first models.LaunchBase
		require.NoError(t, testDB.Collection(COLLECTION_LAUNCH_BASE).FindOne(ctx, bson.M{"external_id": int64(303)}).Decode(&first))

		require.NoError(t, service.GenerateLaunchBaseFromLL2([]int64{303}))

		count, err := testDB.Collection(COLLECTION_LAUNCH_BASE).CountDocuments(ctx, bson.M{"external_id": int64(303)})
		require.NoError(t, err)
		assert.EqualValues(t, 1, count)

		var got models.LaunchBase
		require.NoError(t, testDB.Collection(COLLECTION_LAUNCH_BASE).FindOne(ctx, bson.M{"external_id": int64(303)}).Decode(&got))
		assert.Equal(t, first.ID, got.ID)
	})
}

func dropCollection(t *testing.T, ctx context.Context, testDB *db.MongoDB, name string) {
	t.Helper()
	err := testDB.Collection(name).Drop(ctx)
	require.NoError(t, err)
}
