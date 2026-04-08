package ll2

import (
	"context"

	"github.com/vamosdalian/launchdate-backend/internal/models"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const LL2LAUNCH = "ll2_launch"

func (s *LL2Service) GetLaunchesFromAPI(limit, offset int) (*models.LL2Response, error) {
	var launches *models.LL2Response
	err := s.GetDataFromAPI("launches", limit, offset, &launches, "ordering=net")

	return launches, err
}

func (s *LL2Service) GetLaunchFromAPI(launchId string) (*models.LL2LaunchDetailed, error) {
	var launch *models.LL2LaunchDetailed
	err := s.GetDataFromAPI("launches/"+launchId, 1, 0, &launch)

	return launch, err
}

func (s *LL2Service) GetLaunchUpComingFromAPI(launchId string, limit, offset int) (*models.LL2Response, error) {
	var launches *models.LL2Response
	err := s.GetDataFromAPI("launches/upcoming/"+launchId, limit, offset, &launches, "ordering=net")

	return launches, err
}

func (s *LL2Service) SaveLaunchesToDB(launches []*models.LL2LaunchDetailed) error {
	for _, launch := range launches {
		filter := map[string]any{
			"id": launch.ID,
		}
		update := map[string]any{
			"$set": launch,
		}
		opts := options.Update().SetUpsert(true)
		_, err := s.mongoClient.Collection(LL2LAUNCH).UpdateOne(context.Background(),
			filter, update, opts)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *LL2Service) GetLaunchesFromDB(limit, offset int) (models.LL2LaunchList, error) {
	collection := s.mongoClient.Collection(LL2LAUNCH)

	count, err := collection.EstimatedDocumentCount(context.Background())
	if err != nil {
		return models.LL2LaunchList{}, err
	}

	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(offset))
	findOptions.SetSort(map[string]int{"net": 1}) // Sort by net ascending

	cursor, err := collection.Find(context.Background(), map[string]any{}, findOptions)
	if err != nil {
		return models.LL2LaunchList{}, err
	}
	defer cursor.Close(context.Background())

	var launches []models.LL2LaunchNormal
	for cursor.Next(context.Background()) {
		var launch models.LL2LaunchNormal
		if err := cursor.Decode(&launch); err != nil {
			return models.LL2LaunchList{}, err
		}
		launches = append(launches, launch)
	}

	if err := cursor.Err(); err != nil {
		return models.LL2LaunchList{}, err
	}

	return models.LL2LaunchList{Count: int(count), Launches: launches}, nil
}

func (s *LL2Service) GetLaunchByIDFromDB(id string) (*models.LL2LaunchDetailed, error) {
	collection := s.mongoClient.Collection(LL2LAUNCH)

	var launch models.LL2LaunchDetailed
	err := collection.FindOne(context.Background(), map[string]any{"id": id}).Decode(&launch)
	if err != nil {
		return nil, err
	}

	return &launch, nil
}

func (s *LL2Service) GetLatestLaunchFromDB() (*models.LL2LaunchDetailed, error) {
	collection := s.mongoClient.Collection(LL2LAUNCH)

	findOptions := options.FindOne()
	findOptions.SetSort(map[string]int{"net": -1})

	var launch models.LL2LaunchDetailed
	err := collection.FindOne(context.Background(), map[string]any{}, findOptions).Decode(&launch)
	if err != nil {
		return nil, err
	}

	return &launch, nil
}
