package core

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"github.com/vamosdalian/launchdate-backend/internal/db"
	"github.com/vamosdalian/launchdate-backend/internal/models"
)

func newTestMainService(t *testing.T) (*MainService, *db.MongoDB, func()) {
	t.Helper()

	ctx := context.Background()
	mongoContainer, err := mongodb.Run(ctx, "mongo:6")
	if err != nil {
		t.Fatal(err)
	}

	uri, err := mongoContainer.ConnectionString(ctx)
	if err != nil {
		t.Fatal(err)
	}

	testDB, cleanupDB, err := db.NewMongoDB(uri, "test_public_lists_db")
	if err != nil {
		t.Fatal(err)
	}

	cleanup := func() {
		cleanupDB()
		if err := mongoContainer.Terminate(ctx); err != nil {
			t.Fatal(err)
		}
	}

	return NewMainService(testDB), testDB, cleanup
}

func clearCoreCollections(t *testing.T, testDB *db.MongoDB, collections ...string) {
	t.Helper()
	for _, collection := range collections {
		err := testDB.Database.Collection(collection).Drop(context.Background())
		if err != nil {
			t.Logf("failed to drop collection %s: %v", collection, err)
		}
	}
}

func TestMainService_PublicRocketListSupportsSearchAndPaging(t *testing.T) {
	service, testDB, cleanup := newTestMainService(t)
	defer cleanup()
	ctx := context.Background()

	clearCoreCollections(t, testDB, COLLECTION_ROCKET, COLLECTION_LL2_LAUNCHER)

	for i := 1; i <= 22; i++ {
		name := fmt.Sprintf("Rocket %02d", i)
		fullName := fmt.Sprintf("Rocket %02d Full", i)
		if i == 1 {
			name = "Falcon 9"
			fullName = "Falcon 9 Block 5"
		}
		if i == 2 {
			name = "Falcon Heavy"
			fullName = "Falcon Heavy"
		}

		_, err := testDB.Collection(COLLECTION_ROCKET).InsertOne(ctx, models.Rocket{
			ID:         int64(1000 + i),
			ExternalID: int64(i),
			ThumbImage: fmt.Sprintf("rocket-%02d.jpg", i),
		})
		assert.NoError(t, err)

		_, err = testDB.Collection(COLLECTION_LL2_LAUNCHER).InsertOne(ctx, models.LL2LauncherConfigDetailed{
			LL2LauncherConfigNormal: models.LL2LauncherConfigNormal{
				LL2LauncherConfigList: models.LL2LauncherConfigList{
					ID:       i,
					Name:     name,
					FullName: fullName,
					Variant:  fmt.Sprintf("Variant %02d", i),
				},
			},
			TotalLaunchCount: 100 - i,
		})
		assert.NoError(t, err)
	}

	searchPage, err := service.GetPublicRockets(0, "falcon")
	assert.NoError(t, err)
	assert.Equal(t, 2, searchPage.Count)
	assert.Len(t, searchPage.Rockets, 2)
	assert.Equal(t, "Falcon 9", searchPage.Rockets[0].Name)
	assert.Equal(t, "Falcon Heavy", searchPage.Rockets[1].Name)

	pagedResult, err := service.GetPublicRockets(1, "")
	assert.NoError(t, err)
	assert.Equal(t, 22, pagedResult.Count)
	assert.Len(t, pagedResult.Rockets, 2)
}

func TestMainService_PublicRocketListImageFallbackOrder(t *testing.T) {
	service, testDB, cleanup := newTestMainService(t)
	defer cleanup()
	ctx := context.Background()

	clearCoreCollections(t, testDB, COLLECTION_ROCKET, COLLECTION_LL2_LAUNCHER)

	_, err := testDB.Collection(COLLECTION_ROCKET).InsertMany(ctx, []interface{}{
		models.Rocket{
			ID:         1001,
			ExternalID: 1,
			ThumbImage: "rocket-thumb.jpg",
			MainImage:  "rocket-main.jpg",
		},
		models.Rocket{
			ID:         1002,
			ExternalID: 2,
			MainImage:  "rocket-main-only.jpg",
		},
		models.Rocket{
			ID:         1003,
			ExternalID: 3,
		},
	})
	assert.NoError(t, err)

	_, err = testDB.Collection(COLLECTION_LL2_LAUNCHER).InsertMany(ctx, []interface{}{
		models.LL2LauncherConfigDetailed{
			LL2LauncherConfigNormal: models.LL2LauncherConfigNormal{
				LL2LauncherConfigList: models.LL2LauncherConfigList{ID: 1, Name: "Rocket One"},
				Image:                 models.LL2Image{ImageURL: "ll2-rocket-1.jpg", ThumbnailURL: "ll2-rocket-1-thumb.jpg"},
			},
			TotalLaunchCount: 300,
		},
		models.LL2LauncherConfigDetailed{
			LL2LauncherConfigNormal: models.LL2LauncherConfigNormal{
				LL2LauncherConfigList: models.LL2LauncherConfigList{ID: 2, Name: "Rocket Two"},
				Image:                 models.LL2Image{ImageURL: "ll2-rocket-2.jpg", ThumbnailURL: "ll2-rocket-2-thumb.jpg"},
			},
			TotalLaunchCount: 200,
		},
		models.LL2LauncherConfigDetailed{
			LL2LauncherConfigNormal: models.LL2LauncherConfigNormal{
				LL2LauncherConfigList: models.LL2LauncherConfigList{ID: 3, Name: "Rocket Three"},
				Image:                 models.LL2Image{ImageURL: "ll2-rocket-3.jpg", ThumbnailURL: "ll2-rocket-3-thumb.jpg"},
			},
			TotalLaunchCount: 100,
		},
	})
	assert.NoError(t, err)

	result, err := service.GetPublicRockets(0, "")
	assert.NoError(t, err)
	assert.Len(t, result.Rockets, 3)
	assert.Equal(t, "rocket-thumb.jpg", result.Rockets[0].ThumbImage)
	assert.Equal(t, "rocket-main-only.jpg", result.Rockets[1].ThumbImage)
	assert.Equal(t, "ll2-rocket-3-thumb.jpg", result.Rockets[2].ThumbImage)
}

func TestMainService_PublicRocketDetailUsesStoredLauncherDetailFields(t *testing.T) {
	service, testDB, cleanup := newTestMainService(t)
	defer cleanup()
	ctx := context.Background()

	clearCoreCollections(t, testDB, COLLECTION_ROCKET, COLLECTION_LL2_LAUNCHER, COLLECTION_LL2_LAUNCH, COLLECTION_AGENCY, COLLECTION_LL2_AGENCY)

	_, err := testDB.Collection(COLLECTION_ROCKET).InsertOne(ctx, models.Rocket{
		ID:         1001,
		ExternalID: 136,
		MainImage:  "core-main.jpg",
	})
	assert.NoError(t, err)

	_, err = testDB.Collection(COLLECTION_LL2_LAUNCHER).InsertOne(ctx, models.LL2LauncherConfigDetailed{
		LL2LauncherConfigNormal: models.LL2LauncherConfigNormal{
			LL2LauncherConfigList: models.LL2LauncherConfigList{
				ID:   136,
				Name: "Falcon 9",
			},
			Active:   true,
			Reusable: true,
			Image:    models.LL2Image{ImageURL: "ll2-main.jpg", ThumbnailURL: "ll2-thumb.jpg"},
			InfoURL:  "https://example.com/falcon9",
			WikiURL:  "https://example.com/wiki/falcon9",
		},
		Description:        "Two-stage partially reusable launch vehicle.",
		Length:             70,
		Diameter:           3.7,
		LaunchCost:         67000000,
		LaunchMass:         549054,
		LeoCapacity:        22800,
		GtoCapacity:        8300,
		GeoCapacity:        4000,
		SsoCapacity:        8340,
		ToThrust:           7607,
		TotalLaunchCount:   400,
		SuccessfulLaunches: 396,
		FailedLaunches:     4,
		AttemptedLandings:  350,
		SuccessfulLandings: 330,
		FailedLandings:     20,
	})
	assert.NoError(t, err)

	rocket, err := service.GetPublicRocket(1001)
	assert.NoError(t, err)
	assert.Equal(t, "Falcon 9", rocket.Name)
	assert.Equal(t, "Two-stage partially reusable launch vehicle.", rocket.Description)
	assert.True(t, rocket.Active)
	assert.True(t, rocket.Reusable)
	assert.Equal(t, "core-main.jpg", rocket.MainImage)
	assert.Equal(t, 70.0, rocket.Length)
	assert.Equal(t, 3.7, rocket.Diameter)
	assert.Equal(t, 67000000.0, rocket.LaunchCost)
	assert.Equal(t, 549054.0, rocket.LaunchMass)
	assert.Equal(t, 22800.0, rocket.LeoCapacity)
	assert.Equal(t, 8300.0, rocket.GtoCapacity)
	assert.Equal(t, 4000.0, rocket.GeoCapacity)
	assert.Equal(t, 8340.0, rocket.SsoCapacity)
	assert.Equal(t, 7607.0, rocket.LiftoffThrust)
	assert.Equal(t, 400, rocket.TotalLaunches)
	assert.Equal(t, 396, rocket.SuccessLaunches)
	assert.Equal(t, 4, rocket.FailureLaunches)
	assert.Equal(t, 350, rocket.TotalLandings)
	assert.Equal(t, 330, rocket.SuccessLandings)
	assert.Equal(t, 20, rocket.FailureLandings)
	assert.Empty(t, rocket.Launches)
}

func TestMainService_PublicCompanyListSupportsSearchAndPaging(t *testing.T) {
	service, testDB, cleanup := newTestMainService(t)
	defer cleanup()
	ctx := context.Background()

	clearCoreCollections(t, testDB, COLLECTION_AGENCY, COLLECTION_LL2_AGENCY)

	for i := 1; i <= 22; i++ {
		name := fmt.Sprintf("Company %02d", i)
		if i == 1 {
			name = "SpaceX"
		}
		if i == 2 {
			name = "Space Forge"
		}

		_, err := testDB.Collection(COLLECTION_AGENCY).InsertOne(ctx, models.Agency{
			ID:         int64(2000 + i),
			ExternalID: int64(i),
			ThumbImage: fmt.Sprintf("company-%02d.jpg", i),
		})
		assert.NoError(t, err)

		_, err = testDB.Collection(COLLECTION_LL2_AGENCY).InsertOne(ctx, models.LL2AgencyDetailed{
			LL2AgencyNormal: models.LL2AgencyNormal{
				LL2AgencyMini: models.LL2AgencyMini{
					ID:   i,
					Name: name,
					Type: models.LL2AgencyType{Name: "Launch Provider"},
				},
				Country:       []models.LL2Country{{Name: "United States"}},
				Description:   fmt.Sprintf("Description %02d", i),
				Administrator: "Admin",
				FoundingYear:  2000 + i,
			},
		})
		assert.NoError(t, err)
	}

	searchPage, err := service.GetPublicCompanies(0, "space", false)
	assert.NoError(t, err)
	assert.Equal(t, 2, searchPage.Count)
	assert.Len(t, searchPage.Companies, 2)
	assert.Equal(t, "Space Forge", searchPage.Companies[0].Name)
	assert.Equal(t, "SpaceX", searchPage.Companies[1].Name)

	pagedResult, err := service.GetPublicCompanies(1, "", false)
	assert.NoError(t, err)
	assert.Equal(t, 22, pagedResult.Count)
	assert.Len(t, pagedResult.Companies, 2)
}

func TestMainService_PublicCompanyListSupportsHomeVisibilityFilter(t *testing.T) {
	service, testDB, cleanup := newTestMainService(t)
	defer cleanup()
	ctx := context.Background()

	clearCoreCollections(t, testDB, COLLECTION_AGENCY, COLLECTION_LL2_AGENCY)

	_, err := testDB.Collection(COLLECTION_AGENCY).InsertMany(ctx, []interface{}{
		models.Agency{
			ID:         9001,
			ExternalID: 101,
			ShowOnHome: true,
			SocialUrl:  []models.SocialUrl{{Name: "Website", URL: "https://www.spacex.com"}},
		},
		models.Agency{ID: 9002, ExternalID: 102, ShowOnHome: false},
		models.Agency{ID: 9003, ExternalID: 103, ShowOnHome: true},
	})
	assert.NoError(t, err)

	_, err = testDB.Collection(COLLECTION_LL2_AGENCY).InsertMany(ctx, []interface{}{
		models.LL2AgencyDetailed{
			LL2AgencyNormal: models.LL2AgencyNormal{
				LL2AgencyMini: models.LL2AgencyMini{ID: 101, Name: "SpaceX"},
				Description:   "Launch provider",
				FoundingYear:  2002,
			},
			InfoURL: "https://ll2.example.com/spacex",
		},
		models.LL2AgencyDetailed{
			LL2AgencyNormal: models.LL2AgencyNormal{
				LL2AgencyMini: models.LL2AgencyMini{ID: 102, Name: "Hidden Company"},
				Description:   "Hidden from home",
				FoundingYear:  2010,
			},
		},
		models.LL2AgencyDetailed{
			LL2AgencyNormal: models.LL2AgencyNormal{
				LL2AgencyMini: models.LL2AgencyMini{ID: 103, Name: "Rocket Lab"},
				Description:   "Featured provider",
				FoundingYear:  2006,
			},
			InfoURL: "https://ll2.example.com/rocket-lab",
		},
	})
	assert.NoError(t, err)

	result, err := service.GetPublicCompanies(0, "", true)
	assert.NoError(t, err)
	assert.Equal(t, 2, result.Count)
	assert.Len(t, result.Companies, 2)
	assert.Equal(t, "Rocket Lab", result.Companies[0].Name)
	assert.Equal(t, "SpaceX", result.Companies[1].Name)
	assert.Equal(t, "https://ll2.example.com/rocket-lab", result.Companies[0].Website)
	assert.Equal(t, "https://www.spacex.com", result.Companies[1].Website)
}

func TestMainService_PublicLaunchBaseListSupportsSearchAndPaging(t *testing.T) {
	service, testDB, cleanup := newTestMainService(t)
	defer cleanup()
	ctx := context.Background()

	clearCoreCollections(t, testDB, COLLECTION_LAUNCH_BASE, COLLECTION_LL2_LOCATION)

	for i := 1; i <= 22; i++ {
		name := fmt.Sprintf("Site %02d", i)
		country := "USA"
		if i == 1 {
			name = "Cape Canaveral"
		}
		if i == 2 {
			name = "Cape Test Range"
		}
		if i > 2 {
			country = "Canada"
		}

		_, err := testDB.Collection(COLLECTION_LAUNCH_BASE).InsertOne(ctx, models.LaunchBase{
			ID:         int64(3000 + i),
			ExternalID: int64(i),
		})
		assert.NoError(t, err)

		_, err = testDB.Collection(COLLECTION_LL2_LOCATION).InsertOne(ctx, models.LL2LocationSerializerWithPads{
			LL2Location: models.LL2Location{
				ID:           i,
				Name:         name,
				Country:      models.LL2Country{Name: country},
				TimezoneName: fmt.Sprintf("Timezone %02d", i),
				Description:  fmt.Sprintf("Description %02d", i),
				CelestialBody: models.LL2CelestialBodyDetailed{
					LL2CelestialBodyMini: models.LL2CelestialBodyMini{Name: "Earth"},
				},
				Latitude:  float64(i),
				Longitude: float64(i) * -1,
			},
		})
		assert.NoError(t, err)
	}

	searchPage, err := service.GetPublicLaunchBases(0, "cape")
	assert.NoError(t, err)
	assert.Equal(t, 2, searchPage.Count)
	assert.Len(t, searchPage.LaunchBases, 2)
	assert.Equal(t, "Cape Canaveral", searchPage.LaunchBases[0].Name)
	assert.Equal(t, "Cape Test Range", searchPage.LaunchBases[1].Name)

	pagedResult, err := service.GetPublicLaunchBases(1, "")
	assert.NoError(t, err)
	assert.Equal(t, 22, pagedResult.Count)
	assert.Len(t, pagedResult.LaunchBases, 2)
}

func TestMainService_PublicLaunchListSupportsSearchAndPaging(t *testing.T) {
	service, testDB, cleanup := newTestMainService(t)
	defer cleanup()
	ctx := context.Background()

	clearCoreCollections(t, testDB, COLLECTION_LAUNCH, COLLECTION_LL2_LAUNCH, COLLECTION_ROCKET, COLLECTION_AGENCY, COLLECTION_LAUNCH_BASE)

	_, err := testDB.Collection(COLLECTION_ROCKET).InsertOne(ctx, models.Rocket{ID: 4001, ExternalID: 501, ThumbImage: "rocket-thumb.jpg"})
	assert.NoError(t, err)
	_, err = testDB.Collection(COLLECTION_AGENCY).InsertOne(ctx, models.Agency{ID: 5001, ExternalID: 601, ThumbImage: "company-thumb.jpg"})
	assert.NoError(t, err)
	_, err = testDB.Collection(COLLECTION_LAUNCH_BASE).InsertOne(ctx, models.LaunchBase{ID: 6001, ExternalID: 701})
	assert.NoError(t, err)

	for i := 1; i <= 22; i++ {
		name := fmt.Sprintf("Mission %02d", i)
		if i == 1 {
			name = "Featured Demo"
		}
		if i == 2 {
			name = "Featured Return"
		}

		externalID := fmt.Sprintf("launch-%02d", i)
		launchTime := time.Now().UTC().Add(-time.Duration(i) * time.Hour).Format(time.RFC3339)

		_, err := testDB.Collection(COLLECTION_LAUNCH).InsertOne(ctx, models.Launch{
			ID:         int64(7000 + i),
			ExternalID: externalID,
			ThumbImage: fmt.Sprintf("launch-%02d.jpg", i),
		})
		assert.NoError(t, err)

		_, err = testDB.Collection(COLLECTION_LL2_LAUNCH).InsertOne(ctx, models.LL2LaunchDetailed{
			LL2LaunchNormal: models.LL2LaunchNormal{
				LL2LaunchBasic: models.LL2LaunchBasic{
					ID:   externalID,
					Name: name,
					Net:  launchTime,
					Status: models.LL2Status{
						ID:   3,
						Name: "Launch Successful",
					},
				},
				LaunchServiceProvider: models.LL2AgencyMini{ID: 601, Name: "SpaceX"},
				Rocket: models.LL2RocketNormal{
					Configuration: models.LL2LauncherConfigList{ID: 501, Name: "Falcon 9", FullName: "Falcon 9 Block 5"},
				},
				Pad: models.LL2Pad{
					Location: models.LL2Location{
						ID:           701,
						Name:         "Cape Canaveral",
						TimezoneName: "America/New_York",
						Country:      models.LL2Country{Name: "United States"},
						Latitude:     28.5,
						Longitude:    -80.5,
					},
				},
			},
		})
		assert.NoError(t, err)
	}

	searchPage, err := service.GetPublicLaunches(0, "featured")
	assert.NoError(t, err)
	assert.Equal(t, 2, searchPage.Count)
	assert.Len(t, searchPage.Launches, 2)
	assert.Equal(t, "Featured Demo", searchPage.Launches[0].Name)
	assert.Equal(t, "Featured Return", searchPage.Launches[1].Name)

	pagedResult, err := service.GetPublicLaunches(1, "")
	assert.NoError(t, err)
	assert.Equal(t, 22, pagedResult.Count)
	assert.Len(t, pagedResult.Launches, 2)
}

func TestMainService_PublicLaunchListOrdersFutureThenPast(t *testing.T) {
	service, testDB, cleanup := newTestMainService(t)
	defer cleanup()
	ctx := context.Background()

	clearCoreCollections(t, testDB, COLLECTION_LAUNCH, COLLECTION_LL2_LAUNCH, COLLECTION_ROCKET, COLLECTION_AGENCY, COLLECTION_LAUNCH_BASE)

	_, err := testDB.Collection(COLLECTION_ROCKET).InsertOne(ctx, models.Rocket{ID: 8101, ExternalID: 901, ThumbImage: "rocket-thumb.jpg"})
	assert.NoError(t, err)
	_, err = testDB.Collection(COLLECTION_AGENCY).InsertOne(ctx, models.Agency{ID: 8201, ExternalID: 902, ThumbImage: "company-thumb.jpg"})
	assert.NoError(t, err)
	_, err = testDB.Collection(COLLECTION_LAUNCH_BASE).InsertOne(ctx, models.LaunchBase{ID: 8301, ExternalID: 903})
	assert.NoError(t, err)

	testLaunches := []struct {
		id         string
		name       string
		launchTime time.Time
	}{
		{id: "future-far", name: "Future Far", launchTime: time.Now().UTC().Add(10 * 24 * time.Hour)},
		{id: "future-near", name: "Future Near", launchTime: time.Now().UTC().Add(2 * 24 * time.Hour)},
		{id: "past-near", name: "Past Near", launchTime: time.Now().UTC().Add(-2 * time.Hour)},
		{id: "past-far", name: "Past Far", launchTime: time.Now().UTC().Add(-48 * time.Hour)},
	}

	for index, launch := range testLaunches {
		_, err := testDB.Collection(COLLECTION_LAUNCH).InsertOne(ctx, models.Launch{
			ID:         int64(8400 + index),
			ExternalID: launch.id,
			ThumbImage: fmt.Sprintf("%s.jpg", launch.id),
		})
		assert.NoError(t, err)

		_, err = testDB.Collection(COLLECTION_LL2_LAUNCH).InsertOne(ctx, models.LL2LaunchDetailed{
			LL2LaunchNormal: models.LL2LaunchNormal{
				LL2LaunchBasic: models.LL2LaunchBasic{
					ID:   launch.id,
					Name: launch.name,
					Net:  launch.launchTime.Format(time.RFC3339),
					Status: models.LL2Status{
						ID:   1,
						Name: "Go for Launch",
					},
				},
				LaunchServiceProvider: models.LL2AgencyMini{ID: 902, Name: "Provider"},
				Rocket: models.LL2RocketNormal{
					Configuration: models.LL2LauncherConfigList{ID: 901, Name: "Vehicle", FullName: "Vehicle Full"},
				},
				Pad: models.LL2Pad{
					Location: models.LL2Location{
						ID:           903,
						Name:         "Test Site",
						TimezoneName: "UTC",
						Country:      models.LL2Country{Name: "Testland"},
					},
				},
			},
		})
		assert.NoError(t, err)
	}

	result, err := service.GetPublicLaunches(0, "")
	assert.NoError(t, err)
	assert.Equal(t, 4, result.Count)
	assert.Len(t, result.Launches, 4)
	assert.Equal(t, "Future Far", result.Launches[0].Name)
	assert.Equal(t, "Future Near", result.Launches[1].Name)
	assert.Equal(t, "Past Near", result.Launches[2].Name)
	assert.Equal(t, "Past Far", result.Launches[3].Name)
}
