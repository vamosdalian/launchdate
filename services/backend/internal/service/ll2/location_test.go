package ll2

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vamosdalian/launchdate-backend/internal/models"
)

func TestLL2Service_LocationDB(t *testing.T) {
	s := newTestLL2Service()
	clearCollections(t, LL2LOCATION, LL2PAD)

	t.Run("Locations", func(t *testing.T) {
		locations := []*models.LL2LocationSerializerWithPads{
			{
				LL2Location: models.LL2Location{
					ID:   1,
					Name: "Location One",
				},
			},
			{
				LL2Location: models.LL2Location{
					ID:   2,
					Name: "Location Two",
				},
			},
		}

		err := s.SaveLocationsToDB(locations)
		assert.NoError(t, err)

		list, err := s.GetLocationsFromDB(10, 0)
		assert.NoError(t, err)
		assert.Equal(t, 2, list.Count)
		assert.Len(t, list.Locations, 2)
		assert.Equal(t, 1, list.Locations[0].ID)

		single, err := s.GetLocationByIDFromDB(2)
		assert.NoError(t, err)
		assert.Equal(t, "Location Two", single.Name)

		_, err = s.GetLocationByIDFromDB(999)
		assert.Error(t, err)
	})

	t.Run("Pads", func(t *testing.T) {
		pads := []*models.LL2Pad{
			{
				LL2PadSerializerNoLocation: models.LL2PadSerializerNoLocation{
					Id:   100,
					Name: "Pad One",
				},
			},
			{
				LL2PadSerializerNoLocation: models.LL2PadSerializerNoLocation{
					Id:   200,
					Name: "Pad Two",
				},
			},
		}

		err := s.SavePadsToDB(pads)
		assert.NoError(t, err)

		list, err := s.GetPadsFromDB(10, 0)
		assert.NoError(t, err)
		assert.Equal(t, 2, list.Count)
		assert.Len(t, list.Pads, 2)
		assert.Equal(t, 100, list.Pads[0].Id)

		single, err := s.GetPadByIDFromDB(200)
		assert.NoError(t, err)
		assert.Equal(t, "Pad Two", single.Name)

		_, err = s.GetPadByIDFromDB(999)
		assert.Error(t, err)
	})

	t.Run("GetLocationsFromApi", func(t *testing.T) {
		resp, err := s.GetLocationsFromApi(10, 0)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Greater(t, resp.Count, 0)
		assert.NotEmpty(t, resp.Results)
		assert.Equal(t, 9, resp.Results[0].ID)
	})

	t.Run("GetPadsFromAPI", func(t *testing.T) {
		resp, err := s.GetPadsFromAPI(10, 0)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Greater(t, resp.Count, 0)
		assert.NotEmpty(t, resp.Results)
		assert.Equal(t, 63, resp.Results[0].Id)
	})

	t.Run("GetLocationFromAPI", func(t *testing.T) {
		loc, err := s.GetLocationFromAPI(9)
		assert.NoError(t, err)
		assert.NotNil(t, loc)
		assert.Equal(t, 9, loc.ID)
		assert.Equal(t, "Naro Space Center, South Korea", loc.Name)
	})

	t.Run("GetPadFromAPI", func(t *testing.T) {
		pad, err := s.GetPadFromAPI(63)
		assert.NoError(t, err)
		assert.NotNil(t, pad)
		assert.Equal(t, 63, pad.Id)
		assert.True(t, pad.Active)
	})
}
