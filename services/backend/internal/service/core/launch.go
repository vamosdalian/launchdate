package core

import (
	"context"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/vamosdalian/launchdate-backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const COLLECTION_LAUNCH = "launch"
const COLLECTION_LL2_LAUNCH = "ll2_launch"

func (s *MainService) GenerateLaunchFromLL2(ll2ids []string) error {
	colleciton := s.mc.Collection(COLLECTION_LAUNCH)

	var docs []any
	for _, ll2id := range ll2ids {
		docs = append(docs, models.Launch{
			ID:         s.sn.Generate().Int64(),
			ExternalID: ll2id,
		})
	}
	_, err := colleciton.InsertMany(context.Background(), docs)

	return err
}

type LaunchQuery struct {
	Limit         int
	Offset        int
	Name          string
	Status        string
	LaunchService string
	Rocket        string
	Mission       string
	SortBy        string
	SortOrder     int
}

func (q LaunchQuery) sortFieldAndOrder() (string, int) {
	switch strings.ToLower(strings.TrimSpace(q.SortBy)) {
	case "name":
		if q.SortOrder < 0 {
			return "name", -1
		}
		return "name", 1
	default:
		if q.SortOrder < 0 {
			return "net", -1
		}
		return "net", 1
	}
}

func buildLaunchFilter(q LaunchQuery) bson.M {
	filter := bson.M{}
	if name := strings.TrimSpace(q.Name); name != "" {
		filter["name"] = bson.M{"$regex": name, "$options": "i"}
	}
	if status := strings.TrimSpace(q.Status); status != "" {
		filter["status.name"] = bson.M{"$regex": status, "$options": "i"}
	}
	if provider := strings.TrimSpace(q.LaunchService); provider != "" {
		filter["launch_service_provider.name"] = bson.M{"$regex": provider, "$options": "i"}
	}
	if rocket := strings.TrimSpace(q.Rocket); rocket != "" {
		filter["rocket.configuration.full_name"] = bson.M{"$regex": rocket, "$options": "i"}
	}
	if mission := strings.TrimSpace(q.Mission); mission != "" {
		filter["mission.name"] = bson.M{"$regex": mission, "$options": "i"}
	}
	return filter
}

func (m *MainService) GetLaunches(q LaunchQuery) (models.LaunchList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ll2Collection := m.mc.Collection(COLLECTION_LL2_LAUNCH)
	filter := buildLaunchFilter(q)

	total, err := ll2Collection.CountDocuments(ctx, filter)
	if err != nil {
		return models.LaunchList{}, err
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
		return models.LaunchList{}, err
	}
	defer cursor.Close(ctx)

	ll2Docs := make([]models.LL2LaunchDetailed, 0)
	externalIDs := make([]string, 0)

	for cursor.Next(ctx) {
		var doc models.LL2LaunchDetailed
		if err := cursor.Decode(&doc); err != nil {
			return models.LaunchList{}, err
		}
		ll2Docs = append(ll2Docs, doc)
		if doc.ID != "" {
			externalIDs = append(externalIDs, doc.ID)
		}
	}

	if err := cursor.Err(); err != nil {
		return models.LaunchList{}, err
	}

	launches := make(map[string]models.Launch, len(externalIDs))
	if len(externalIDs) > 0 {
		launchCursor, err := m.mc.Collection(COLLECTION_LAUNCH).
			Find(ctx, bson.M{"external_id": bson.M{"$in": externalIDs}})
		if err != nil {
			return models.LaunchList{}, err
		}
		defer launchCursor.Close(ctx)

		for launchCursor.Next(ctx) {
			var launch models.Launch
			if err := launchCursor.Decode(&launch); err != nil {
				return models.LaunchList{}, err
			}
			launches[launch.ExternalID] = launch
		}

		if err := launchCursor.Err(); err != nil {
			return models.LaunchList{}, err
		}
	}

	serializers := make([]models.LaunchSerializer, 0, len(ll2Docs))
	for _, doc := range ll2Docs {
		serializers = append(serializers, models.LaunchSerializer{
			Launch: launches[doc.ID],
			Data:   doc.LL2LaunchNormal,
		})
	}

	return models.LaunchList{
		Count:    int(total),
		Launches: serializers,
	}, nil
}

func (m *MainService) GetLaunch(id int64) (models.LaunchSerializer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := m.mc.Collection(COLLECTION_LAUNCH)

	var launch models.Launch
	if err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&launch); err != nil {
		return models.LaunchSerializer{}, err
	}

	var data models.LL2LaunchNormal
	if launch.ExternalID != "" {
		var doc models.LL2LaunchDetailed
		err := m.mc.Collection(COLLECTION_LL2_LAUNCH).FindOne(ctx, bson.M{"id": launch.ExternalID}).Decode(&doc)
		if err != nil && err != mongo.ErrNoDocuments {
			return models.LaunchSerializer{}, err
		}
		if err == nil {
			data = doc.LL2LaunchNormal
		}
	}

	return models.LaunchSerializer{
		Launch: launch,
		Data:   data,
	}, nil
}

func (m *MainService) UpdateLaunch(l *models.Launch) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := m.mc.Collection(COLLECTION_LAUNCH)

	update := bson.M{
		"$set": bson.M{
			"background_image": l.BackgroundImage,
			"image_list":       l.ImageList,
			"thumb_image":      l.ThumbImage,
		},
	}

	result, err := collection.UpdateOne(ctx, bson.M{"id": l.ID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (m *MainService) GetPublicLaunch(id int64) (models.PublicLaunchDetail, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Get Launch
	var launch models.Launch
	err := m.mc.Collection(COLLECTION_LAUNCH).FindOne(ctx, bson.M{"id": id}).Decode(&launch)
	if err != nil {
		return models.PublicLaunchDetail{}, err
	}

	// 2. Get LL2Launch
	var ll2Launch models.LL2LaunchDetailed
	if launch.ExternalID != "" {
		err = m.mc.Collection(COLLECTION_LL2_LAUNCH).FindOne(ctx, bson.M{"id": launch.ExternalID}).Decode(&ll2Launch)
		if err != nil && err != mongo.ErrNoDocuments {
			return models.PublicLaunchDetail{}, err
		}
	}

	// 3. Get Rocket
	var rocket models.Rocket
	if ll2Launch.Rocket.Configuration.ID != 0 {
		// Assuming external_id in rocket collection matches LL2 ID
		err = m.mc.Collection(COLLECTION_AGENCY).FindOne(ctx, bson.M{"external_id": ll2Launch.Rocket.Configuration.ID}).Decode(&rocket)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				logrus.Errorf("couldn't find rocket for launch %d with ll2 rocket id %d", id, ll2Launch.Rocket.Configuration.ID)
			} else {
				logrus.Errorf("failed to get rocket for launch %d: %v", id, err)
			}
		}
	}

	// 4. Get Agency
	var agency models.Agency
	if ll2Launch.LaunchServiceProvider.ID != 0 {
		err = m.mc.Collection(COLLECTION_AGENCY).FindOne(ctx, bson.M{"external_id": ll2Launch.LaunchServiceProvider.ID}).Decode(&agency)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				logrus.Errorf("couldn't find agency for launch %d with ll2 agency id %d", id, ll2Launch.LaunchServiceProvider.ID)
			} else {
				logrus.Errorf("failed to get agency for launch %d: %v", id, err)
			}
		}
	}

	// 5. Construct Response
	detail := models.PublicLaunchDetail{
		ID:              launch.ID,
		Name:            ll2Launch.Name,
		LaunchTime:      ll2Launch.Net,
		Status:          ll2Launch.Status.ID,
		BackgroundImage: launch.BackgroundImage,
		ImageList:       launch.ImageList,
		RocketInfo: models.PublicCompactRocket{
			ID:         rocket.ID,
			Name:       ll2Launch.Rocket.Configuration.Name,
			ThumbImage: rocket.ThumbImage,
		},
		AgencyInfo: models.PublicCompactAgency{
			ID:         agency.ID,
			Name:       ll2Launch.LaunchServiceProvider.Name,
			ThumbImage: agency.ThumbImage,
		},
		LocationInfo: models.PublicCompactLocation{
			ID:   int64(ll2Launch.Pad.Location.ID),
			Name: ll2Launch.Pad.Location.Name,
			Lat:  ll2Launch.Pad.Location.Latitude,
			Lon:  ll2Launch.Pad.Location.Longitude,
		},
		MissionInfo:   []models.Mission{},
		TimelineEvent: []models.TimelineEvent{},
	}

	for _, event := range ll2Launch.Timeline {
		detail.TimelineEvent = append(detail.TimelineEvent, models.TimelineEvent{
			RelativeTime: event.RelativeTime,
			Abbrev:       event.Type.Abbrev,
			Description:  event.Type.Description,
		})
	}

	if ll2Launch.Mission.ID != 0 {
		detail.MissionInfo = append(detail.MissionInfo, models.Mission{
			ID:          int64(ll2Launch.Mission.ID),
			Name:        ll2Launch.Mission.Name,
			Description: ll2Launch.Mission.Description,
		})
	}

	return detail, nil
}

func (m *MainService) GetPublicLaunches(page int) (models.PublicLaunchList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	now := time.Now().UTC().Format(time.RFC3339)
	ll2Collection := m.mc.Collection(COLLECTION_LL2_LAUNCH)

	var futureLaunches []models.LL2LaunchDetailed

	// 1. Future launches (next 14 days) - Only for page 0
	if page == 0 {
		futureLimit := time.Now().UTC().AddDate(0, 0, 14).Format(time.RFC3339)
		futureFilter := bson.M{
			"net": bson.M{
				"$gte": now,
				"$lte": futureLimit,
			},
		}
		futureCursor, err := ll2Collection.Find(ctx, futureFilter, options.Find().SetSort(bson.D{{Key: "net", Value: 1}}))
		if err != nil {
			return models.PublicLaunchList{}, err
		}
		defer futureCursor.Close(ctx)

		if err = futureCursor.All(ctx, &futureLaunches); err != nil {
			return models.PublicLaunchList{}, err
		}
	}

	// 2. Past launches (latest 20)
	pastFilter := bson.M{
		"net": bson.M{
			"$lt": now,
		},
	}
	limit := int64(20)
	skip := int64(page) * limit

	pastCursor, err := ll2Collection.Find(ctx, pastFilter, options.Find().SetSort(bson.D{{Key: "net", Value: -1}}).SetLimit(limit).SetSkip(skip))
	if err != nil {
		return models.PublicLaunchList{}, err
	}
	defer pastCursor.Close(ctx)

	var pastLaunches []models.LL2LaunchDetailed
	if err = pastCursor.All(ctx, &pastLaunches); err != nil {
		return models.PublicLaunchList{}, err
	}

	// Combine results
	allLL2Launches := append(futureLaunches, pastLaunches...)
	if len(allLL2Launches) == 0 {
		return models.PublicLaunchList{Count: 0, Launches: []models.PublicCompactLaunch{}}, nil
	}

	// 3. Get internal launch data
	externalIDs := make([]string, 0, len(allLL2Launches))
	for _, l := range allLL2Launches {
		if l.ID != "" {
			externalIDs = append(externalIDs, l.ID)
		}
	}

	launchMap := make(map[string]models.Launch)
	if len(externalIDs) > 0 {
		launchCursor, err := m.mc.Collection(COLLECTION_LAUNCH).Find(ctx, bson.M{"external_id": bson.M{"$in": externalIDs}})
		if err != nil {
			return models.PublicLaunchList{}, err
		}
		defer launchCursor.Close(ctx)

		for launchCursor.Next(ctx) {
			var l models.Launch
			if err := launchCursor.Decode(&l); err != nil {
				return models.PublicLaunchList{}, err
			}
			launchMap[l.ExternalID] = l
		}
	}

	// 4. Construct response
	publicLaunches := make([]models.PublicCompactLaunch, 0, len(allLL2Launches))
	for _, ll2 := range allLL2Launches {
		internalLaunch, ok := launchMap[ll2.ID]
		if !ok {
			logrus.Errorf("couldn't find internal launch for ll2 launch id %s", ll2.ID)
		}

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

	return models.PublicLaunchList{
		Count:    len(publicLaunches),
		Launches: publicLaunches,
	}, nil
}
