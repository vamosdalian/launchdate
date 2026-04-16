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

const COLLECTION_ROCKET = "rocket"
const COLLECTION_LL2_LAUNCHER = "ll2_launcher"

type rocketAggregateDoc struct {
	RocketDoc                      models.Rocket `bson:"rocket_doc"`
	models.LL2LauncherConfigNormal `bson:",inline"`
}

type RocketQuery struct {
	Limit     int
	Offset    int
	Search    string
	FullName  string
	Name      string
	Variant   string
	SortBy    string
	SortOrder int
}

func (q RocketQuery) sortFieldAndOrder() (string, int) {
	switch strings.ToLower(strings.TrimSpace(q.SortBy)) {
	case "total_launch_count":
		if q.SortOrder < 0 {
			return "total_launch_count", -1
		}
		return "total_launch_count", 1
	case "full_name", "fullname":
		if q.SortOrder < 0 {
			return "full_name", -1
		}
		return "full_name", 1
	default:
		if q.SortOrder < 0 {
			return "rocket_doc.id", -1
		}
		return "rocket_doc.id", 1
	}
}

func buildRocketFilter(q RocketQuery) bson.M {
	filters := []bson.M{{"id": bson.M{"$ne": nil}}}
	if searchClause := buildTextSearchClause(q.Search, "full_name", "name", "variant"); len(searchClause) > 0 {
		filters = append(filters, searchClause)
	}
	if fullName := strings.TrimSpace(q.FullName); fullName != "" {
		filters = append(filters, bson.M{"full_name": bson.M{"$regex": fullName, "$options": "i"}})
	}
	if name := strings.TrimSpace(q.Name); name != "" {
		filters = append(filters, bson.M{"name": bson.M{"$regex": name, "$options": "i"}})
	}
	if variant := strings.TrimSpace(q.Variant); variant != "" {
		filters = append(filters, bson.M{"variant": bson.M{"$regex": variant, "$options": "i"}})
	}
	return combineFilters(filters...)
}

func rocketAggregationPipeline(filter bson.M) mongo.Pipeline {
	return mongo.Pipeline{
		bson.D{{Key: "$match", Value: filter}},
		bson.D{{Key: "$lookup", Value: bson.M{
			"from":         COLLECTION_ROCKET,
			"localField":   "id",
			"foreignField": "external_id",
			"as":           "rocket_doc",
		}}},
		bson.D{{Key: "$unwind", Value: bson.M{"path": "$rocket_doc", "preserveNullAndEmptyArrays": false}}},
	}
}

func (m *MainService) findRocketByExternalID(ctx context.Context, externalID int64) (models.Rocket, error) {
	var rocket models.Rocket
	err := m.mc.Collection(COLLECTION_ROCKET).FindOne(ctx, bson.M{"external_id": externalID}).Decode(&rocket)
	if err == nil {
		return rocket, nil
	}
	if err != mongo.ErrNoDocuments {
		return models.Rocket{}, err
	}

	return models.Rocket{}, mongo.ErrNoDocuments
}

func (m *MainService) findRocketByPublicID(ctx context.Context, id int64) (models.Rocket, error) {
	lookups := []struct {
		collectionName string
		filter         bson.M
	}{
		{collectionName: COLLECTION_ROCKET, filter: bson.M{"external_id": id}},
		{collectionName: COLLECTION_ROCKET, filter: bson.M{"id": id}},
	}

	for _, lookup := range lookups {
		var rocket models.Rocket
		err := m.mc.Collection(lookup.collectionName).FindOne(ctx, lookup.filter).Decode(&rocket)
		if err == nil {
			return rocket, nil
		}
		if err != mongo.ErrNoDocuments {
			return models.Rocket{}, err
		}
	}

	return models.Rocket{}, mongo.ErrNoDocuments
}

func (s *MainService) GenerateRocketsFromLL2(launchers []int64) error {
	return s.ensureInt64ExternalIDs(COLLECTION_ROCKET, launchers)
}

func (m *MainService) GetRockets(q RocketQuery) (models.RocketList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	launcherCollection := m.mc.Collection(COLLECTION_LL2_LAUNCHER)
	filter := buildRocketFilter(q)
	sortField, sortOrder := q.sortFieldAndOrder()

	basePipeline := rocketAggregationPipeline(filter)
	countPipeline := append(append(mongo.Pipeline{}, basePipeline...), bson.D{{Key: "$count", Value: "total"}})

	countCursor, err := launcherCollection.Aggregate(ctx, countPipeline)
	if err != nil {
		return models.RocketList{}, err
	}
	defer countCursor.Close(ctx)

	var total int64
	if countCursor.Next(ctx) {
		var countDoc struct {
			Total int64 `bson:"total"`
		}
		if err := countCursor.Decode(&countDoc); err != nil {
			return models.RocketList{}, err
		}
		total = countDoc.Total
	}
	if err := countCursor.Err(); err != nil {
		return models.RocketList{}, err
	}

	if total == 0 {
		return models.RocketList{Count: 0, Rockets: []models.RocketSerializer{}}, nil
	}

	dataPipeline := append(append(mongo.Pipeline{}, basePipeline...), bson.D{{Key: "$sort", Value: bson.D{{Key: sortField, Value: sortOrder}}}})
	if q.Offset > 0 {
		dataPipeline = append(dataPipeline, bson.D{{Key: "$skip", Value: q.Offset}})
	}
	if q.Limit > 0 {
		dataPipeline = append(dataPipeline, bson.D{{Key: "$limit", Value: q.Limit}})
	}

	dataCursor, err := launcherCollection.Aggregate(ctx, dataPipeline)
	if err != nil {
		return models.RocketList{}, err
	}
	defer dataCursor.Close(ctx)

	serializers := make([]models.RocketSerializer, 0)
	for dataCursor.Next(ctx) {
		var doc rocketAggregateDoc
		if err := dataCursor.Decode(&doc); err != nil {
			return models.RocketList{}, err
		}
		serializers = append(serializers, models.RocketSerializer{
			Rocket: doc.RocketDoc,
			Data:   doc.LL2LauncherConfigNormal,
		})
	}

	if err := dataCursor.Err(); err != nil {
		return models.RocketList{}, err
	}

	return models.RocketList{
		Count:   int(total),
		Rockets: serializers,
	}, nil
}

func (m *MainService) GetRocket(id int64) (models.RocketSerializer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := m.mc.Collection(COLLECTION_ROCKET)

	var rocket models.Rocket
	if err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&rocket); err != nil {
		return models.RocketSerializer{}, err
	}

	var data models.LL2LauncherConfigNormal
	if rocket.ExternalID != 0 {
		filter := bson.M{"id": int(rocket.ExternalID)}
		err := m.mc.Collection(COLLECTION_LL2_LAUNCHER).FindOne(ctx, filter).Decode(&data)
		if err != nil && err != mongo.ErrNoDocuments {
			return models.RocketSerializer{}, err
		}
	}

	return models.RocketSerializer{
		Rocket: rocket,
		Data:   data,
	}, nil
}

func (m *MainService) UpdateRocket(r *models.Rocket) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := m.mc.Collection(COLLECTION_ROCKET)

	update := bson.M{
		"$set": r,
	}

	result, err := collection.UpdateOne(ctx, bson.M{"id": r.ID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	m.bumpPublicCacheDomains(publicDomainRocket, publicDomainLaunch)

	return nil
}

func (m *MainService) GetPublicRockets(page int, search string) (models.PublicRocketPage, error) {
	return loadPublicCached(m.publicCache, m.publicListCacheOptions(publicDomainRocket, page, search, nil), func() (models.PublicRocketPage, error) {
		const pageSize = 20

		rocketList, err := m.GetRockets(RocketQuery{
			Limit:     pageSize,
			Offset:    page * pageSize,
			Search:    search,
			SortBy:    "total_launch_count",
			SortOrder: -1,
		})
		if err != nil {
			return models.PublicRocketPage{}, err
		}

		rockets := make([]models.PublicRocketListItem, 0, len(rocketList.Rockets))
		for _, rocket := range rocketList.Rockets {
			if rocket.Rocket.ID == 0 {
				continue
			}
			rockets = append(rockets, buildPublicRocketListItem(rocket.Rocket, rocket.Data))
		}

		return models.PublicRocketPage{
			Count:   rocketList.Count,
			Rockets: rockets,
		}, nil
	})
}

func (m *MainService) GetPublicRocket(id int64) (models.PublicRocketView, error) {
	return loadPublicCached(m.publicCache, m.publicDetailCacheOptions(publicDomainRocket, id), func() (models.PublicRocketView, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		rocket, err := m.findRocketByCoreID(ctx, id)
		if err != nil {
			return models.PublicRocketView{}, err
		}

		var ll2Data models.LL2LauncherConfigDetailed
		err = m.mc.Collection(COLLECTION_LL2_LAUNCHER).FindOne(ctx, bson.M{"id": rocket.ExternalID}).Decode(&ll2Data)
		if err != nil {
			return models.PublicRocketView{}, err
		}

		var company models.PublicCompanyRef
		if ll2Data.Manufacturer.ID != 0 {
			var agency models.Agency
			err := m.mc.Collection(COLLECTION_AGENCY).FindOne(ctx, bson.M{"external_id": ll2Data.Manufacturer.ID}).Decode(&agency)
			if err != nil && err != mongo.ErrNoDocuments {
				return models.PublicRocketView{}, err
			}

			var agencyDoc *models.LL2AgencyDetailed
			var ll2Agency models.LL2AgencyDetailed
			err = m.mc.Collection(COLLECTION_LL2_AGENCY).FindOne(ctx, bson.M{"id": ll2Data.Manufacturer.ID}).Decode(&ll2Agency)
			if err != nil && err != mongo.ErrNoDocuments {
				return models.PublicRocketView{}, err
			}
			if err == nil {
				agencyDoc = &ll2Agency
			}

			company = buildPublicCompanyRef(agency, ll2Data.Manufacturer.Name, agencyDoc)
		}

		launchLimit := 10
		launchFilter := bson.M{"rocket.configuration.id": rocket.ExternalID}
		launchOpts := options.Find().SetSort(bson.D{{Key: "net", Value: -1}}).SetLimit(int64(launchLimit))

		launchCursor, err := m.mc.Collection(COLLECTION_LL2_LAUNCH).Find(ctx, launchFilter, launchOpts)
		if err != nil {
			return models.PublicRocketView{}, err
		}
		defer launchCursor.Close(ctx)

		var ll2Launches []models.LL2LaunchDetailed
		if err = launchCursor.All(ctx, &ll2Launches); err != nil {
			return models.PublicRocketView{}, err
		}

		var publicLaunches []models.PublicLaunchSummary
		if len(ll2Launches) > 0 {
			externalIDs := make([]string, 0, len(ll2Launches))
			agencyExternalIDs := make([]int64, 0, len(ll2Launches))
			baseExternalIDs := make([]int64, 0, len(ll2Launches))
			for _, l := range ll2Launches {
				externalIDs = append(externalIDs, l.ID)
				if l.LaunchServiceProvider.ID != 0 {
					agencyExternalIDs = append(agencyExternalIDs, int64(l.LaunchServiceProvider.ID))
				}
				if l.Pad.Location.ID != 0 {
					baseExternalIDs = append(baseExternalIDs, int64(l.Pad.Location.ID))
				}
			}

			launchMap, err := m.loadLaunchMapByExternalIDs(ctx, externalIDs)
			if err != nil {
				return models.PublicRocketView{}, err
			}
			agencyMap, err := m.loadAgencyMapByExternalIDs(ctx, agencyExternalIDs)
			if err != nil {
				return models.PublicRocketView{}, err
			}
			agencyDocMap, err := m.loadLL2AgencyMapByExternalIDs(ctx, agencyExternalIDs)
			if err != nil {
				return models.PublicRocketView{}, err
			}
			baseMap, err := m.loadLaunchBaseMapByExternalIDs(ctx, baseExternalIDs)
			if err != nil {
				return models.PublicRocketView{}, err
			}
			rocketMap := map[int64]models.Rocket{rocket.ExternalID: rocket}
			launcherMap := map[int64]models.LL2LauncherConfigNormal{rocket.ExternalID: ll2Data.LL2LauncherConfigNormal}

			for _, ll2 := range ll2Launches {
				launchSummary, include := m.buildPublicLaunchSummary(launchMap[ll2.ID], ll2, rocketMap, launcherMap, agencyMap, agencyDocMap, baseMap)
				if !include {
					continue
				}
				publicLaunches = append(publicLaunches, launchSummary)
			}
		} else {
			publicLaunches = []models.PublicLaunchSummary{}
		}

		return models.PublicRocketView{
			ID:              publicID(rocket.ID),
			Name:            ll2Data.Name,
			Description:     ll2Data.Description,
			Active:          ll2Data.Active,
			Reusable:        ll2Data.Reusable,
			LaunchImage:     resolveRocketLaunch(rocket, ll2Data.Image),
			MainImage:       resolveRocketMain(rocket, ll2Data.Image),
			ImageList:       rocket.ImageList,
			Company:         company,
			Launches:        publicLaunches,
			LaunchCost:      float64(ll2Data.LaunchCost),
			Diameter:        ll2Data.Diameter,
			Length:          ll2Data.Length,
			LiftoffThrust:   ll2Data.ToThrust,
			LaunchMass:      ll2Data.LaunchMass,
			LeoCapacity:     ll2Data.LeoCapacity,
			GtoCapacity:     ll2Data.GtoCapacity,
			GeoCapacity:     ll2Data.GeoCapacity,
			SsoCapacity:     ll2Data.SsoCapacity,
			TotalLaunches:   ll2Data.TotalLaunchCount,
			SuccessLaunches: ll2Data.SuccessfulLaunches,
			FailureLaunches: ll2Data.FailedLaunches,
			TotalLandings:   ll2Data.AttemptedLandings,
			SuccessLandings: ll2Data.SuccessfulLandings,
			FailureLandings: ll2Data.FailedLandings,
		}, nil
	})
}
