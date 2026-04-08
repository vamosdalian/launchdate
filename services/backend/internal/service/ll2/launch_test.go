package ll2

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vamosdalian/launchdate-backend/internal/models"
)

func TestLL2Service_LaunchDB(t *testing.T) {
	s := newTestLL2Service()
	clearCollections(t, LL2LAUNCH)

	t.Run("SaveLaunchesToDB and GetLaunchesFromDB", func(t *testing.T) {
		launches := []*models.LL2LaunchDetailed{
			{
				LL2LaunchNormal: models.LL2LaunchNormal{
					LL2LaunchBasic: models.LL2LaunchBasic{
						ID:   "launch-1",
						Name: "Launch One",
						Net:  "2023-01-01T00:00:00Z",
					},
				},
			},
			{
				LL2LaunchNormal: models.LL2LaunchNormal{
					LL2LaunchBasic: models.LL2LaunchBasic{
						ID:   "launch-2",
						Name: "Launch Two",
						Net:  "2023-01-02T00:00:00Z",
					},
				},
			},
		}

		err := s.SaveLaunchesToDB(launches)
		assert.NoError(t, err)

		list, err := s.GetLaunchesFromDB(10, 0)
		assert.NoError(t, err)
		assert.Equal(t, 2, list.Count)
		assert.Len(t, list.Launches, 2)
		assert.Equal(t, "launch-1", list.Launches[0].ID)

		single, err := s.GetLaunchByIDFromDB("launch-2")
		assert.NoError(t, err)
		assert.Equal(t, "Launch Two", single.Name)

		_, err = s.GetLaunchByIDFromDB("launch-999")
		assert.Error(t, err)
	})

	t.Run("GetLaunchesFromAPI", func(t *testing.T) {
		resp, err := s.GetLaunchesFromAPI(10, 0)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Greater(t, resp.Count, 0)
		assert.NotEmpty(t, resp.Results)
		assert.Equal(t, "6f18e0ac-ef45-4658-89b1-1cbf48c821ae", resp.Results[0].ID)
	})

	t.Run("GetLaunchFromAPI", func(t *testing.T) {
		launch, err := s.GetLaunchFromAPI("6f18e0ac-ef45-4658-89b1-1cbf48c821ae")
		assert.NoError(t, err)
		assert.NotNil(t, launch)
		assert.Equal(t, "6f18e0ac-ef45-4658-89b1-1cbf48c821ae", launch.ID)
		assert.Equal(t, "Falcon 9 Block 5 | Starlink Group 11-4", launch.Name)
	})
}
