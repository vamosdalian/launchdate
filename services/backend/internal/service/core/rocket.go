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
	FullName  string
	Name      string
	Variant   string
	SortBy    string
	SortOrder int
}

func (q RocketQuery) sortFieldAndOrder() (string, int) {
	switch strings.ToLower(strings.TrimSpace(q.SortBy)) {
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
	filter := bson.M{
		"id": bson.M{"$ne": nil},
	}
	if fullName := strings.TrimSpace(q.FullName); fullName != "" {
		filter["full_name"] = bson.M{"$regex": fullName, "$options": "i"}
	}
	if name := strings.TrimSpace(q.Name); name != "" {
		filter["name"] = bson.M{"$regex": name, "$options": "i"}
	}
	if variant := strings.TrimSpace(q.Variant); variant != "" {
		filter["variant"] = bson.M{"$regex": variant, "$options": "i"}
	}
	return filter
}

func rocketAggregationPipeline(filter bson.M) mongo.Pipeline {
	return mongo.Pipeline{
		bson.D{{Key: "$match", Value: filter}},
		bson.D{{Key: "$lookup", Value: bson.M{
			"from":         COLLECTION_AGENCY,
			"localField":   "id",
			"foreignField": "external_id",
			"as":           "rocket_doc",
		}}},
		bson.D{{Key: "$unwind", Value: bson.M{"path": "$rocket_doc", "preserveNullAndEmptyArrays": false}}},
	}
}

func (s *MainService) GenerateRocketsFromLL2(launchers []int64) error {
	colleciton := s.mc.Collection(COLLECTION_ROCKET)

	var docs []any
	for _, ll2id := range launchers {
		docs = append(docs, models.Rocket{
			ID:         s.sn.Generate().Int64(),
			ExternalID: ll2id,
		})
	}
	_, err := colleciton.InsertMany(context.Background(), docs)

	return err
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

	collection := m.mc.Collection(COLLECTION_AGENCY)

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

	collection := m.mc.Collection(COLLECTION_AGENCY)

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

	return nil
}

func (m *MainService) GetPublicRockets(page int) (models.PublicRocketList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	limit := 20
	offset := page * limit

	launcherCollection := m.mc.Collection(COLLECTION_LL2_LAUNCHER)

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$lookup", Value: bson.M{
			"from":         COLLECTION_AGENCY,
			"localField":   "id",
			"foreignField": "external_id",
			"as":           "rocket_doc",
		}}},
		bson.D{{Key: "$unwind", Value: bson.M{"path": "$rocket_doc", "preserveNullAndEmptyArrays": false}}},
		bson.D{{Key: "$sort", Value: bson.D{{Key: "total_launch_count", Value: -1}}}},
		bson.D{{Key: "$skip", Value: offset}},
		bson.D{{Key: "$limit", Value: limit}},
	}

	cursor, err := launcherCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return models.PublicRocketList{}, err
	}
	defer cursor.Close(ctx)

	var rockets []models.PublicCompactRocket
	for cursor.Next(ctx) {
		var doc rocketAggregateDoc
		if err := cursor.Decode(&doc); err != nil {
			return models.PublicRocketList{}, err
		}
		rockets = append(rockets, models.PublicCompactRocket{
			ID:         doc.RocketDoc.ID,
			Name:       doc.LL2LauncherConfigNormal.Name,
			ThumbImage: doc.RocketDoc.ThumbImage,
		})
	}

	if err := cursor.Err(); err != nil {
		return models.PublicRocketList{}, err
	}

	countPipeline := mongo.Pipeline{
		bson.D{{Key: "$lookup", Value: bson.M{
			"from":         COLLECTION_AGENCY,
			"localField":   "id",
			"foreignField": "external_id",
			"as":           "rocket_doc",
		}}},
		bson.D{{Key: "$unwind", Value: bson.M{"path": "$rocket_doc", "preserveNullAndEmptyArrays": false}}},
		bson.D{{Key: "$count", Value: "total"}},
	}
	countCursor, err := launcherCollection.Aggregate(ctx, countPipeline)
	if err != nil {
		return models.PublicRocketList{}, err
	}
	defer countCursor.Close(ctx)

	var total int64
	if countCursor.Next(ctx) {
		var countDoc struct {
			Total int64 `bson:"total"`
		}
		if err := countCursor.Decode(&countDoc); err != nil {
			return models.PublicRocketList{}, err
		}
		total = countDoc.Total
	}

	return models.PublicRocketList{
		Count:   int(total),
		Rockets: rockets,
	}, nil
}

func (m *MainService) GetPublicRocket(id int64) (models.PublicRocketDetail, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var rocket models.Rocket
	err := m.mc.Collection(COLLECTION_AGENCY).FindOne(ctx, bson.M{"id": id}).Decode(&rocket)
	if err != nil {
		return models.PublicRocketDetail{}, err
	}

	var ll2Data models.LL2LauncherConfigDetailed
	err = m.mc.Collection(COLLECTION_LL2_LAUNCHER).FindOne(ctx, bson.M{"id": rocket.ExternalID}).Decode(&ll2Data)
	if err != nil {
		return models.PublicRocketDetail{}, err
	}

	var agencyInfo models.PublicCompactAgency
	if ll2Data.Manufacturer.ID != 0 {
		var agency models.Agency
		err := m.mc.Collection(COLLECTION_AGENCY).FindOne(ctx, bson.M{"external_id": ll2Data.Manufacturer.ID}).Decode(&agency)
		if err == nil {
			agencyInfo = models.PublicCompactAgency{
				ID:         agency.ID,
				Name:       ll2Data.Manufacturer.Name,
				ThumbImage: agency.ThumbImage,
			}
		} else {
			// If internal agency not found, at least provide the name from LL2
			agencyInfo = models.PublicCompactAgency{
				Name: ll2Data.Manufacturer.Name,
			}
		}
	}

	launchLimit := 10
	launchFilter := bson.M{"rocket.configuration.id": rocket.ExternalID}
	launchOpts := options.Find().SetSort(bson.D{{Key: "net", Value: -1}}).SetLimit(int64(launchLimit))

	launchCursor, err := m.mc.Collection(COLLECTION_LL2_LAUNCH).Find(ctx, launchFilter, launchOpts)
	if err != nil {
		return models.PublicRocketDetail{}, err
	}
	defer launchCursor.Close(ctx)

	var ll2Launches []models.LL2LaunchDetailed
	if err = launchCursor.All(ctx, &ll2Launches); err != nil {
		return models.PublicRocketDetail{}, err
	}

	var publicLaunches []models.PublicCompactLaunch
	if len(ll2Launches) > 0 {
		externalIDs := make([]string, 0, len(ll2Launches))
		for _, l := range ll2Launches {
			externalIDs = append(externalIDs, l.ID)
		}

		launchMap := make(map[string]models.Launch)
		lCursor, err := m.mc.Collection(COLLECTION_LAUNCH).Find(ctx, bson.M{"external_id": bson.M{"$in": externalIDs}})
		if err == nil {
			defer lCursor.Close(ctx)
			for lCursor.Next(ctx) {
				var l models.Launch
				if err := lCursor.Decode(&l); err == nil {
					launchMap[l.ExternalID] = l
				}
			}
		}

		for _, ll2 := range ll2Launches {
			internalLaunch := launchMap[ll2.ID]
			publicLaunches = append(publicLaunches, models.PublicCompactLaunch{
				ID:         internalLaunch.ID,
				Name:       ll2.Name,
				LaunchTime: ll2.Net,
				Status:     ll2.Status.ID,
				ThumbImage: internalLaunch.ThumbImage,
				RocketName: ll2.Rocket.Configuration.Name,
				AgencyName: ll2.LaunchServiceProvider.Name,
				Location:   ll2.Pad.Location.Name,
			})
		}
	} else {
		publicLaunches = []models.PublicCompactLaunch{}
	}

	return models.PublicRocketDetail{
		ID:              rocket.ID,
		Name:            ll2Data.Name,
		Description:     ll2Data.Description,
		Active:          ll2Data.Active,
		Reusable:        ll2Data.Reusable,
		LaunchImage:     rocket.LaunchImage,
		MainImage:       rocket.MainImage,
		ImageList:       rocket.ImageList,
		AgencyInfo:      agencyInfo,
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
}
