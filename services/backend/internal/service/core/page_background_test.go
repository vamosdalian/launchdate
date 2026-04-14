package core

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"github.com/vamosdalian/launchdate-backend/internal/db"
	"github.com/vamosdalian/launchdate-backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
)

func TestMainService_PageBackgrounds(t *testing.T) {
	ctx := context.Background()

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

	testDB, cleanup, err := db.NewMongoDB(uri, "test_page_background_db")
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	service := NewMainService(testDB)
	testDB.Collection(COLLECTION_PAGE_BACKGROUND).Drop(ctx)

	assert.NoError(t, service.EnsurePageBackgroundIndexes())

	backgrounds, err := service.GetPageBackgrounds()
	assert.NoError(t, err)
	assert.Len(t, backgrounds, len(models.PageBackgroundDefinitions))
	assert.False(t, backgrounds[0].Configured)

	updated, err := service.UpsertPageBackground(models.PageBackgroundKeyHome, "https://cdn.example.com/home.jpg")
	assert.NoError(t, err)
	assert.True(t, updated.Configured)
	assert.Equal(t, models.PageBackgroundKeyHome, updated.PageKey)

	var stored models.PageBackground
	err = testDB.Collection(COLLECTION_PAGE_BACKGROUND).FindOne(ctx, bson.M{"page_key": models.PageBackgroundKeyHome}).Decode(&stored)
	assert.NoError(t, err)
	assert.Equal(t, "https://cdn.example.com/home.jpg", stored.BackgroundImage)

	cleared, err := service.UpsertPageBackground(models.PageBackgroundKeyHome, "   ")
	assert.NoError(t, err)
	assert.False(t, cleared.Configured)

	count, err := testDB.Collection(COLLECTION_PAGE_BACKGROUND).CountDocuments(ctx, bson.M{"page_key": models.PageBackgroundKeyHome})
	assert.NoError(t, err)
	assert.EqualValues(t, 0, count)

	_, err = service.UpsertPageBackground("invalid", "https://cdn.example.com/invalid.jpg")
	assert.ErrorIs(t, err, ErrInvalidPageBackgroundKey)
}
