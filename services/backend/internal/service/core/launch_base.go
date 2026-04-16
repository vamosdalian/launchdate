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
	return s.ensureInt64ExternalIDs(COLLECTION_LAUNCH_BASE, ll2ids)
}

type LaunchBaseQuery struct {
	Limit         int
	Offset        int
	Search        string
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
	filters := make([]bson.M, 0, 4)
	if searchClause := buildTextSearchClause(q.Search, "name", "country.name", "timezone_name", "celestial_body.name"); len(searchClause) > 0 {
		filters = append(filters, searchClause)
	}
	if name := strings.TrimSpace(q.Name); name != "" {
		filters = append(filters, bson.M{"name": bson.M{"$regex": name, "$options": "i"}})
	}
	if body := strings.TrimSpace(q.CelestialBody); body != "" {
		filters = append(filters, bson.M{"celestial_body.name": bson.M{"$regex": body, "$options": "i"}})
	}
	if country := strings.TrimSpace(q.Country); country != "" {
		filters = append(filters, bson.M{"country.name": bson.M{"$regex": country, "$options": "i"}})
	}
	return combineFilters(filters...)
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

func (m *MainService) GetPublicLaunchBases(page int, search string) (models.PublicLaunchBasePage, error) {
	return loadPublicCached(m.publicCache, m.publicListCacheOptions(publicDomainLaunchBase, page, search, nil), func() (models.PublicLaunchBasePage, error) {
		const pageSize = 20

		launchBaseList, err := m.GetLaunchBases(LaunchBaseQuery{
			Limit:     pageSize,
			Offset:    page * pageSize,
			Search:    search,
			SortBy:    "name",
			SortOrder: 1,
		})
		if err != nil {
			return models.PublicLaunchBasePage{}, err
		}

		bases := make([]models.PublicLaunchBaseListItem, 0, len(launchBaseList.Launches))
		for _, launchBase := range launchBaseList.Launches {
			if launchBase.ID == 0 {
				continue
			}
			bases = append(bases, buildPublicLaunchBaseListItem(models.LaunchBase{ID: launchBase.ID}, launchBase.Data))
		}

		return models.PublicLaunchBasePage{
			Count:       launchBaseList.Count,
			LaunchBases: bases,
		}, nil
	})
}

func (m *MainService) GetPublicLaunchBase(id int64) (models.PublicLaunchBaseView, error) {
	return loadPublicCached(m.publicCache, m.publicDetailCacheOptions(publicDomainLaunchBase, id), func() (models.PublicLaunchBaseView, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var base models.LaunchBase
		err := m.mc.Collection(COLLECTION_LAUNCH_BASE).FindOne(ctx, bson.M{"id": id}).Decode(&base)
		if err != nil {
			return models.PublicLaunchBaseView{}, err
		}

		var doc models.LL2LocationSerializerWithPads
		if err := m.mc.Collection(COLLECTION_LL2_LOCATION).FindOne(ctx, bson.M{"id": int(base.ExternalID)}).Decode(&doc); err != nil {
			return models.PublicLaunchBaseView{}, err
		}

		baseFilter := bson.M{"pad.location.id": int(base.ExternalID)}
		launchCursor, err := m.mc.Collection(COLLECTION_LL2_LAUNCH).Find(
			ctx,
			baseFilter,
			options.Find().SetSort(bson.D{{Key: "net", Value: -1}}).SetLimit(6),
		)
		if err != nil {
			return models.PublicLaunchBaseView{}, err
		}
		defer launchCursor.Close(ctx)

		ll2Launches := make([]models.LL2LaunchDetailed, 0)
		launchExternalIDs := make([]string, 0)
		rocketExternalIDs := make([]int64, 0)
		agencyExternalIDs := make([]int64, 0)
		for launchCursor.Next(ctx) {
			var ll2Launch models.LL2LaunchDetailed
			if err := launchCursor.Decode(&ll2Launch); err != nil {
				return models.PublicLaunchBaseView{}, err
			}
			ll2Launches = append(ll2Launches, ll2Launch)
			launchExternalIDs = append(launchExternalIDs, ll2Launch.ID)
			if ll2Launch.Rocket.Configuration.ID != 0 {
				rocketExternalIDs = append(rocketExternalIDs, int64(ll2Launch.Rocket.Configuration.ID))
			}
			if ll2Launch.LaunchServiceProvider.ID != 0 {
				agencyExternalIDs = append(agencyExternalIDs, int64(ll2Launch.LaunchServiceProvider.ID))
			}
		}
		if err := launchCursor.Err(); err != nil {
			return models.PublicLaunchBaseView{}, err
		}

		launchMap, err := m.loadLaunchMapByExternalIDs(ctx, launchExternalIDs)
		if err != nil {
			return models.PublicLaunchBaseView{}, err
		}
		rocketMap, err := m.loadRocketMapByExternalIDs(ctx, rocketExternalIDs)
		if err != nil {
			return models.PublicLaunchBaseView{}, err
		}
		launcherMap, err := m.loadLauncherConfigMapByExternalIDs(ctx, rocketExternalIDs)
		if err != nil {
			return models.PublicLaunchBaseView{}, err
		}
		agencyMap, err := m.loadAgencyMapByExternalIDs(ctx, agencyExternalIDs)
		if err != nil {
			return models.PublicLaunchBaseView{}, err
		}
		agencyDocMap, err := m.loadLL2AgencyMapByExternalIDs(ctx, agencyExternalIDs)
		if err != nil {
			return models.PublicLaunchBaseView{}, err
		}
		baseMap := map[int64]models.LaunchBase{base.ExternalID: base}

		launches := make([]models.PublicLaunchSummary, 0, len(ll2Launches))
		for _, ll2Launch := range ll2Launches {
			launchSummary, include := m.buildPublicLaunchSummary(launchMap[ll2Launch.ID], ll2Launch, rocketMap, launcherMap, agencyMap, agencyDocMap, baseMap)
			if !include {
				continue
			}
			launches = append(launches, launchSummary)
		}

		launchCount64, err := m.mc.Collection(COLLECTION_LL2_LAUNCH).CountDocuments(ctx, baseFilter)
		if err != nil {
			return models.PublicLaunchBaseView{}, err
		}
		upcomingCount64, err := m.mc.Collection(COLLECTION_LL2_LAUNCH).CountDocuments(ctx, bson.M{
			"pad.location.id": int(base.ExternalID),
			"net":             bson.M{"$gte": time.Now().UTC().Format(time.RFC3339)},
		})
		if err != nil {
			return models.PublicLaunchBaseView{}, err
		}
		successfulCount64, err := m.mc.Collection(COLLECTION_LL2_LAUNCH).CountDocuments(ctx, bson.M{
			"pad.location.id": int(base.ExternalID),
			"status.id":       3,
		})
		if err != nil {
			return models.PublicLaunchBaseView{}, err
		}
		failedCount64, err := m.mc.Collection(COLLECTION_LL2_LAUNCH).CountDocuments(ctx, bson.M{
			"pad.location.id": int(base.ExternalID),
			"status.id":       4,
		})
		if err != nil {
			return models.PublicLaunchBaseView{}, err
		}

		successRate := 0
		resolvedCount := successfulCount64 + failedCount64
		if resolvedCount > 0 {
			successRate = int((float64(successfulCount64) * 100 / float64(resolvedCount)) + 0.5)
		}

		basic := buildPublicLaunchBaseListItem(base, doc)
		return models.PublicLaunchBaseView{
			ID:          basic.ID,
			Name:        basic.Name,
			Location:    basic.Location,
			Country:     basic.Country,
			Description: basic.Description,
			ImageURL:    basic.ImageURL,
			Latitude:    basic.Latitude,
			Longitude:   basic.Longitude,
			Launches:    launches,
			Stats: models.PublicLaunchBaseStats{
				LaunchCount:         int(launchCount64),
				UpcomingLaunchCount: int(upcomingCount64),
				SuccessfulLaunches:  int(successfulCount64),
				FailedLaunches:      int(failedCount64),
				SuccessRate:         successRate,
			},
		}, nil
	})
}
