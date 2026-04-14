package core

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/vamosdalian/launchdate-backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const COLLECTION_PAGE_BACKGROUND = "page_background"

var ErrInvalidPageBackgroundKey = errors.New("invalid page background key")

func (s *MainService) EnsurePageBackgroundIndexes() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := s.mc.Collection(COLLECTION_PAGE_BACKGROUND).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "page_key", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("page_key_unique"),
	})

	return err
}

func (s *MainService) GetPageBackgrounds() ([]models.PageBackground, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := s.mc.Collection(COLLECTION_PAGE_BACKGROUND).Find(ctx, bson.M{
		"page_key": bson.M{"$in": pageBackgroundKeys()},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	records := make(map[string]models.PageBackground, len(models.PageBackgroundDefinitions))
	for cursor.Next(ctx) {
		var record models.PageBackground
		if err := cursor.Decode(&record); err != nil {
			return nil, err
		}
		records[record.PageKey] = decoratePageBackground(record)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	pageBackgrounds := make([]models.PageBackground, 0, len(models.PageBackgroundDefinitions))
	for _, definition := range models.PageBackgroundDefinitions {
		if record, ok := records[definition.Key]; ok {
			pageBackgrounds = append(pageBackgrounds, record)
			continue
		}

		pageBackgrounds = append(pageBackgrounds, decoratePageBackground(models.PageBackground{PageKey: definition.Key}))
	}

	return pageBackgrounds, nil
}

func (s *MainService) UpsertPageBackground(pageKey, backgroundImage string) (models.PageBackground, error) {
	pageKey = strings.TrimSpace(pageKey)
	if !models.IsValidPageBackgroundKey(pageKey) {
		return models.PageBackground{}, ErrInvalidPageBackgroundKey
	}

	backgroundImage = strings.TrimSpace(backgroundImage)
	if backgroundImage == "" {
		if err := s.deletePageBackground(pageKey); err != nil {
			return models.PageBackground{}, err
		}

		return decoratePageBackground(models.PageBackground{PageKey: pageKey}), nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	now := time.Now().UTC()
	update := bson.M{
		"$set": bson.M{
			"background_image": backgroundImage,
			"updated_at":       now,
		},
		"$setOnInsert": bson.M{
			"id":         s.sn.Generate().Int64(),
			"page_key":   pageKey,
			"created_at": now,
		},
	}

	_, err := s.mc.Collection(COLLECTION_PAGE_BACKGROUND).UpdateOne(
		ctx,
		bson.M{"page_key": pageKey},
		update,
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return models.PageBackground{}, err
	}

	var record models.PageBackground
	if err := s.mc.Collection(COLLECTION_PAGE_BACKGROUND).FindOne(ctx, bson.M{"page_key": pageKey}).Decode(&record); err != nil {
		return models.PageBackground{}, err
	}

	return decoratePageBackground(record), nil
}

func (s *MainService) deletePageBackground(pageKey string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := s.mc.Collection(COLLECTION_PAGE_BACKGROUND).DeleteOne(ctx, bson.M{"page_key": pageKey})
	return err
}

func pageBackgroundKeys() []string {
	keys := make([]string, 0, len(models.PageBackgroundDefinitions))
	for _, definition := range models.PageBackgroundDefinitions {
		keys = append(keys, definition.Key)
	}

	return keys
}

func decoratePageBackground(record models.PageBackground) models.PageBackground {
	record.DisplayName = models.PageBackgroundDisplayName(record.PageKey)
	record.Configured = strings.TrimSpace(record.BackgroundImage) != ""
	return record
}
