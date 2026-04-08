package ll2

import (
	"context"
	"fmt"

	"github.com/vamosdalian/launchdate-backend/internal/models"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const LL2AGENCY = "ll2_agency"

func (s *LL2Service) GetAngecyFromAPI(limit, offset int) (*models.LL2AngecyResponse, error) {
	var launches *models.LL2AngecyResponse
	err := s.GetDataFromAPI("agencies", limit, offset, &launches)
	return launches, err
}

func (s *LL2Service) GetAgencyFromAPI(id int) (*models.LL2AgencyDetailed, error) {
	var agency *models.LL2AgencyDetailed
	err := s.GetDataFromAPI(fmt.Sprintf("agencies/%d", id), 1, 0, &agency)
	return agency, err
}

func (s *LL2Service) SaveAgenciesToDB(agencies []*models.LL2AgencyDetailed) error {
	for _, agency := range agencies {
		filter := map[string]any{
			"id": agency.ID,
		}
		update := map[string]any{
			"$set": agency,
		}
		opts := options.Update().SetUpsert(true)
		_, err := s.mongoClient.Collection(LL2AGENCY).UpdateOne(context.Background(),
			filter, update, opts)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *LL2Service) GetAngecyFromDB(limit, offset int) (models.LL2AgencyList, error) {
	collection := s.mongoClient.Collection(LL2AGENCY)

	count, err := collection.EstimatedDocumentCount(context.Background())
	if err != nil {
		return models.LL2AgencyList{}, err
	}

	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(offset))
	findOptions.SetSort(map[string]int{"id": 1}) // Sort by id ascending

	cursor, err := collection.Find(context.Background(), map[string]any{}, findOptions)
	if err != nil {
		return models.LL2AgencyList{}, err
	}
	defer cursor.Close(context.Background())

	var agencies []models.LL2AgencyDetailed
	for cursor.Next(context.Background()) {
		var agency models.LL2AgencyDetailed
		if err := cursor.Decode(&agency); err != nil {
			return models.LL2AgencyList{}, err
		}
		agencies = append(agencies, agency)
	}

	if err := cursor.Err(); err != nil {
		return models.LL2AgencyList{}, err
	}

	return models.LL2AgencyList{Count: int(count), Agencies: agencies}, nil
}

func (s *LL2Service) GetAgencyByIDFromDB(id int) (*models.LL2AgencyDetailed, error) {
	collection := s.mongoClient.Collection(LL2AGENCY)

	var agency models.LL2AgencyDetailed
	err := collection.FindOne(context.Background(), map[string]any{"id": id}).Decode(&agency)
	if err != nil {
		return nil, err
	}

	return &agency, nil
}
