package core

import (
	"context"
	"strconv"
	"strings"

	"github.com/vamosdalian/launchdate-backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func publicID(id int64) string {
	if id == 0 {
		return ""
	}
	return strconv.FormatInt(id, 10)
}

func normalizePublicLaunchStatus(status models.LL2Status) (models.PublicLaunchStatus, string) {
	label := strings.TrimSpace(status.Abbrev)
	if label == "" {
		label = strings.TrimSpace(status.Name)
	}
	if label == "" {
		label = "Unknown"
	}

	normalized := strings.ToLower(strings.TrimSpace(status.Name + " " + status.Abbrev))
	switch {
	case strings.Contains(normalized, "success"):
		return models.PublicLaunchStatusSuccess, label
	case strings.Contains(normalized, "partial") && strings.Contains(normalized, "failure"):
		return models.PublicLaunchStatusFailure, label
	case strings.Contains(normalized, "failure"):
		return models.PublicLaunchStatusFailure, label
	case strings.Contains(normalized, "cancel"):
		return models.PublicLaunchStatusCancelled, label
	case strings.Contains(normalized, "hold"), strings.Contains(normalized, "delay"):
		return models.PublicLaunchStatusDelayed, label
	case strings.Contains(normalized, "flight"):
		return models.PublicLaunchStatusInFlight, label
	case strings.Contains(normalized, "tbd"), strings.Contains(normalized, "determined"), strings.Contains(normalized, "confirmed"), strings.Contains(normalized, "go"):
		return models.PublicLaunchStatusScheduled, label
	}

	switch status.ID {
	case 1, 2, 8:
		return models.PublicLaunchStatusScheduled, label
	case 3:
		return models.PublicLaunchStatusSuccess, label
	case 4, 7:
		return models.PublicLaunchStatusFailure, label
	case 5:
		return models.PublicLaunchStatusDelayed, label
	case 6:
		return models.PublicLaunchStatusInFlight, label
	case 9:
		return models.PublicLaunchStatusCancelled, label
	default:
		return models.PublicLaunchStatusUnknown, label
	}
}

func (m *MainService) loadLaunchMapByExternalIDs(ctx context.Context, externalIDs []string) (map[string]models.Launch, error) {
	launches := make(map[string]models.Launch)
	if len(externalIDs) == 0 {
		return launches, nil
	}

	cursor, err := m.mc.Collection(COLLECTION_LAUNCH).Find(ctx, bson.M{"external_id": bson.M{"$in": externalIDs}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var launch models.Launch
		if err := cursor.Decode(&launch); err != nil {
			return nil, err
		}
		launches[launch.ExternalID] = launch
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return launches, nil
}

func (m *MainService) loadAgencyMapByExternalIDs(ctx context.Context, externalIDs []int64) (map[int64]models.Agency, error) {
	agencies := make(map[int64]models.Agency)
	if len(externalIDs) == 0 {
		return agencies, nil
	}

	cursor, err := m.mc.Collection(COLLECTION_AGENCY).Find(ctx, bson.M{"external_id": bson.M{"$in": externalIDs}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var agency models.Agency
		if err := cursor.Decode(&agency); err != nil {
			return nil, err
		}
		agencies[agency.ExternalID] = agency
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return agencies, nil
}

func (m *MainService) loadLL2AgencyMapByExternalIDs(ctx context.Context, externalIDs []int64) (map[int64]models.LL2AgencyDetailed, error) {
	agencies := make(map[int64]models.LL2AgencyDetailed)
	if len(externalIDs) == 0 {
		return agencies, nil
	}

	ids := make([]int, 0, len(externalIDs))
	seen := make(map[int64]struct{}, len(externalIDs))
	for _, externalID := range externalIDs {
		if externalID == 0 {
			continue
		}
		if _, ok := seen[externalID]; ok {
			continue
		}
		seen[externalID] = struct{}{}
		ids = append(ids, int(externalID))
	}

	if len(ids) == 0 {
		return agencies, nil
	}

	cursor, err := m.mc.Collection(COLLECTION_LL2_AGENCY).Find(ctx, bson.M{"id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var agency models.LL2AgencyDetailed
		if err := cursor.Decode(&agency); err != nil {
			return nil, err
		}
		agencies[int64(agency.ID)] = agency
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return agencies, nil
}

func (m *MainService) loadLaunchBaseMapByExternalIDs(ctx context.Context, externalIDs []int64) (map[int64]models.LaunchBase, error) {
	bases := make(map[int64]models.LaunchBase)
	if len(externalIDs) == 0 {
		return bases, nil
	}

	cursor, err := m.mc.Collection(COLLECTION_LAUNCH_BASE).Find(ctx, bson.M{"external_id": bson.M{"$in": externalIDs}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var base models.LaunchBase
		if err := cursor.Decode(&base); err != nil {
			return nil, err
		}
		bases[base.ExternalID] = base
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return bases, nil
}

func (m *MainService) loadRocketMapByExternalIDs(ctx context.Context, externalIDs []int64) (map[int64]models.Rocket, error) {
	rockets := make(map[int64]models.Rocket)
	if len(externalIDs) == 0 {
		return rockets, nil
	}

	cursor, err := m.mc.Collection(COLLECTION_ROCKET).Find(ctx, bson.M{"external_id": bson.M{"$in": externalIDs}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var rocket models.Rocket
		if err := cursor.Decode(&rocket); err != nil {
			return nil, err
		}
		rockets[rocket.ExternalID] = rocket
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return rockets, nil
}

func (m *MainService) loadLauncherConfigMapByExternalIDs(ctx context.Context, externalIDs []int64) (map[int64]models.LL2LauncherConfigNormal, error) {
	launchers := make(map[int64]models.LL2LauncherConfigNormal)
	if len(externalIDs) == 0 {
		return launchers, nil
	}

	ids := make([]int, 0, len(externalIDs))
	seen := make(map[int64]struct{}, len(externalIDs))
	for _, externalID := range externalIDs {
		if externalID == 0 {
			continue
		}
		if _, ok := seen[externalID]; ok {
			continue
		}
		seen[externalID] = struct{}{}
		ids = append(ids, int(externalID))
	}

	if len(ids) == 0 {
		return launchers, nil
	}

	cursor, err := m.mc.Collection(COLLECTION_LL2_LAUNCHER).Find(ctx, bson.M{"id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var launcher models.LL2LauncherConfigDetailed
		if err := cursor.Decode(&launcher); err != nil {
			return nil, err
		}
		launchers[int64(launcher.ID)] = launcher.LL2LauncherConfigNormal
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return launchers, nil
}

func resolveLaunchThumb(launch models.Launch, image models.LL2Image) string {
	return firstNonEmpty(
		launch.ThumbImage,
		image.ThumbnailURL,
		image.ImageURL,
		launch.BackgroundImage,
		firstImageURL(launch.ImageList),
	)
}

func resolveLaunchBackground(launch models.Launch, image models.LL2Image) string {
	return firstNonEmpty(
		launch.BackgroundImage,
		image.ImageURL,
		image.ThumbnailURL,
		launch.ThumbImage,
		firstImageURL(launch.ImageList),
	)
}

func resolveRocketThumb(rocket models.Rocket, image models.LL2Image) string {
	return firstNonEmpty(
		rocket.ThumbImage,
		rocket.MainImage,
		rocket.LaunchImage,
		image.ThumbnailURL,
		image.ImageURL,
		firstImageURL(rocket.ImageList),
	)
}

func resolveRocketMain(rocket models.Rocket, image models.LL2Image) string {
	return firstNonEmpty(
		rocket.MainImage,
		image.ImageURL,
		image.ThumbnailURL,
		rocket.LaunchImage,
		rocket.ThumbImage,
		firstImageURL(rocket.ImageList),
	)
}

func resolveRocketLaunch(rocket models.Rocket, image models.LL2Image) string {
	return firstNonEmpty(
		rocket.LaunchImage,
		image.ImageURL,
		image.ThumbnailURL,
		rocket.MainImage,
		rocket.ThumbImage,
		firstImageURL(rocket.ImageList),
	)
}

func buildPublicRocketRef(rocket models.Rocket, config models.LL2LauncherConfigList, image models.LL2Image) models.PublicRocketRef {
	return models.PublicRocketRef{
		ID:         publicID(rocket.ID),
		Name:       config.Name,
		ImageURL:   resolveRocketMain(rocket, image),
		ThumbImage: resolveRocketThumb(rocket, image),
	}
}

func firstImageURL(images []string) string {
	for _, image := range images {
		if strings.TrimSpace(image) != "" {
			return image
		}
	}

	return ""
}

func buildPublicCompanyRef(agency models.Agency, name string, doc *models.LL2AgencyDetailed) models.PublicCompanyRef {
	return models.PublicCompanyRef{
		ID:       publicID(agency.ID),
		Name:     name,
		ImageURL: resolveAgencyLogo(agency, doc),
	}
}

func resolveAgencyLogo(agency models.Agency, doc *models.LL2AgencyDetailed) string {
	if doc == nil {
		return firstNonEmpty(agency.ThumbImage, firstImageURL(agency.Images))
	}

	return firstNonEmpty(
		agency.ThumbImage,
		doc.Logo.ImageURL,
		doc.Logo.ThumbnailURL,
		doc.Image.ImageURL,
		doc.Image.ThumbnailURL,
		doc.SocialLogo.ImageURL,
		doc.SocialLogo.ThumbnailURL,
		firstImageURL(agency.Images),
	)
}

func publicCompanyImageURL(agency models.Agency, doc *models.LL2AgencyDetailed) string {
	return resolveAgencyLogo(agency, doc)
}

func buildPublicLaunchBaseRef(base models.LaunchBase, location models.LL2Location) models.PublicLaunchBaseRef {
	return models.PublicLaunchBaseRef{
		ID:        publicID(base.ID),
		Name:      location.Name,
		Location:  firstNonEmpty(location.TimezoneName, location.CelestialBody.Name, location.Country.Name, location.Name),
		Country:   location.Country.Name,
		Latitude:  location.Latitude,
		Longitude: location.Longitude,
	}
}

func (m *MainService) buildPublicLaunchSummary(
	internalLaunch models.Launch,
	ll2Launch models.LL2LaunchDetailed,
	rocketMap map[int64]models.Rocket,
	launcherMap map[int64]models.LL2LauncherConfigNormal,
	agencyMap map[int64]models.Agency,
	agencyDocMap map[int64]models.LL2AgencyDetailed,
	baseMap map[int64]models.LaunchBase,
) (models.PublicLaunchSummary, bool) {
	if internalLaunch.ID == 0 {
		return models.PublicLaunchSummary{}, false
	}

	status, statusLabel := normalizePublicLaunchStatus(ll2Launch.Status)
	launcherConfig := launcherMap[int64(ll2Launch.Rocket.Configuration.ID)]
	rocket := buildPublicRocketRef(
		rocketMap[int64(ll2Launch.Rocket.Configuration.ID)],
		ll2Launch.Rocket.Configuration,
		launcherConfig.Image,
	)
	companyDoc, hasCompanyDoc := agencyDocMap[int64(ll2Launch.LaunchServiceProvider.ID)]
	company := buildPublicCompanyRef(
		agencyMap[int64(ll2Launch.LaunchServiceProvider.ID)],
		ll2Launch.LaunchServiceProvider.Name,
		func() *models.LL2AgencyDetailed {
			if !hasCompanyDoc {
				return nil
			}
			return &companyDoc
		}(),
	)
	launchBase := buildPublicLaunchBaseRef(baseMap[int64(ll2Launch.Pad.Location.ID)], ll2Launch.Pad.Location)

	return models.PublicLaunchSummary{
		ID:              publicID(internalLaunch.ID),
		Name:            ll2Launch.Name,
		LaunchTime:      ll2Launch.Net,
		Status:          status,
		StatusLabel:     statusLabel,
		ThumbImage:      resolveLaunchThumb(internalLaunch, ll2Launch.Image),
		BackgroundImage: resolveLaunchBackground(internalLaunch, ll2Launch.Image),
		Rocket:          rocket,
		Company:         company,
		LaunchBase:      launchBase,
	}, true
}

func buildPublicRocketListItem(rocket models.Rocket, config models.LL2LauncherConfigNormal) models.PublicRocketListItem {
	return models.PublicRocketListItem{
		ID:         publicID(rocket.ID),
		Name:       config.Name,
		ThumbImage: resolveRocketThumb(rocket, config.Image),
	}
}

func publicRocketListImageURL(rocket models.Rocket, image models.LL2Image) string {
	return resolveRocketThumb(rocket, image)
}

func buildPublicCompanyListItem(agency models.Agency, doc models.LL2AgencyDetailed) models.PublicCompanyListItem {
	return models.PublicCompanyListItem{
		ID:           publicID(agency.ID),
		Name:         doc.Name,
		Description:  doc.Description,
		Founded:      doc.FoundingYear,
		Founder:      firstNonEmpty(doc.Administrator, "Unknown"),
		Headquarters: firstNonEmpty(firstAgencyCountryName(doc.Country), "Unknown"),
		Employees:    0,
		Website:      firstNonEmpty(firstSocialURL(agency.SocialUrl), doc.InfoURL, firstDetailedSocialURL(doc.SocialMediaLinks), doc.WikiURL),
		ImageURL:     resolveAgencyLogo(agency, &doc),
	}
}

func buildPublicLaunchBaseListItem(base models.LaunchBase, doc models.LL2LocationSerializerWithPads) models.PublicLaunchBaseListItem {
	return models.PublicLaunchBaseListItem{
		ID:          publicID(base.ID),
		Name:        doc.Name,
		Location:    firstNonEmpty(doc.TimezoneName, doc.CelestialBody.Name, doc.Country.Name, doc.Name),
		Country:     doc.Country.Name,
		Description: doc.Description,
		ImageURL:    firstNonEmpty(doc.Image.ImageURL, doc.Image.ThumbnailURL, doc.MapImage),
		Latitude:    doc.Latitude,
		Longitude:   doc.Longitude,
	}
}

func findExistingRocket(rockets map[int64]models.Rocket, externalID int64) models.Rocket {
	if rocket, ok := rockets[externalID]; ok {
		return rocket
	}
	return models.Rocket{}
}

func findExistingAgency(agencies map[int64]models.Agency, externalID int64) models.Agency {
	if agency, ok := agencies[externalID]; ok {
		return agency
	}
	return models.Agency{}
}

func findExistingLaunchBase(bases map[int64]models.LaunchBase, externalID int64) models.LaunchBase {
	if base, ok := bases[externalID]; ok {
		return base
	}
	return models.LaunchBase{}
}

func (m *MainService) findRocketByCoreID(ctx context.Context, id int64) (models.Rocket, error) {
	var rocket models.Rocket
	err := m.mc.Collection(COLLECTION_ROCKET).FindOne(ctx, bson.M{"id": id}).Decode(&rocket)
	if err == nil {
		return rocket, nil
	}
	if err != mongo.ErrNoDocuments {
		return models.Rocket{}, err
	}

	return models.Rocket{}, mongo.ErrNoDocuments
}
