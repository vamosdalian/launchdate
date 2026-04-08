package image

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"github.com/vamosdalian/launchdate-backend/internal/db"
	"github.com/vamosdalian/launchdate-backend/internal/models"
	"github.com/vamosdalian/launchdate-backend/internal/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupMongoContainer(t *testing.T) (*db.MongoDB, func()) {
	ctx := context.Background()
	mongodbContainer, err := mongodb.Run(ctx, "mongo:6")
	if err != nil {
		t.Fatalf("failed to start container: %s", err)
	}

	endpoint, err := mongodbContainer.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("failed to get connection string: %s", err)
	}

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(endpoint))
	if err != nil {
		t.Fatalf("failed to connect to mongo: %s", err)
	}

	mongoDB := &db.MongoDB{
		Client:   mongoClient,
		Database: mongoClient.Database("testdb"),
	}

	cleanup := func() {
		if err := mongodbContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}

	return mongoDB, cleanup
}

func TestUploadImage(t *testing.T) {
	// Setup MongoDB
	mongoDB, cleanup := setupMongoContainer(t)
	defer cleanup()

	// Setup S3 Mock
	s3Server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		assert.Equal(t, "PUT", r.Method)
		// With PathStyle, the bucket name is part of the path
		assert.Contains(t, r.URL.Path, "/test-bucket/")
		assert.Equal(t, "image/png", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
	}))
	defer s3Server.Close()

	// Create S3 Client
	s3Client, err := util.CreateS3Client("test-id", "test-key", "auto", s3Server.URL)
	require.NoError(t, err)

	// Create ImageService
	svc := NewImageService(s3Client, mongoDB, "test-bucket", "")

	// Create a dummy image
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	// Fill with some color
	for x := 0; x < 100; x++ {
		for y := 0; y < 100; y++ {
			img.Set(x, y, color.RGBA{255, 0, 0, 255})
		}
	}
	var buf bytes.Buffer
	err = png.Encode(&buf, img)
	require.NoError(t, err)

	// Call UploadImage
	ctx := context.Background()
	key, err := svc.UploadImage(ctx, &buf, "test.png")
	require.NoError(t, err)
	assert.NotEmpty(t, key)
	assert.Contains(t, key, ".png")
	assert.NotContains(t, key, "images/")

	// Verify MongoDB
	var result bson.M
	err = mongoDB.Collection(imageMetaCollectionName).FindOne(ctx, bson.M{"key": key}).Decode(&result)
	require.NoError(t, err)
	assert.Equal(t, "test.png", result["name"])
	assert.NotEmpty(t, result["id"])
	assert.NotEmpty(t, result["upload_time"])
	assert.NotEmpty(t, result["size"])
	assert.Equal(t, int32(100), result["width"])
	assert.Equal(t, int32(100), result["height"])
}

func TestListImages(t *testing.T) {
	mongoDB, cleanup := setupMongoContainer(t)
	defer cleanup()

	domain := "https://cdn.example.com"
	svc := NewImageService(nil, mongoDB, "test-bucket", domain)
	ctx := context.Background()
	ploadTime, _ := time.Parse(time.RFC3339, "2023-01-01T00:00:00Z")
	// Insert some dummy data
	for i := 0; i < 15; i++ {
		_, err := mongoDB.Collection(imageMetaCollectionName).InsertOne(ctx, models.Image{
			ID:         int64(i),
			Key:        "test-key.webp",
			Name:       "test.png",
			UploadTime: ploadTime,
			ThumbImages: []models.ThumbImage{
				{
					Key: "test-thumb.webp",
				},
			},
		})
		require.NoError(t, err)
	}

	// Test pagination
	list, err := svc.ListImages(10, 0)
	require.NoError(t, err)
	assert.Equal(t, 15, list.Count)
	assert.Len(t, list.Images, 10)
	assert.Equal(t, domain+"/test-key.webp", list.Images[0].URL)
	assert.NotEmpty(t, list.Images[0].ThumbImages)
	assert.Equal(t, domain+"/test-thumb.webp", list.Images[0].ThumbImages[0].URL)

	list, err = svc.ListImages(10, 10)
	require.NoError(t, err)
	assert.Equal(t, 15, list.Count)
	assert.Len(t, list.Images, 5)
}

func TestDeleteImage(t *testing.T) {
	// Setup MongoDB
	mongoDB, cleanup := setupMongoContainer(t)
	defer cleanup()

	// Setup S3 Mock
	s3Server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		assert.Equal(t, "DELETE", r.Method)
		// With PathStyle, the bucket name is part of the path
		assert.Contains(t, r.URL.Path, "/test-bucket/")
		w.WriteHeader(http.StatusNoContent)
	}))
	defer s3Server.Close()

	// Create S3 Client
	s3Client, err := util.CreateS3Client("test-id", "test-key", "auto", s3Server.URL)
	require.NoError(t, err)

	// Create ImageService
	svc := NewImageService(s3Client, mongoDB, "test-bucket", "")
	ctx := context.Background()

	// Insert a dummy image metadata
	key := "test-image.webp"
	_, err = mongoDB.Collection("image_meta").InsertOne(ctx, models.Image{
		Key:  key,
		Name: "test.png",
	})
	require.NoError(t, err)

	// Call DeleteImage
	err = svc.DeleteImage(key)
	require.NoError(t, err)

	// Verify MongoDB deletion
	count, err := mongoDB.Collection("image_meta").CountDocuments(ctx, bson.M{"key": key})
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

func TestGenerateThumb(t *testing.T) {
	// Setup MongoDB
	mongoDB, cleanup := setupMongoContainer(t)
	defer cleanup()

	// Create a dummy image for S3 GetObject
	img := image.NewRGBA(image.Rect(0, 0, 200, 200))
	for x := 0; x < 200; x++ {
		for y := 0; y < 200; y++ {
			img.Set(x, y, color.RGBA{0, 255, 0, 255})
		}
	}
	var imgBuf bytes.Buffer
	err := png.Encode(&imgBuf, img)
	require.NoError(t, err)
	imgBytes := imgBuf.Bytes()

	// Setup S3 Mock
	s3Server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			// Serve the dummy image
			w.Header().Set("Content-Type", "image/png")
			w.Write(imgBytes)
			return
		}
		if r.Method == "PUT" {
			// Verify upload
			assert.Equal(t, "image/webp", r.Header.Get("Content-Type"))
			// Check key format
			assert.Contains(t, r.URL.Path, "@100x100.webp")

			// Read body to verify it's not empty
			body, err := io.ReadAll(r.Body)
			assert.NoError(t, err)
			assert.NotEmpty(t, body)

			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer s3Server.Close()

	// Create S3 Client
	s3Client, err := util.CreateS3Client("test-id", "test-key", "auto", s3Server.URL)
	require.NoError(t, err)

	// Create ImageService
	svc := NewImageService(s3Client, mongoDB, "test-bucket", "")

	// Insert initial metadata into MongoDB
	ctx := context.Background()
	originalKey := "12345.webp"
	_, err = mongoDB.Collection(imageMetaCollectionName).InsertOne(ctx, models.Image{
		ID:          12345,
		Key:         originalKey,
		Name:        "test.png",
		Width:       200,
		Height:      200,
		Size:        int64(len(imgBytes)),
		ContentType: "image/png",
		UploadTime:  time.Now(),
		ThumbImages: []models.ThumbImage{},
	})
	require.NoError(t, err)

	// Call GenerateThumb
	err = svc.GenerateThumb(ctx, 12345, 100, 100)
	require.NoError(t, err)

	// Verify MongoDB update
	var result models.Image
	err = mongoDB.Collection(imageMetaCollectionName).FindOne(ctx, bson.M{"id": 12345}).Decode(&result)
	require.NoError(t, err)

	require.Len(t, result.ThumbImages, 1)
	thumb := result.ThumbImages[0]
	assert.Equal(t, fmt.Sprintf("%d@100x100.webp", 12345), thumb.Key)
	assert.Equal(t, 100, thumb.Width)
	assert.Equal(t, 100, thumb.Height)
	assert.Equal(t, "image/webp", thumb.ContentType)
}
