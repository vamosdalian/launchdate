package ll2

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vamosdalian/launchdate-backend/internal/models"
)

func TestLL2Service_LauncherDB(t *testing.T) {
	s := newTestLL2Service()
	clearCollections(t, LL2LAUNCHER, LL2LAUNCHERFAMILY)

	t.Run("Launchers", func(t *testing.T) {
		launchers := []models.LL2LauncherConfigDetailed{
			{
				LL2LauncherConfigNormal: models.LL2LauncherConfigNormal{
					LL2LauncherConfigList: models.LL2LauncherConfigList{
						ID:   1,
						Name: "Launcher One",
					},
				},
				Description: "Launcher One Description",
				Length:      70,
				Diameter:    5.2,
				LaunchMass:  549054,
				LeoCapacity: 22800,
				LaunchCost:  67000000,
				ToThrust:    7607,
			},
			{
				LL2LauncherConfigNormal: models.LL2LauncherConfigNormal{
					LL2LauncherConfigList: models.LL2LauncherConfigList{
						ID:   2,
						Name: "Launcher Two",
					},
				},
				Description: "Launcher Two Description",
				Length:      57,
				Diameter:    3.7,
			},
		}

		err := s.SaveLaunchersToDB(launchers)
		assert.NoError(t, err)

		list, err := s.GetLaunchersFromDB(10, 0)
		assert.NoError(t, err)
		assert.Equal(t, 2, list.Count)
		assert.Len(t, list.Launchers, 2)
		assert.Equal(t, 1, list.Launchers[0].ID)
		assert.Equal(t, 2, list.Launchers[1].ID)

		single, err := s.GetLauncherByIDFromDB(1)
		assert.NoError(t, err)
		assert.Equal(t, "Launcher One", single.Name)
		assert.Equal(t, "Launcher One Description", single.Description)
		assert.Equal(t, 70.0, single.Length)
		assert.Equal(t, 5.2, single.Diameter)
		assert.Equal(t, 549054.0, single.LaunchMass)
		assert.Equal(t, 22800.0, single.LeoCapacity)
		assert.Equal(t, 67000000, single.LaunchCost)
		assert.Equal(t, 7607.0, single.ToThrust)

		_, err = s.GetLauncherByIDFromDB(999)
		assert.Error(t, err)
	})

	t.Run("LauncherFamilies", func(t *testing.T) {
		families := []models.LL2LauncherConfigFamilyDetailed{
			{
				LL2LauncherConfigFamilyNormal: models.LL2LauncherConfigFamilyNormal{
					LL2LauncherConfigFamilyMini: models.LL2LauncherConfigFamilyMini{
						ID:   10,
						Name: "Family One",
					},
				},
			},
			{
				LL2LauncherConfigFamilyNormal: models.LL2LauncherConfigFamilyNormal{
					LL2LauncherConfigFamilyMini: models.LL2LauncherConfigFamilyMini{
						ID:   20,
						Name: "Family Two",
					},
				},
			},
		}

		err := s.SaveLauncherFamiliesToDB(families)
		assert.NoError(t, err)

		list, err := s.GetLauncherFamiliesFromDB(10, 0)
		assert.NoError(t, err)
		assert.Equal(t, 2, list.Count)
		assert.Len(t, list.Families, 2)
		assert.Equal(t, 10, list.Families[0].ID)

		single, err := s.GetLauncherFamilyByIDFromDB(20)
		assert.NoError(t, err)
		assert.Equal(t, "Family Two", single.Name)

		_, err = s.GetLauncherFamilyByIDFromDB(999)
		assert.Error(t, err)
	})

	t.Run("GetLaunchersFromAPI", func(t *testing.T) {
		resp, err := s.GetLaunchersFromAPI(10, 0)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Greater(t, resp.Count, 0)
		assert.NotEmpty(t, resp.Results)
		assert.Equal(t, 136, resp.Results[0].ID)
	})

	t.Run("GetLauncherFamiliesFromAPI", func(t *testing.T) {
		resp, err := s.GetLauncherFamiliesFromAPI(10, 0)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Greater(t, resp.Count, 0)
		assert.NotEmpty(t, resp.Results)
		assert.Equal(t, 56, resp.Results[0].ID)
	})

	t.Run("GetLauncherFromAPI", func(t *testing.T) {
		launcher, err := s.GetLauncherFromAPI(136)
		assert.NoError(t, err)
		assert.NotNil(t, launcher)
		assert.Equal(t, 136, launcher.ID)
		assert.Equal(t, "Angara 1.2", launcher.Name)
	})

	t.Run("GetLauncherFamilyFromAPI", func(t *testing.T) {
		family, err := s.GetLauncherFamilyFromAPI(56)
		assert.NoError(t, err)
		assert.NotNil(t, family)
		assert.Equal(t, 56, family.ID)
		assert.Equal(t, "Angara", family.Name)
	})
}
