package ll2

import (
	"context"
	"fmt"

	"github.com/vamosdalian/launchdate-backend/internal/models"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const LL2LOCATION = "ll2_location"
const LL2PAD = "ll2_pad"

func (s *LL2Service) GetLocationsFromApi(limit, offset int) (*models.LL2LocationResponse, error) {
	var locations *models.LL2LocationResponse
	err := s.GetDataFromAPI("locations", limit, offset, &locations)
	return locations, err
}

func (s *LL2Service) GetLocationFromAPI(id int) (*models.LL2LocationSerializerWithPads, error) {
	var location *models.LL2LocationSerializerWithPads
	err := s.GetDataFromAPI(fmt.Sprintf("locations/%d", id), 1, 0, &location)
	return location, err
}

func (s *LL2Service) SaveLocationsToDB(locations []*models.LL2LocationSerializerWithPads) error {
	for _, location := range locations {
		filter := map[string]any{
			"id": location.ID,
		}
		update := map[string]any{
			"$set": location,
		}
		opts := options.Update().SetUpsert(true)
		collection := s.mongoClient.Collection(LL2LOCATION)
		_, err := collection.UpdateOne(context.Background(), filter, update, opts)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *LL2Service) GetLocationsFromDB(limit, offset int) (models.LL2LocationList, error) {
	collection := s.mongoClient.Collection(LL2LOCATION)
	ctx := context.Background()

	count, err := collection.EstimatedDocumentCount(ctx)
	if err != nil {
		return models.LL2LocationList{}, err
	}

	opts := options.Find()
	opts.SetLimit(int64(limit))
	opts.SetSkip(int64(offset))
	opts.SetSort(map[string]int{"id": 1})
	cursor, err := collection.Find(ctx, struct{}{}, opts)
	if err != nil {
		return models.LL2LocationList{}, err
	}
	defer cursor.Close(ctx)

	var locations []models.LL2LocationSerializerWithPads
	for cursor.Next(ctx) {
		var location models.LL2LocationSerializerWithPads
		if err := cursor.Decode(&location); err != nil {
			return models.LL2LocationList{}, err
		}
		locations = append(locations, location)
	}

	if err := cursor.Err(); err != nil {
		return models.LL2LocationList{}, err
	}

	return models.LL2LocationList{Count: int(count), Locations: locations}, nil
}

func (s *LL2Service) GetLocationByIDFromDB(id int) (*models.LL2LocationSerializerWithPads, error) {
	collection := s.mongoClient.Collection(LL2LOCATION)
	ctx := context.Background()

	var location models.LL2LocationSerializerWithPads
	err := collection.FindOne(ctx, map[string]any{"id": id}).Decode(&location)
	if err != nil {
		return nil, err
	}

	return &location, nil
}

func (s *LL2Service) GetPadsFromAPI(limit, offset int) (*models.LL2PadResponse, error) {
	var pads *models.LL2PadResponse
	err := s.GetDataFromAPI("pads", limit, offset, &pads)
	return pads, err
}

func (s *LL2Service) GetPadFromAPI(id int) (*models.LL2Pad, error) {
	var pad *models.LL2Pad
	err := s.GetDataFromAPI(fmt.Sprintf("pads/%d", id), 1, 0, &pad)
	return pad, err
}

func (s *LL2Service) SavePadsToDB(pads []*models.LL2Pad) error {
	for _, pad := range pads {
		filter := map[string]any{
			"id": pad.Id,
		}
		update := map[string]any{
			"$set": pad,
		}
		opts := options.Update().SetUpsert(true)
		collection := s.mongoClient.Collection(LL2PAD)
		_, err := collection.UpdateOne(context.Background(), filter, update, opts)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *LL2Service) GetPadsFromDB(limit, offset int) (models.LL2PadList, error) {
	collection := s.mongoClient.Collection(LL2PAD)
	ctx := context.Background()

	count, err := collection.EstimatedDocumentCount(ctx)
	if err != nil {
		return models.LL2PadList{}, err
	}

	opts := options.Find()
	opts.SetLimit(int64(limit))
	opts.SetSkip(int64(offset))
	opts.SetSort(map[string]int{"id": 1})
	cursor, err := collection.Find(ctx, struct{}{}, opts)
	if err != nil {
		return models.LL2PadList{}, err
	}
	defer cursor.Close(ctx)

	var pads []models.LL2Pad
	for cursor.Next(ctx) {
		var pad models.LL2Pad
		if err := cursor.Decode(&pad); err != nil {
			return models.LL2PadList{}, err
		}
		pads = append(pads, pad)
	}

	if err := cursor.Err(); err != nil {
		return models.LL2PadList{}, err
	}

	return models.LL2PadList{Count: int(count), Pads: pads}, nil
}

func (s *LL2Service) GetPadByIDFromDB(id int) (*models.LL2Pad, error) {
	collection := s.mongoClient.Collection(LL2PAD)
	ctx := context.Background()

	var pad models.LL2Pad
	err := collection.FindOne(ctx, map[string]any{"id": id}).Decode(&pad)
	if err != nil {
		return nil, err
	}

	return &pad, nil
}
