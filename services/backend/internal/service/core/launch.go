package core

import (
	"context"
	"sort"
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
	return s.ensureStringExternalIDs(COLLECTION_LAUNCH, ll2ids)
}

type LaunchQuery struct {
	Limit         int
	Offset        int
	Search        string
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
	filters := make([]bson.M, 0, 6)
	if searchClause := buildTextSearchClause(q.Search, "name", "launch_service_provider.name", "rocket.configuration.full_name", "mission.name", "pad.location.name", "pad.name"); len(searchClause) > 0 {
		filters = append(filters, searchClause)
	}
	if name := strings.TrimSpace(q.Name); name != "" {
		filters = append(filters, bson.M{"name": bson.M{"$regex": name, "$options": "i"}})
	}
	if status := strings.TrimSpace(q.Status); status != "" {
		filters = append(filters, bson.M{"status.name": bson.M{"$regex": status, "$options": "i"}})
	}
	if provider := strings.TrimSpace(q.LaunchService); provider != "" {
		filters = append(filters, bson.M{"launch_service_provider.name": bson.M{"$regex": provider, "$options": "i"}})
	}
	if rocket := strings.TrimSpace(q.Rocket); rocket != "" {
		filters = append(filters, bson.M{"rocket.configuration.full_name": bson.M{"$regex": rocket, "$options": "i"}})
	}
	if mission := strings.TrimSpace(q.Mission); mission != "" {
		filters = append(filters, bson.M{"mission.name": bson.M{"$regex": mission, "$options": "i"}})
	}
	return combineFilters(filters...)
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

	m.bumpPublicCacheDomains(publicDomainLaunch, publicDomainRocket, publicDomainCompany, publicDomainLaunchBase)

	return nil
}

func (m *MainService) GetPublicLaunch(id int64) (models.PublicLaunchView, error) {
	return loadPublicCached(m.publicCache, m.publicDetailCacheOptions(publicDomainLaunch, id), func() (models.PublicLaunchView, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var launch models.Launch
		err := m.mc.Collection(COLLECTION_LAUNCH).FindOne(ctx, bson.M{"id": id}).Decode(&launch)
		if err != nil {
			return models.PublicLaunchView{}, err
		}

		var ll2Launch models.LL2LaunchDetailed
		if launch.ExternalID != "" {
			err = m.mc.Collection(COLLECTION_LL2_LAUNCH).FindOne(ctx, bson.M{"id": launch.ExternalID}).Decode(&ll2Launch)
			if err != nil && err != mongo.ErrNoDocuments {
				return models.PublicLaunchView{}, err
			}
		}

		rocketMap, err := m.loadRocketMapByExternalIDs(ctx, []int64{int64(ll2Launch.Rocket.Configuration.ID)})
		if err != nil {
			return models.PublicLaunchView{}, err
		}
		launcherMap, err := m.loadLauncherConfigMapByExternalIDs(ctx, []int64{int64(ll2Launch.Rocket.Configuration.ID)})
		if err != nil {
			return models.PublicLaunchView{}, err
		}
		agencyMap, err := m.loadAgencyMapByExternalIDs(ctx, []int64{int64(ll2Launch.LaunchServiceProvider.ID)})
		if err != nil {
			return models.PublicLaunchView{}, err
		}
		agencyDocMap, err := m.loadLL2AgencyMapByExternalIDs(ctx, []int64{int64(ll2Launch.LaunchServiceProvider.ID)})
		if err != nil {
			return models.PublicLaunchView{}, err
		}
		baseMap, err := m.loadLaunchBaseMapByExternalIDs(ctx, []int64{int64(ll2Launch.Pad.Location.ID)})
		if err != nil {
			return models.PublicLaunchView{}, err
		}

		launchSummary, ok := m.buildPublicLaunchSummary(
			launch,
			ll2Launch,
			rocketMap,
			launcherMap,
			agencyMap,
			agencyDocMap,
			baseMap,
		)
		if !ok {
			return models.PublicLaunchView{}, mongo.ErrNoDocuments
		}

		detail := models.PublicLaunchView{
			ID:              launchSummary.ID,
			Name:            launchSummary.Name,
			LaunchTime:      launchSummary.LaunchTime,
			Status:          launchSummary.Status,
			StatusLabel:     launchSummary.StatusLabel,
			BackgroundImage: launchSummary.BackgroundImage,
			ImageList:       launch.ImageList,
			Rocket:          launchSummary.Rocket,
			Company:         launchSummary.Company,
			LaunchBase:      launchSummary.LaunchBase,
			Missions:        []models.PublicMissionSummary{},
			Timeline:        []models.PublicTimelineEntry{},
		}

		for _, event := range ll2Launch.Timeline {
			detail.Timeline = append(detail.Timeline, models.PublicTimelineEntry{
				RelativeTime: event.RelativeTime,
				Abbrev:       event.Type.Abbrev,
				Description:  event.Type.Description,
			})
		}

		if ll2Launch.Mission.ID != 0 {
			detail.Missions = append(detail.Missions, models.PublicMissionSummary{
				Name:        ll2Launch.Mission.Name,
				Description: ll2Launch.Mission.Description,
			})
		}

		return detail, nil
	})
}

func (m *MainService) GetPublicLaunches(page int, search string) (models.PublicLaunchPage, error) {
	return loadPublicCached(m.publicCache, m.publicListCacheOptions(publicDomainLaunch, page, search, nil), func() (models.PublicLaunchPage, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		now := time.Now().UTC().Format(time.RFC3339)
		ll2Collection := m.mc.Collection(COLLECTION_LL2_LAUNCH)
		baseFilter := buildLaunchFilter(LaunchQuery{Search: search})
		const pageSize = int64(20)
		var err error

		var futureLaunches []models.LL2LaunchDetailed
		var futureCount int64

		if page == 0 {
			futureLimit := time.Now().UTC().AddDate(0, 0, 14).Format(time.RFC3339)
			futureFilter := combineFilters(baseFilter, bson.M{
				"net": bson.M{
					"$gte": now,
					"$lte": futureLimit,
				},
			})
			futureCount, err = ll2Collection.CountDocuments(ctx, futureFilter)
			if err != nil {
				return models.PublicLaunchPage{}, err
			}
			futureCursor, err := ll2Collection.Find(ctx, futureFilter, options.Find().SetSort(bson.D{{Key: "net", Value: -1}}))
			if err != nil {
				return models.PublicLaunchPage{}, err
			}
			defer futureCursor.Close(ctx)

			if err = futureCursor.All(ctx, &futureLaunches); err != nil {
				return models.PublicLaunchPage{}, err
			}
		}

		pastFilter := combineFilters(baseFilter, bson.M{
			"net": bson.M{
				"$lt": now,
			},
		})
		pastCount, err := ll2Collection.CountDocuments(ctx, pastFilter)
		if err != nil {
			return models.PublicLaunchPage{}, err
		}
		skip := int64(page) * pageSize

		pastCursor, err := ll2Collection.Find(ctx, pastFilter, options.Find().SetSort(bson.D{{Key: "net", Value: -1}}).SetLimit(pageSize).SetSkip(skip))
		if err != nil {
			return models.PublicLaunchPage{}, err
		}
		defer pastCursor.Close(ctx)

		var pastLaunches []models.LL2LaunchDetailed
		if err = pastCursor.All(ctx, &pastLaunches); err != nil {
			return models.PublicLaunchPage{}, err
		}

		allLL2Launches := append(futureLaunches, pastLaunches...)
		if len(allLL2Launches) == 0 {
			return models.PublicLaunchPage{Count: 0, Launches: []models.PublicLaunchSummary{}}, nil
		}

		externalIDs := make([]string, 0, len(allLL2Launches))
		for _, l := range allLL2Launches {
			if l.ID != "" {
				externalIDs = append(externalIDs, l.ID)
			}
		}

		launchMap, err := m.loadLaunchMapByExternalIDs(ctx, externalIDs)
		if err != nil {
			return models.PublicLaunchPage{}, err
		}

		rocketExternalIDs := make([]int64, 0, len(allLL2Launches))
		agencyExternalIDs := make([]int64, 0, len(allLL2Launches))
		baseExternalIDs := make([]int64, 0, len(allLL2Launches))
		for _, ll2 := range allLL2Launches {
			if ll2.Rocket.Configuration.ID != 0 {
				rocketExternalIDs = append(rocketExternalIDs, int64(ll2.Rocket.Configuration.ID))
			}
			if ll2.LaunchServiceProvider.ID != 0 {
				agencyExternalIDs = append(agencyExternalIDs, int64(ll2.LaunchServiceProvider.ID))
			}
			if ll2.Pad.Location.ID != 0 {
				baseExternalIDs = append(baseExternalIDs, int64(ll2.Pad.Location.ID))
			}
		}

		rocketMap, err := m.loadRocketMapByExternalIDs(ctx, rocketExternalIDs)
		if err != nil {
			return models.PublicLaunchPage{}, err
		}
		launcherMap, err := m.loadLauncherConfigMapByExternalIDs(ctx, rocketExternalIDs)
		if err != nil {
			return models.PublicLaunchPage{}, err
		}
		agencyMap, err := m.loadAgencyMapByExternalIDs(ctx, agencyExternalIDs)
		if err != nil {
			return models.PublicLaunchPage{}, err
		}
		agencyDocMap, err := m.loadLL2AgencyMapByExternalIDs(ctx, agencyExternalIDs)
		if err != nil {
			return models.PublicLaunchPage{}, err
		}
		baseMap, err := m.loadLaunchBaseMapByExternalIDs(ctx, baseExternalIDs)
		if err != nil {
			return models.PublicLaunchPage{}, err
		}

		publicLaunches := make([]models.PublicLaunchSummary, 0, len(allLL2Launches))
		for _, ll2 := range allLL2Launches {
			internalLaunch, ok := launchMap[ll2.ID]
			if !ok {
				logrus.Errorf("couldn't find internal launch for ll2 launch id %s", ll2.ID)
				continue
			}

			launchSummary, include := m.buildPublicLaunchSummary(internalLaunch, ll2, rocketMap, launcherMap, agencyMap, agencyDocMap, baseMap)
			if !include {
				continue
			}
			publicLaunches = append(publicLaunches, launchSummary)
		}

		nowTime := time.Now().UTC()
		sort.SliceStable(publicLaunches, func(leftIndex, rightIndex int) bool {
			leftLaunchTime, leftErr := time.Parse(time.RFC3339, publicLaunches[leftIndex].LaunchTime)
			rightLaunchTime, rightErr := time.Parse(time.RFC3339, publicLaunches[rightIndex].LaunchTime)
			if leftErr != nil || rightErr != nil {
				return publicLaunches[leftIndex].LaunchTime < publicLaunches[rightIndex].LaunchTime
			}

			leftIsUpcoming := !leftLaunchTime.Before(nowTime)
			rightIsUpcoming := !rightLaunchTime.Before(nowTime)

			if leftIsUpcoming != rightIsUpcoming {
				return leftIsUpcoming
			}

			if leftIsUpcoming {
				return leftLaunchTime.After(rightLaunchTime)
			}

			return leftLaunchTime.After(rightLaunchTime)
		})

		return models.PublicLaunchPage{
			Count:    int(futureCount + pastCount),
			Launches: publicLaunches,
		}, nil
	})
}
