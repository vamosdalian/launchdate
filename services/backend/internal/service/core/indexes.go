package core

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *MainService) EnsurePublicIndexes() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	indexSpecs := []struct {
		collection string
		models     []mongo.IndexModel
	}{
		{
			collection: COLLECTION_LAUNCH,
			models: []mongo.IndexModel{
				indexModel("id_unique", bson.D{{Key: "id", Value: 1}}, true),
				indexModel("external_id_unique", bson.D{{Key: "external_id", Value: 1}}, true),
			},
		},
		{
			collection: COLLECTION_ROCKET,
			models: []mongo.IndexModel{
				indexModel("id_unique", bson.D{{Key: "id", Value: 1}}, true),
				indexModel("external_id_unique", bson.D{{Key: "external_id", Value: 1}}, true),
			},
		},
		{
			collection: COLLECTION_AGENCY,
			models: []mongo.IndexModel{
				indexModel("id_unique", bson.D{{Key: "id", Value: 1}}, true),
				indexModel("external_id_unique", bson.D{{Key: "external_id", Value: 1}}, true),
				indexModel("show_on_home", bson.D{{Key: "show_on_home", Value: 1}}, false),
			},
		},
		{
			collection: COLLECTION_LAUNCH_BASE,
			models: []mongo.IndexModel{
				indexModel("id_unique", bson.D{{Key: "id", Value: 1}}, true),
				indexModel("external_id_unique", bson.D{{Key: "external_id", Value: 1}}, true),
			},
		},
		{
			collection: COLLECTION_LL2_LAUNCH,
			models: []mongo.IndexModel{
				indexModel("id_unique", bson.D{{Key: "id", Value: 1}}, true),
				indexModel("net", bson.D{{Key: "net", Value: 1}}, false),
				indexModel("rocket_configuration_id", bson.D{{Key: "rocket.configuration.id", Value: 1}}, false),
				indexModel("pad_location_id", bson.D{{Key: "pad.location.id", Value: 1}}, false),
				indexModel("launch_service_provider_id", bson.D{{Key: "launch_service_provider.id", Value: 1}}, false),
			},
		},
		{
			collection: COLLECTION_LL2_LAUNCHER,
			models: []mongo.IndexModel{
				indexModel("id_unique", bson.D{{Key: "id", Value: 1}}, true),
				indexModel("manufacturer_id", bson.D{{Key: "manufacturer.id", Value: 1}}, false),
			},
		},
		{
			collection: COLLECTION_LL2_LOCATION,
			models: []mongo.IndexModel{
				indexModel("id_unique", bson.D{{Key: "id", Value: 1}}, true),
			},
		},
	}

	for _, spec := range indexSpecs {
		if len(spec.models) == 0 {
			continue
		}
		if _, err := s.mc.Collection(spec.collection).Indexes().CreateMany(ctx, spec.models); err != nil {
			return err
		}
	}

	return nil
}

func indexModel(name string, keys bson.D, unique bool) mongo.IndexModel {
	opts := options.Index().SetName(name)
	if unique {
		opts.SetUnique(true)
	}

	return mongo.IndexModel{
		Keys:    keys,
		Options: opts,
	}
}
