package ll2

import (
	"context"
	"fmt"

	"github.com/vamosdalian/launchdate-backend/internal/models"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const LL2LAUNCHER = "ll2_launcher"
const LL2LAUNCHERFAMILY = "ll2_launcher_family"

func (s *LL2Service) GetLaunchersFromAPI(limit, offset int) (*models.LL2LauncherResponse, error) {
	var launches *models.LL2LauncherResponse
	err := s.GetDataFromAPI("launcher_configurations", limit, offset, &launches)
	return launches, err
}

func (s *LL2Service) GetLauncherFromAPI(id int) (*models.LL2LauncherConfigDetailed, error) {
	var launcher *models.LL2LauncherConfigDetailed
	err := s.GetDataFromAPI(fmt.Sprintf("launcher_configurations/%d", id), 1, 0, &launcher)
	return launcher, err
}

func (s *LL2Service) SaveLaunchersToDB(launchers []models.LL2LauncherConfigNormal) error {
	for _, launcher := range launchers {
		filter := map[string]any{
			"id": launcher.ID,
		}
		update := map[string]any{
			"$set": launcher,
		}
		opts := options.Update().SetUpsert(true)
		_, err := s.mongoClient.Collection(LL2LAUNCHER).UpdateOne(context.Background(), filter, update, opts)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *LL2Service) GetLaunchersFromDB(limit, offset int) (models.LL2LauncherList, error) {
	collection := s.mongoClient.Collection(LL2LAUNCHER)
	ctx := context.Background()

	count, err := collection.EstimatedDocumentCount(ctx)
	if err != nil {
		return models.LL2LauncherList{}, err
	}

	opts := options.Find()
	opts.SetLimit(int64(limit))
	opts.SetSkip(int64(offset))
	opts.SetSort(map[string]int{"id": 1})
	cursor, err := collection.Find(ctx, struct{}{}, opts)
	if err != nil {
		return models.LL2LauncherList{}, err
	}
	defer cursor.Close(ctx)

	var launchers []models.LL2LauncherConfigNormal
	for cursor.Next(ctx) {
		var launcher models.LL2LauncherConfigNormal
		if err := cursor.Decode(&launcher); err != nil {
			return models.LL2LauncherList{}, err
		}
		launchers = append(launchers, launcher)
	}

	if err := cursor.Err(); err != nil {
		return models.LL2LauncherList{}, err
	}

	return models.LL2LauncherList{Count: int(count), Launchers: launchers}, nil
}

func (s *LL2Service) GetLauncherByIDFromDB(id int) (*models.LL2LauncherConfigDetailed, error) {
	collection := s.mongoClient.Collection(LL2LAUNCHER)
	ctx := context.Background()

	var launcher models.LL2LauncherConfigDetailed
	if err := collection.FindOne(ctx, map[string]any{"id": id}).Decode(&launcher); err != nil {
		return nil, err
	}

	return &launcher, nil
}

func (s *LL2Service) GetLauncherFamilyFromAPI(id int) (*models.LL2LauncherConfigFamilyDetailed, error) {
	var family *models.LL2LauncherConfigFamilyDetailed
	err := s.GetDataFromAPI(fmt.Sprintf("launcher_configuration_families/%d", id), 1, 0, &family)
	return family, err
}

func (s *LL2Service) GetLauncherFamiliesFromAPI(limit, offset int) (*models.LL2LauncherFamilyResponse, error) {
	var families *models.LL2LauncherFamilyResponse
	err := s.GetDataFromAPI("launcher_configuration_families", limit, offset, &families)
	return families, err
}

func (s *LL2Service) SaveLauncherFamiliesToDB(families []models.LL2LauncherConfigFamilyDetailed) error {
	for _, family := range families {
		filter := map[string]any{
			"id": family.ID,
		}
		update := map[string]any{
			"$set": family,
		}
		opts := options.Update().SetUpsert(true)
		_, err := s.mongoClient.Collection(LL2LAUNCHERFAMILY).UpdateOne(context.Background(), filter, update, opts)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *LL2Service) GetLauncherFamiliesFromDB(limit, offset int) (models.LL2LauncherFamilyList, error) {
	collection := s.mongoClient.Collection("ll2_launcher_family")
	ctx := context.Background()

	count, err := collection.EstimatedDocumentCount(ctx)
	if err != nil {
		return models.LL2LauncherFamilyList{}, err
	}

	opts := options.Find()
	opts.SetLimit(int64(limit))
	opts.SetSkip(int64(offset))
	opts.SetSort(map[string]int{"id": 1})
	cursor, err := collection.Find(ctx, struct{}{}, opts)
	if err != nil {
		return models.LL2LauncherFamilyList{}, err
	}
	defer cursor.Close(ctx)

	var families []models.LL2LauncherConfigFamilyDetailed
	for cursor.Next(ctx) {
		var family models.LL2LauncherConfigFamilyDetailed
		if err := cursor.Decode(&family); err != nil {
			return models.LL2LauncherFamilyList{}, err
		}
		families = append(families, family)
	}

	if err := cursor.Err(); err != nil {
		return models.LL2LauncherFamilyList{}, err
	}

	return models.LL2LauncherFamilyList{Count: int(count), Families: families}, nil
}

func (s *LL2Service) GetLauncherFamilyByIDFromDB(id int) (*models.LL2LauncherConfigFamilyDetailed, error) {
	collection := s.mongoClient.Collection(LL2LAUNCHERFAMILY)
	ctx := context.Background()

	var family models.LL2LauncherConfigFamilyDetailed
	if err := collection.FindOne(ctx, map[string]any{"id": id}).Decode(&family); err != nil {
		return nil, err
	}

	return &family, nil
}
