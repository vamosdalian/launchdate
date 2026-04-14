package core

import (
	"context"

	"github.com/bwmarrin/snowflake"
	"github.com/vamosdalian/launchdate-backend/internal/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MainService struct {
	mc *db.MongoDB
	sn *snowflake.Node
}

func NewMainService(mc *db.MongoDB) *MainService {
	node, _ := snowflake.NewNode(0)
	return &MainService{
		mc: mc,
		sn: node,
	}
}

func (s *MainService) ensureStringExternalIDs(collectionName string, externalIDs []string) error {
	collection := s.mc.Collection(collectionName)
	seen := make(map[string]struct{}, len(externalIDs))
	models := make([]mongo.WriteModel, 0, len(externalIDs))

	for _, externalID := range externalIDs {
		if externalID == "" {
			continue
		}
		if _, ok := seen[externalID]; ok {
			continue
		}
		seen[externalID] = struct{}{}

		models = append(models, mongo.NewUpdateOneModel().
			SetFilter(bson.M{"external_id": externalID}).
			SetUpdate(bson.M{"$setOnInsert": bson.M{
				"id":          s.sn.Generate().Int64(),
				"external_id": externalID,
			}}).
			SetUpsert(true))
	}

	if len(models) == 0 {
		return nil
	}

	_, err := collection.BulkWrite(context.Background(), models, options.BulkWrite().SetOrdered(false))
	return err
}

func (s *MainService) ensureInt64ExternalIDs(collectionName string, externalIDs []int64) error {
	collection := s.mc.Collection(collectionName)
	seen := make(map[int64]struct{}, len(externalIDs))
	models := make([]mongo.WriteModel, 0, len(externalIDs))

	for _, externalID := range externalIDs {
		if externalID == 0 {
			continue
		}
		if _, ok := seen[externalID]; ok {
			continue
		}
		seen[externalID] = struct{}{}

		models = append(models, mongo.NewUpdateOneModel().
			SetFilter(bson.M{"external_id": externalID}).
			SetUpdate(bson.M{"$setOnInsert": bson.M{
				"id":          s.sn.Generate().Int64(),
				"external_id": externalID,
			}}).
			SetUpsert(true))
	}

	if len(models) == 0 {
		return nil
	}

	_, err := collection.BulkWrite(context.Background(), models, options.BulkWrite().SetOrdered(false))
	return err
}
