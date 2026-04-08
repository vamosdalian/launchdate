package image

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"time"

	"github.com/HugoSmits86/nativewebp"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/bwmarrin/snowflake"
	"github.com/disintegration/imaging"
	"github.com/vamosdalian/launchdate-backend/internal/db"
	"github.com/vamosdalian/launchdate-backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	imageMetaCollectionName = "image_meta"
)

type ImageService struct {
	bucketName  string
	imageDoamin string
	os          *s3.Client
	db          *db.MongoDB
	sn          *snowflake.Node
}

func NewImageService(os *s3.Client, db *db.MongoDB, bucketName string, domain string) *ImageService {
	node, _ := snowflake.NewNode(0)
	return &ImageService{
		bucketName:  bucketName,
		imageDoamin: domain,
		os:          os,
		db:          db,
		sn:          node,
	}
}

func (i *ImageService) UploadImage(ctx context.Context, file io.Reader, filename string) (string, error) {
	data, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	cfg, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("failed to decode image config: %w", err)
	}

	id := i.sn.Generate().Int64()
	ext := "." + format
	if format == "jpeg" {
		ext = ".jpg"
	}
	key := fmt.Sprintf("%d%s", id, ext)
	contentType := fmt.Sprintf("image/%s", format)

	// Upload to S3
	_, err = i.os.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(i.bucketName),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to s3: %w", err)
	}

	// Save metadata to MongoDB
	meta := models.Image{
		ID:          id,
		Key:         key,
		Name:        filename,
		Width:       cfg.Width,
		Height:      cfg.Height,
		Size:        int64(len(data)),
		ContentType: contentType,
		UploadTime:  time.Now(),
		ThumbImages: []models.ThumbImage{},
	}

	_, err = i.db.Collection(imageMetaCollectionName).InsertOne(ctx, meta)
	if err != nil {
		return "", fmt.Errorf("failed to save metadata: %w", err)
	}

	return key, nil
}

func (i *ImageService) ListImages(limit, offset int) (models.ImageList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := i.db.Collection(imageMetaCollectionName)

	// Count total documents
	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return models.ImageList{}, fmt.Errorf("failed to count images: %w", err)
	}

	// Find documents with pagination
	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.M{"upload_time": -1}) // Sort by upload time desc

	cursor, err := collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return models.ImageList{}, fmt.Errorf("failed to list images: %w", err)
	}
	defer cursor.Close(ctx)

	var images []models.Image
	if err = cursor.All(ctx, &images); err != nil {
		return models.ImageList{}, fmt.Errorf("failed to decode images: %w", err)
	}

	// Ensure images is not nil for JSON serialization
	if images == nil {
		images = []models.Image{}
	}

	// Populate URLs
	for idx := range images {
		images[idx].URL = fmt.Sprintf("%s/%s", i.imageDoamin, images[idx].Key)
		for tIdx := range images[idx].ThumbImages {
			images[idx].ThumbImages[tIdx].URL = fmt.Sprintf("%s/%s", i.imageDoamin, images[idx].ThumbImages[tIdx].Key)
		}
	}

	return models.ImageList{
		Count:  int(count),
		Images: images,
	}, nil
}

func (i *ImageService) DeleteImage(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get metadata to find thumbnails
	var meta models.Image
	if err := i.db.Collection(imageMetaCollectionName).FindOne(ctx, bson.M{"key": key}).Decode(&meta); err != nil {
		return fmt.Errorf("failed to find image metadata: %w", err)
	}

	// Delete from S3
	_, err := i.os.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(i.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete image from s3: %w", err)
	}

	// Delete thumbnails from S3
	for _, thumb := range meta.ThumbImages {
		_, err = i.os.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(i.bucketName),
			Key:    aws.String(thumb.Key),
		})
		if err != nil {
			return fmt.Errorf("failed to delete thumbnail from s3: %w", err)
		}
	}

	// Delete metadata from MongoDB
	_, err = i.db.Collection(imageMetaCollectionName).DeleteOne(ctx, bson.M{"key": key})
	if err != nil {
		return fmt.Errorf("failed to delete image metadata: %w", err)
	}

	return nil
}

func (i *ImageService) GenerateThumb(ctx context.Context, id int64, width, height int) error {
	// 1. Get image metadata
	var meta models.Image
	err := i.db.Collection(imageMetaCollectionName).FindOne(ctx, bson.M{"id": id}).Decode(&meta)
	if err != nil {
		return fmt.Errorf("failed to find image metadata: %w", err)
	}

	// Check if thumbnail already exists
	for _, thumb := range meta.ThumbImages {
		if thumb.Width == width && thumb.Height == height {
			return nil
		}
	}

	// 2. Download original image from S3
	resp, err := i.os.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(i.bucketName),
		Key:    aws.String(meta.Key),
	})
	if err != nil {
		return fmt.Errorf("failed to download image from s3: %w", err)
	}
	defer resp.Body.Close()

	// 3. Decode image
	img, err := imaging.Decode(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to decode image: %w", err)
	}

	// 4. Resize image
	resizedImg := imaging.Resize(img, width, height, imaging.Lanczos)

	// 5. Encode to WebP
	var buf bytes.Buffer
	if err := nativewebp.Encode(&buf, resizedImg, nil); err != nil {
		return fmt.Errorf("failed to encode thumbnail to webp: %w", err)
	}

	// 6. Upload thumbnail to S3
	thumbKey := fmt.Sprintf("%d@%dx%d.webp", id, width, height)
	_, err = i.os.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(i.bucketName),
		Key:         aws.String(thumbKey),
		Body:        bytes.NewReader(buf.Bytes()),
		ContentType: aws.String("image/webp"),
	})
	if err != nil {
		return fmt.Errorf("failed to upload thumbnail to s3: %w", err)
	}

	// 7. Update MongoDB
	thumbMeta := models.ThumbImage{
		Key:         thumbKey,
		Width:       width,
		Height:      height,
		Size:        int64(buf.Len()),
		ContentType: "image/webp",
		UploadTime:  time.Now(),
	}

	_, err = i.db.Collection(imageMetaCollectionName).UpdateOne(ctx, bson.M{"id": id}, bson.M{
		"$push": bson.M{"thumb_images": thumbMeta},
	})
	if err != nil {
		return fmt.Errorf("failed to update metadata with thumbnail: %w", err)
	}

	return nil
}
