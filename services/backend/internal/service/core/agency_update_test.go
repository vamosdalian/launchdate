package core

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vamosdalian/launchdate-backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
)

func TestMainService_UpdateAgencyPersistsShowOnHome(t *testing.T) {
	service, testDB, cleanup := newTestMainService(t)
	defer cleanup()
	ctx := context.Background()

	clearCoreCollections(t, testDB, COLLECTION_AGENCY)

	_, err := testDB.Collection(COLLECTION_AGENCY).InsertMany(ctx, []interface{}{
		models.Agency{
			ID:         91001,
			ExternalID: 101,
			ThumbImage: "thumb-before.jpg",
			Images:     []string{"image-before.jpg"},
			SocialUrl:  []models.SocialUrl{{Name: "Website", URL: "https://before.example.com"}},
			ShowOnHome: false,
		},
		models.Agency{
			ID:         91002,
			ExternalID: 102,
			ThumbImage: "other-thumb.jpg",
			Images:     []string{"other-image.jpg"},
			SocialUrl:  []models.SocialUrl{{Name: "Website", URL: "https://other.example.com"}},
			ShowOnHome: false,
		},
	})
	require.NoError(t, err)

	err = service.UpdateAgency(&models.Agency{
		ID:         91001,
		ThumbImage: "thumb-after.jpg",
		Images:     []string{"image-after.jpg"},
		SocialUrl:  []models.SocialUrl{{Name: "Website", URL: "https://after.example.com"}},
		ShowOnHome: true,
	})
	require.NoError(t, err)

	var stored models.Agency
	err = testDB.Collection(COLLECTION_AGENCY).FindOne(ctx, bson.M{"id": int64(91001)}).Decode(&stored)
	require.NoError(t, err)

	assert.Equal(t, "thumb-after.jpg", stored.ThumbImage)
	assert.Equal(t, []string{"image-after.jpg"}, stored.Images)
	assert.Equal(t, []models.SocialUrl{{Name: "Website", URL: "https://after.example.com"}}, stored.SocialUrl)
	assert.True(t, stored.ShowOnHome)

	var untouched models.Agency
	err = testDB.Collection(COLLECTION_AGENCY).FindOne(ctx, bson.M{"id": int64(91002)}).Decode(&untouched)
	require.NoError(t, err)
	assert.Equal(t, "other-thumb.jpg", untouched.ThumbImage)
	assert.Equal(t, []string{"other-image.jpg"}, untouched.Images)
	assert.Equal(t, []models.SocialUrl{{Name: "Website", URL: "https://other.example.com"}}, untouched.SocialUrl)
	assert.False(t, untouched.ShowOnHome)
}
