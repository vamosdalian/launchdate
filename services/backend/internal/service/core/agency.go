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

const COLLECTION_AGENCY = "agency"
const COLLECTION_LL2_AGENCY = "ll2_agency"

func (s *MainService) GenerateAgencyFromLL2(ll2ids []int64) error {
	colleciton := s.mc.Collection(COLLECTION_AGENCY)

	var docs []any
	for _, ll2id := range ll2ids {
		docs = append(docs, models.Agency{
			ID:         s.sn.Generate().Int64(),
			ExternalID: ll2id,
		})
	}
	_, err := colleciton.InsertMany(context.Background(), docs)

	return err
}

type AgencyQuery struct {
	Limit     int
	Offset    int
	Name      string
	Type      string
	Country   string
	SortBy    string
	SortOrder int
}

func (q AgencyQuery) sortFieldAndOrder() (string, int) {
	switch strings.ToLower(strings.TrimSpace(q.SortBy)) {
	case "founding_year":
		if q.SortOrder < 0 {
			return "founding_year", -1
		}
		return "founding_year", 1
	default:
		if q.SortOrder < 0 {
			return "name", -1
		}
		return "name", 1
	}
}

func buildAgencyFilter(q AgencyQuery) bson.M {
	filter := bson.M{}
	if name := strings.TrimSpace(q.Name); name != "" {
		filter["name"] = bson.M{"$regex": name, "$options": "i"}
	}
	if agencyType := strings.TrimSpace(q.Type); agencyType != "" {
		filter["type.name"] = bson.M{"$regex": agencyType, "$options": "i"}
	}
	if country := strings.TrimSpace(q.Country); country != "" {
		filter["country.name"] = bson.M{"$regex": country, "$options": "i"}
	}
	return filter
}

func (s *MainService) GetAgencies(q AgencyQuery) (models.AgencyList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ll2Collection := s.mc.Collection(COLLECTION_LL2_AGENCY)
	filter := buildAgencyFilter(q)

	total, err := ll2Collection.CountDocuments(ctx, filter)
	if err != nil {
		return models.AgencyList{}, err
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
		return models.AgencyList{}, err
	}
	defer cursor.Close(ctx)

	ll2Docs := make([]models.LL2AgencyDetailed, 0)
	externalIDs := make([]int64, 0)

	for cursor.Next(ctx) {
		var doc models.LL2AgencyDetailed
		if err := cursor.Decode(&doc); err != nil {
			return models.AgencyList{}, err
		}
		ll2Docs = append(ll2Docs, doc)
		if doc.ID != 0 {
			externalIDs = append(externalIDs, int64(doc.ID))
		}
	}

	if err := cursor.Err(); err != nil {
		return models.AgencyList{}, err
	}

	agencyMap := make(map[int64]models.Agency)
	if len(externalIDs) > 0 {
		agencyCursor, err := s.mc.Collection(COLLECTION_AGENCY).
			Find(ctx, bson.M{"external_id": bson.M{"$in": externalIDs}})
		if err != nil {
			return models.AgencyList{}, err
		}
		defer agencyCursor.Close(ctx)

		for agencyCursor.Next(ctx) {
			var agency models.Agency
			if err := agencyCursor.Decode(&agency); err != nil {
				return models.AgencyList{}, err
			}
			agencyMap[agency.ExternalID] = agency
		}

		if err := agencyCursor.Err(); err != nil {
			return models.AgencyList{}, err
		}
	}

	serializers := make([]models.AgencySerializer, 0, len(ll2Docs))
	for _, doc := range ll2Docs {
		agency := agencyMap[int64(doc.ID)]
		serializers = append(serializers, models.AgencySerializer{
			Agency: agency,
			Data:   doc.LL2AgencyNormal,
		})
	}

	return models.AgencyList{
		Count:    int(total),
		Agencies: serializers,
	}, nil
}

func (s *MainService) GetAgency(id int64) (models.AgencySerializer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := s.mc.Collection(COLLECTION_AGENCY)

	var agency models.Agency
	if err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&agency); err != nil {
		return models.AgencySerializer{}, err
	}

	var ll2Data models.LL2AgencyNormal
	if agency.ExternalID != 0 {
		var doc models.LL2AgencyDetailed
		err := s.mc.Collection(COLLECTION_LL2_AGENCY).FindOne(ctx, bson.M{"id": int(agency.ExternalID)}).Decode(&doc)
		if err != nil && err != mongo.ErrNoDocuments {
			return models.AgencySerializer{}, err
		}
		if err == nil {
			ll2Data = doc.LL2AgencyNormal
		}
	}

	return models.AgencySerializer{
		Agency: agency,
		Data:   ll2Data,
	}, nil
}

func (m *MainService) UpdateAgency(a *models.Agency) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := m.mc.Collection(COLLECTION_AGENCY)

	update := bson.M{
		"$set": bson.M{
			"thumb_image": a.ThumbImage,
			"images":      a.Images,
			"social_url":  a.SocialUrl,
		},
	}

	result, err := collection.UpdateOne(ctx, bson.M{"id": a.ID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
