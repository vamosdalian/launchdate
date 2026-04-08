package ll2

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vamosdalian/launchdate-backend/internal/models"
)

func TestLL2Service_AgencyDB(t *testing.T) {
	s := newTestLL2Service()
	clearCollections(t, LL2AGENCY)

	t.Run("SaveAgenciesToDB and GetAngecyFromDB", func(t *testing.T) {
		agencies := []*models.LL2AgencyDetailed{
			{
				LL2AgencyNormal: models.LL2AgencyNormal{
					LL2AgencyMini: models.LL2AgencyMini{
						ID:   1,
						Name: "Agency One",
					},
				},
				TotalLaunchCount: 10,
			},
			{
				LL2AgencyNormal: models.LL2AgencyNormal{
					LL2AgencyMini: models.LL2AgencyMini{
						ID:   2,
						Name: "Agency Two",
					},
				},
				TotalLaunchCount: 20,
			},
		}

		err := s.SaveAgenciesToDB(agencies)
		assert.NoError(t, err)

		list, err := s.GetAngecyFromDB(10, 0)
		assert.NoError(t, err)
		assert.Equal(t, 2, list.Count)
		assert.Len(t, list.Agencies, 2)
		assert.Equal(t, 1, list.Agencies[0].ID)
		assert.Equal(t, 2, list.Agencies[1].ID)

		// Test Limit and Offset
		listShort, err := s.GetAngecyFromDB(1, 0)
		assert.NoError(t, err)
		assert.Len(t, listShort.Agencies, 1)

		listOffset, err := s.GetAngecyFromDB(1, 1)
		assert.NoError(t, err)
		assert.Len(t, listOffset.Agencies, 1)
		assert.Equal(t, 2, listOffset.Agencies[0].ID)
	})

	t.Run("GetAgencyByIDFromDB", func(t *testing.T) {
		agency, err := s.GetAgencyByIDFromDB(1)
		assert.NoError(t, err)
		assert.NotNil(t, agency)
		assert.Equal(t, "Agency One", agency.Name)

		_, err = s.GetAgencyByIDFromDB(999)
		assert.Error(t, err)
	})

	t.Run("GetAngecyFromAPI", func(t *testing.T) {
		// Mock API is now global and automatically sets URL in newTestLL2Service
		// No manual server start/stop needed here
		resp, err := s.GetAngecyFromAPI(10, 0)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		// Basic validation based on agency.json content
		// We expect results to be populated.
		// Since we mock the file, the limit/offsetparams in GeFromAPI are ignored by our mock server logic,
		// it just returns the file content.
		assert.Greater(t, resp.Count, 0)
		assert.NotEmpty(t, resp.Results)
		assert.Equal(t, 225, resp.Results[0].ID)
	})

	t.Run("GetAgencyFromAPI_Single", func(t *testing.T) {
		agency, err := s.GetAgencyFromAPI(225)
		assert.NoError(t, err)
		assert.NotNil(t, agency)
		assert.Equal(t, 225, agency.ID)
		assert.Equal(t, "1worldspace", agency.Name)
	})
}
