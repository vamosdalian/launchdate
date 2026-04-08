package core

import (
	"context"
	"strings"
	"time"

	"github.com/vamosdalian/launchdate-backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const COLLECTION_LAUNCH_BASE = "launch_base"
const COLLECTION_LL2_LOCATION = "ll2_location"

func (s *MainService) GenerateLaunchBaseFromLL2(ll2ids []int64) error {
	colleciton := s.mc.Collection(COLLECTION_LAUNCH_BASE)

	var docs []any
	for _, ll2id := range ll2ids {
		docs = append(docs, models.LaunchBase{
			ID:         s.sn.Generate().Int64(),
			ExternalID: ll2id,
		})
	}
	_, err := colleciton.InsertMany(context.Background(), docs)

	return err
}

type LaunchBaseQuery struct {
	Limit         int
	Offset        int
	Name          string
	CelestialBody string
	Country       string
	SortBy        string
	SortOrder     int
}

func (q LaunchBaseQuery) sortFieldAndOrder() (string, int) {
	switch strings.ToLower(strings.TrimSpace(q.SortBy)) {
	case "name":
		if q.SortOrder < 0 {
			return "name", -1
		}
		return "name", 1
	default:
		if q.SortOrder < 0 {
			return "id", -1
		}
		return "id", 1
	}
}

func buildLaunchBaseFilter(q LaunchBaseQuery) bson.M {
	filter := bson.M{}
	if name := strings.TrimSpace(q.Name); name != "" {
		filter["name"] = bson.M{"$regex": name, "$options": "i"}
	}
	if body := strings.TrimSpace(q.CelestialBody); body != "" {
		filter["celestial_body.name"] = bson.M{"$regex": body, "$options": "i"}
	}
	if country := strings.TrimSpace(q.Country); country != "" {
		filter["country.name"] = bson.M{"$regex": country, "$options": "i"}
	}
	return filter
}

func (m *MainService) LoadLaunchBasesFromLL2() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	ll2Locations := m.mc.Collection(COLLECTION_LL2_LOCATION)
	opts := options.Find().SetProjection(bson.M{"id": 1})
	cursor, err := ll2Locations.Find(ctx, bson.M{"id": bson.M{"$ne": nil}}, opts)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	launchBaseCollection := m.mc.Collection(COLLECTION_LAUNCH_BASE)
	seen := make(map[int64]struct{})

	for cursor.Next(ctx) {
		var document struct {
			ID int `bson:"id"`
		}
		if err := cursor.Decode(&document); err != nil {
			return err
		}

		externalID := int64(document.ID)
		if externalID == 0 {
			continue
		}
		if _, exists := seen[externalID]; exists {
			continue
		}
		seen[externalID] = struct{}{}

		filter := bson.M{"external_id": externalID}
		update := bson.M{
			"$setOnInsert": bson.M{
				"id":          m.sn.Generate().Int64(),
				"external_id": externalID,
			},
		}

		if _, err := launchBaseCollection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true)); err != nil {
			return err
		}
	}

	if err := cursor.Err(); err != nil {
		return err
	}

	return nil
}

func (m *MainService) GetLaunchBases(q LaunchBaseQuery) (models.LaunchBaseList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ll2Collection := m.mc.Collection(COLLECTION_LL2_LOCATION)
	filter := buildLaunchBaseFilter(q)

	total, err := ll2Collection.CountDocuments(ctx, filter)
	if err != nil {
		return models.LaunchBaseList{}, err
	}

	findOpts := options.Find()
	if q.Limit > 0 {
		findOpts.SetLimit(int64(q.Limit))
	}
	if q.Offset > 0 {
		findOpts.SetSkip(int64(q.Offset))
	}
	sortField, sortOrder := q.sortFieldAndOrder()
	findOpts.SetSort(bson.D{{Key: sortField, Value: sortOrder}})

	cursor, err := ll2Collection.Find(ctx, filter, findOpts)
	if err != nil {
		return models.LaunchBaseList{}, err
	}
	defer cursor.Close(ctx)

	ll2Docs := make([]models.LL2LocationSerializerWithPads, 0)
	locationIDs := make([]int, 0)

	for cursor.Next(ctx) {
		var doc models.LL2LocationSerializerWithPads
		if err := cursor.Decode(&doc); err != nil {
			return models.LaunchBaseList{}, err
		}
		ll2Docs = append(ll2Docs, doc)
		locationIDs = append(locationIDs, doc.ID)
	}

	if err := cursor.Err(); err != nil {
		return models.LaunchBaseList{}, err
	}

	if len(ll2Docs) == 0 {
		return models.LaunchBaseList{Count: int(total), Launches: []models.LaunchBaseSerializer{}}, nil
	}

	baseIDs := make(map[int64]int64, len(locationIDs))
	if len(locationIDs) > 0 {
		externalIDs := make([]int64, 0, len(locationIDs))
		for _, id := range locationIDs {
			externalIDs = append(externalIDs, int64(id))
		}
		baseCursor, err := m.mc.Collection(COLLECTION_LAUNCH_BASE).
			Find(ctx, bson.M{"external_id": bson.M{"$in": externalIDs}})
		if err != nil {
			return models.LaunchBaseList{}, err
		}
		defer baseCursor.Close(ctx)

		for baseCursor.Next(ctx) {
			var base models.LaunchBase
			if err := baseCursor.Decode(&base); err != nil {
				return models.LaunchBaseList{}, err
			}
			baseIDs[base.ExternalID] = base.ID
		}

		if err := baseCursor.Err(); err != nil {
			return models.LaunchBaseList{}, err
		}
	}

	serializers := make([]models.LaunchBaseSerializer, 0, len(ll2Docs))
	for _, doc := range ll2Docs {
		serializers = append(serializers, models.LaunchBaseSerializer{
			ID:   baseIDs[int64(doc.ID)],
			Data: doc,
		})
	}

	return models.LaunchBaseList{
		Count:    int(total),
		Launches: serializers,
	}, nil
}

func (m *MainService) GetLaunchBase(id int64) (models.LaunchBaseSerializer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := m.mc.Collection(COLLECTION_LAUNCH_BASE)

	var base models.LaunchBase
	if err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&base); err != nil {
		return models.LaunchBaseSerializer{}, err
	}

	var data models.LL2LocationSerializerWithPads
	if base.ExternalID != 0 {
		var doc models.LL2LocationSerializerWithPads
		err := m.mc.Collection(COLLECTION_LL2_LOCATION).FindOne(ctx, bson.M{"id": int(base.ExternalID)}).Decode(&doc)
		if err != nil && err != mongo.ErrNoDocuments {
			return models.LaunchBaseSerializer{}, err
		}
		if err == nil {
			data = doc
		}
	}

	return models.LaunchBaseSerializer{
		ID:   base.ID,
		Data: data,
	}, nil
}
