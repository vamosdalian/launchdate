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
	return s.ensureInt64ExternalIDs(COLLECTION_AGENCY, ll2ids)
}

type AgencyQuery struct {
	Limit      int
	Offset     int
	Search     string
	Name       string
	Type       string
	Country    string
	SortBy     string
	SortOrder  int
	ShowOnHome *bool
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

func buildAgencyFilter(q AgencyQuery, externalIDs []int) bson.M {
	filters := make([]bson.M, 0, 4)
	if searchClause := buildTextSearchClause(q.Search, "name", "type.name", "country.name"); len(searchClause) > 0 {
		filters = append(filters, searchClause)
	}
	if len(externalIDs) > 0 {
		filters = append(filters, bson.M{"id": bson.M{"$in": externalIDs}})
	}
	if name := strings.TrimSpace(q.Name); name != "" {
		filters = append(filters, bson.M{"name": bson.M{"$regex": name, "$options": "i"}})
	}
	if agencyType := strings.TrimSpace(q.Type); agencyType != "" {
		filters = append(filters, bson.M{"type.name": bson.M{"$regex": agencyType, "$options": "i"}})
	}
	if country := strings.TrimSpace(q.Country); country != "" {
		filters = append(filters, bson.M{"country.name": bson.M{"$regex": country, "$options": "i"}})
	}
	return combineFilters(filters...)
}

func (s *MainService) loadAgencyExternalIDsByHomeFlag(ctx context.Context, showOnHome bool) ([]int, error) {
	cursor, err := s.mc.Collection(COLLECTION_AGENCY).Find(
		ctx,
		bson.M{"show_on_home": showOnHome},
		options.Find().SetProjection(bson.M{"external_id": 1}),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	externalIDs := make([]int, 0)
	for cursor.Next(ctx) {
		var agency models.Agency
		if err := cursor.Decode(&agency); err != nil {
			return nil, err
		}
		if agency.ExternalID == 0 {
			continue
		}
		externalIDs = append(externalIDs, int(agency.ExternalID))
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return externalIDs, nil
}

func (s *MainService) GetAgencies(q AgencyQuery) (models.AgencyList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ll2Collection := s.mc.Collection(COLLECTION_LL2_AGENCY)
	allowedExternalIDs := []int(nil)
	if q.ShowOnHome != nil {
		var err error
		allowedExternalIDs, err = s.loadAgencyExternalIDsByHomeFlag(ctx, *q.ShowOnHome)
		if err != nil {
			return models.AgencyList{}, err
		}
		if len(allowedExternalIDs) == 0 {
			return models.AgencyList{Count: 0, Agencies: []models.AgencySerializer{}}, nil
		}
	}

	filter := buildAgencyFilter(q, allowedExternalIDs)

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
			"thumb_image":  a.ThumbImage,
			"images":       a.Images,
			"social_url":   a.SocialUrl,
			"show_on_home": a.ShowOnHome,
		},
	}

	result, err := collection.UpdateOne(ctx, bson.M{"id": a.ID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	m.bumpPublicCacheDomains(publicDomainCompany, publicDomainRocket, publicDomainLaunch)

	return nil
}

func firstAgencyCountryName(countries []models.LL2Country) string {
	if len(countries) == 0 {
		return ""
	}
	return countries[0].Name
}

func isPreferredWebsiteLabel(value string) bool {
	label := strings.ToLower(strings.TrimSpace(value))
	switch label {
	case "website", "homepage", "official", "official website", "site", "official site":
		return true
	default:
		return false
	}
}

func firstSocialURL(urls []models.SocialUrl) string {
	for _, social := range urls {
		if isPreferredWebsiteLabel(social.Name) && strings.TrimSpace(social.URL) != "" {
			return social.URL
		}
	}
	return ""
}

func firstDetailedSocialURL(urls []models.LL2SocialMediaLink) string {
	for _, social := range urls {
		if isPreferredWebsiteLabel(social.SocialMedia.Name) && strings.TrimSpace(social.URL) != "" {
			return social.URL
		}
	}
	return ""
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func (s *MainService) GetPublicCompanies(page int, search string, showOnHomeOnly bool) (models.PublicCompanyPage, error) {
	extra := map[string]string{}
	if showOnHomeOnly {
		extra["homepage_only"] = "true"
	}

	return loadPublicCached(s.publicCache, s.publicListCacheOptions(publicDomainCompany, page, search, extra), func() (models.PublicCompanyPage, error) {
		const pageSize = 20
		var showOnHome *bool
		if showOnHomeOnly {
			showOnHome = &showOnHomeOnly
		}

		agencyList, err := s.GetAgencies(AgencyQuery{
			Limit:      pageSize,
			Offset:     page * pageSize,
			Search:     search,
			SortBy:     "name",
			SortOrder:  1,
			ShowOnHome: showOnHome,
		})
		if err != nil {
			return models.PublicCompanyPage{}, err
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		externalIDs := make([]int64, 0, len(agencyList.Agencies))
		for _, agency := range agencyList.Agencies {
			if agency.ExternalID == 0 {
				continue
			}
			externalIDs = append(externalIDs, agency.ExternalID)
		}

		agencyDocMap, err := s.loadLL2AgencyMapByExternalIDs(ctx, externalIDs)
		if err != nil {
			return models.PublicCompanyPage{}, err
		}

		companies := make([]models.PublicCompanyListItem, 0, len(agencyList.Agencies))
		for _, agency := range agencyList.Agencies {
			if agency.Agency.ID == 0 {
				continue
			}
			doc, ok := agencyDocMap[agency.ExternalID]
			if !ok {
				doc = models.LL2AgencyDetailed{LL2AgencyNormal: agency.Data}
			}
			companies = append(companies, buildPublicCompanyListItem(agency.Agency, doc))
		}

		return models.PublicCompanyPage{
			Count:     agencyList.Count,
			Companies: companies,
		}, nil
	})
}

func (s *MainService) GetPublicCompany(id int64) (models.PublicCompanyView, error) {
	return loadPublicCached(s.publicCache, s.publicDetailCacheOptions(publicDomainCompany, id), func() (models.PublicCompanyView, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var agency models.Agency
		err := s.mc.Collection(COLLECTION_AGENCY).FindOne(ctx, bson.M{"id": id}).Decode(&agency)
		if err != nil {
			return models.PublicCompanyView{}, err
		}

		var doc models.LL2AgencyDetailed
		if err := s.mc.Collection(COLLECTION_LL2_AGENCY).FindOne(ctx, bson.M{"id": int(agency.ExternalID)}).Decode(&doc); err != nil {
			return models.PublicCompanyView{}, err
		}

		launcherCursor, err := s.mc.Collection(COLLECTION_LL2_LAUNCHER).Find(
			ctx,
			bson.M{"manufacturer.id": int(agency.ExternalID)},
			options.Find().SetSort(bson.D{{Key: "total_launch_count", Value: -1}}).SetLimit(12),
		)
		if err != nil {
			return models.PublicCompanyView{}, err
		}
		defer launcherCursor.Close(ctx)

		ll2Launchers := make([]models.LL2LauncherConfigNormal, 0)
		rocketExternalIDs := make([]int64, 0)
		for launcherCursor.Next(ctx) {
			var launcher models.LL2LauncherConfigNormal
			if err := launcherCursor.Decode(&launcher); err != nil {
				return models.PublicCompanyView{}, err
			}
			ll2Launchers = append(ll2Launchers, launcher)
			if launcher.ID != 0 {
				rocketExternalIDs = append(rocketExternalIDs, int64(launcher.ID))
			}
		}
		if err := launcherCursor.Err(); err != nil {
			return models.PublicCompanyView{}, err
		}

		rocketMap, err := s.loadRocketMapByExternalIDs(ctx, rocketExternalIDs)
		if err != nil {
			return models.PublicCompanyView{}, err
		}

		rockets := make([]models.PublicRocketListItem, 0, len(ll2Launchers))
		for _, launcher := range ll2Launchers {
			rocket := findExistingRocket(rocketMap, int64(launcher.ID))
			if rocket.ID == 0 {
				continue
			}
			rockets = append(rockets, buildPublicRocketListItem(rocket, launcher))
		}

		launchCursor, err := s.mc.Collection(COLLECTION_LL2_LAUNCH).Find(
			ctx,
			bson.M{"launch_service_provider.id": int(agency.ExternalID)},
			options.Find().SetSort(bson.D{{Key: "net", Value: -1}}).SetLimit(6),
		)
		if err != nil {
			return models.PublicCompanyView{}, err
		}
		defer launchCursor.Close(ctx)

		ll2Launches := make([]models.LL2LaunchDetailed, 0)
		launchExternalIDs := make([]string, 0)
		launchRocketExternalIDs := make([]int64, 0)
		launchBaseExternalIDs := make([]int64, 0)
		for launchCursor.Next(ctx) {
			var ll2Launch models.LL2LaunchDetailed
			if err := launchCursor.Decode(&ll2Launch); err != nil {
				return models.PublicCompanyView{}, err
			}
			ll2Launches = append(ll2Launches, ll2Launch)
			launchExternalIDs = append(launchExternalIDs, ll2Launch.ID)
			if ll2Launch.Rocket.Configuration.ID != 0 {
				launchRocketExternalIDs = append(launchRocketExternalIDs, int64(ll2Launch.Rocket.Configuration.ID))
			}
			if ll2Launch.Pad.Location.ID != 0 {
				launchBaseExternalIDs = append(launchBaseExternalIDs, int64(ll2Launch.Pad.Location.ID))
			}
		}
		if err := launchCursor.Err(); err != nil {
			return models.PublicCompanyView{}, err
		}

		launchMap, err := s.loadLaunchMapByExternalIDs(ctx, launchExternalIDs)
		if err != nil {
			return models.PublicCompanyView{}, err
		}
		launchRocketMap, err := s.loadRocketMapByExternalIDs(ctx, launchRocketExternalIDs)
		if err != nil {
			return models.PublicCompanyView{}, err
		}
		launcherMap, err := s.loadLauncherConfigMapByExternalIDs(ctx, launchRocketExternalIDs)
		if err != nil {
			return models.PublicCompanyView{}, err
		}
		baseMap, err := s.loadLaunchBaseMapByExternalIDs(ctx, launchBaseExternalIDs)
		if err != nil {
			return models.PublicCompanyView{}, err
		}
		agencyMap := map[int64]models.Agency{agency.ExternalID: agency}
		agencyDocMap := map[int64]models.LL2AgencyDetailed{agency.ExternalID: doc}

		launches := make([]models.PublicLaunchSummary, 0, len(ll2Launches))
		for _, ll2Launch := range ll2Launches {
			launchSummary, include := s.buildPublicLaunchSummary(launchMap[ll2Launch.ID], ll2Launch, launchRocketMap, launcherMap, agencyMap, agencyDocMap, baseMap)
			if !include {
				continue
			}
			launches = append(launches, launchSummary)
		}

		basic := buildPublicCompanyListItem(agency, doc)
		return models.PublicCompanyView{
			ID:           basic.ID,
			Name:         basic.Name,
			Description:  basic.Description,
			Founded:      basic.Founded,
			Founder:      basic.Founder,
			Headquarters: basic.Headquarters,
			Employees:    basic.Employees,
			Website:      basic.Website,
			ImageURL:     basic.ImageURL,
			Rockets:      rockets,
			Launches:     launches,
			Stats: models.PublicCompanyStats{
				RocketCount:        len(rockets),
				LaunchCount:        doc.TotalLaunchCount,
				SuccessfulLaunches: doc.SuccessfulLaunches,
				FailedLaunches:     doc.FailedLaunches,
				PendingLaunches:    doc.PendingLaunches,
			},
		}, nil
	})
}
