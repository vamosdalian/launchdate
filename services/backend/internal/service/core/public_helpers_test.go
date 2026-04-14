package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vamosdalian/launchdate-backend/internal/models"
)

func TestNormalizePublicLaunchStatus(t *testing.T) {
	tests := []struct {
		name           string
		status         models.LL2Status
		expectedStatus models.PublicLaunchStatus
		expectedLabel  string
	}{
		{
			name: "scheduled from go abbreviation",
			status: models.LL2Status{
				Name:   "To Be Determined",
				Abbrev: "Go",
			},
			expectedStatus: models.PublicLaunchStatusScheduled,
			expectedLabel:  "Go",
		},
		{
			name: "success from status id",
			status: models.LL2Status{
				ID:   3,
				Name: "Launch Successful",
			},
			expectedStatus: models.PublicLaunchStatusSuccess,
			expectedLabel:  "Launch Successful",
		},
		{
			name: "cancelled from status text",
			status: models.LL2Status{
				Name:   "Launch Cancelled",
				Abbrev: "Scrubbed",
			},
			expectedStatus: models.PublicLaunchStatusCancelled,
			expectedLabel:  "Scrubbed",
		},
		{
			name: "unknown fallback label",
			status: models.LL2Status{
				ID: 0,
			},
			expectedStatus: models.PublicLaunchStatusUnknown,
			expectedLabel:  "Unknown",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			status, label := normalizePublicLaunchStatus(test.status)
			assert.Equal(t, test.expectedStatus, status)
			assert.Equal(t, test.expectedLabel, label)
		})
	}
}

func TestBuildPublicLaunchSummaryUsesCoreIDs(t *testing.T) {
	service := &MainService{}
	internalLaunch := models.Launch{
		ID:              91001,
		ExternalID:      "ll2-launch-1",
		ThumbImage:      "launch-thumb.jpg",
		BackgroundImage: "launch-background.jpg",
	}
	ll2Launch := models.LL2LaunchDetailed{
		LL2LaunchNormal: models.LL2LaunchNormal{
			LL2LaunchBasic: models.LL2LaunchBasic{
				Name: "Falcon 9 | Demo Payload",
				Net:  "2026-04-14T06:13:00Z",
				Status: models.LL2Status{
					ID:   1,
					Name: "Go for Launch",
				},
			},
			Rocket: models.LL2RocketNormal{
				Configuration: models.LL2LauncherConfigList{
					ID:   501,
					Name: "Falcon 9",
				},
			},
			LaunchServiceProvider: models.LL2AgencyMini{
				ID:   601,
				Name: "SpaceX",
			},
			Pad: models.LL2Pad{
				Location: models.LL2Location{
					ID:           701,
					Name:         "Cape Canaveral SFS, FL, USA",
					TimezoneName: "America/New_York",
					Latitude:     28.488889,
					Longitude:    -80.577778,
					Country: models.LL2Country{
						Name: "United States of America",
					},
				},
			},
		},
	}

	rocketMap := map[int64]models.Rocket{
		501: {
			ID:         11001,
			ExternalID: 501,
			MainImage:  "rocket-main.jpg",
			ThumbImage: "rocket-thumb.jpg",
		},
	}
	agencyMap := map[int64]models.Agency{
		601: {
			ID:         12001,
			ExternalID: 601,
			ThumbImage: "company-logo.png",
		},
	}
	launcherMap := map[int64]models.LL2LauncherConfigNormal{
		501: {
			LL2LauncherConfigList: models.LL2LauncherConfigList{
				ID:   501,
				Name: "Falcon 9",
			},
			Image: models.LL2Image{ImageURL: "ll2-rocket.jpg", ThumbnailURL: "ll2-rocket-thumb.jpg"},
		},
	}
	baseMap := map[int64]models.LaunchBase{
		701: {
			ID:         13001,
			ExternalID: 701,
		},
	}

	summary, include := service.buildPublicLaunchSummary(internalLaunch, ll2Launch, rocketMap, launcherMap, agencyMap, nil, baseMap)
	assert.True(t, include)
	assert.Equal(t, "91001", summary.ID)
	assert.Equal(t, "Falcon 9 | Demo Payload", summary.Name)
	assert.Equal(t, "2026-04-14T06:13:00Z", summary.LaunchTime)
	assert.Equal(t, models.PublicLaunchStatusScheduled, summary.Status)
	assert.Equal(t, "Go for Launch", summary.StatusLabel)
	assert.Equal(t, "launch-thumb.jpg", summary.ThumbImage)
	assert.Equal(t, "launch-background.jpg", summary.BackgroundImage)

	assert.Equal(t, "11001", summary.Rocket.ID)
	assert.Equal(t, "Falcon 9", summary.Rocket.Name)
	assert.Equal(t, "rocket-main.jpg", summary.Rocket.ImageURL)
	assert.Equal(t, "rocket-thumb.jpg", summary.Rocket.ThumbImage)

	assert.Equal(t, "12001", summary.Company.ID)
	assert.Equal(t, "SpaceX", summary.Company.Name)
	assert.Equal(t, "company-logo.png", summary.Company.ImageURL)

	assert.Equal(t, "13001", summary.LaunchBase.ID)
	assert.Equal(t, "Cape Canaveral SFS, FL, USA", summary.LaunchBase.Name)
	assert.Equal(t, "America/New_York", summary.LaunchBase.Location)
	assert.Equal(t, "United States of America", summary.LaunchBase.Country)
	assert.Equal(t, 28.488889, summary.LaunchBase.Latitude)
	assert.Equal(t, -80.577778, summary.LaunchBase.Longitude)
}

func TestBuildPublicLaunchSummarySkipsMissingCoreLaunch(t *testing.T) {
	service := &MainService{}
	summary, include := service.buildPublicLaunchSummary(models.Launch{}, models.LL2LaunchDetailed{}, nil, nil, nil, nil, nil)
	assert.False(t, include)
	assert.Equal(t, models.PublicLaunchSummary{}, summary)
}

func TestBuildPublicLaunchSummaryUsesLL2ImageFallbacks(t *testing.T) {
	service := &MainService{}
	internalLaunch := models.Launch{ID: 91002, ExternalID: "ll2-launch-2"}
	ll2Launch := models.LL2LaunchDetailed{
		LL2LaunchNormal: models.LL2LaunchNormal{
			LL2LaunchBasic: models.LL2LaunchBasic{
				Name:   "Fallback Mission",
				Net:    "2026-04-14T07:00:00Z",
				Image:  models.LL2Image{ImageURL: "launch-image.jpg", ThumbnailURL: "launch-thumb.jpg"},
				Status: models.LL2Status{ID: 1, Name: "Go for Launch"},
			},
			Rocket:                models.LL2RocketNormal{Configuration: models.LL2LauncherConfigList{ID: 777, Name: "Vehicle"}},
			LaunchServiceProvider: models.LL2AgencyMini{ID: 888, Name: "Agency"},
			Pad:                   models.LL2Pad{Location: models.LL2Location{ID: 999, Name: "Site"}},
		},
	}
	launcherMap := map[int64]models.LL2LauncherConfigNormal{
		777: {
			LL2LauncherConfigList: models.LL2LauncherConfigList{ID: 777, Name: "Vehicle"},
			Image:                 models.LL2Image{ImageURL: "rocket-image.jpg", ThumbnailURL: "rocket-thumb.jpg"},
		},
	}
	agencyDocMap := map[int64]models.LL2AgencyDetailed{
		888: {
			LL2AgencyNormal: models.LL2AgencyNormal{
				LL2AgencyMini: models.LL2AgencyMini{ID: 888, Name: "Agency"},
				Logo:          models.LL2Image{ImageURL: "agency-logo.jpg", ThumbnailURL: "agency-logo-thumb.jpg"},
			},
		},
	}
	baseMap := map[int64]models.LaunchBase{999: {ID: 13002, ExternalID: 999}}

	summary, include := service.buildPublicLaunchSummary(
		internalLaunch,
		ll2Launch,
		map[int64]models.Rocket{777: {ID: 11002, ExternalID: 777}},
		launcherMap,
		map[int64]models.Agency{888: {ID: 12002, ExternalID: 888}},
		agencyDocMap,
		baseMap,
	)
	assert.True(t, include)
	assert.Equal(t, "launch-thumb.jpg", summary.ThumbImage)
	assert.Equal(t, "launch-image.jpg", summary.BackgroundImage)
	assert.Equal(t, "rocket-image.jpg", summary.Rocket.ImageURL)
	assert.Equal(t, "rocket-thumb.jpg", summary.Rocket.ThumbImage)
	assert.Equal(t, "agency-logo.jpg", summary.Company.ImageURL)
}

func TestBuildPublicCompanyListItemPrefersExplicitWebsiteURL(t *testing.T) {
	item := buildPublicCompanyListItem(
		models.Agency{
			ID: 4001,
			SocialUrl: []models.SocialUrl{
				{Name: "Twitter", URL: "https://twitter.com/example"},
				{Name: "Website", URL: "https://www.example.com"},
			},
		},
		models.LL2AgencyDetailed{
			LL2AgencyNormal: models.LL2AgencyNormal{
				LL2AgencyMini: models.LL2AgencyMini{Name: "Example Agency"},
			},
			InfoURL: "https://ll2.example.com/info",
			WikiURL: "https://ll2.example.com/wiki",
			SocialMediaLinks: []models.LL2SocialMediaLink{{
				SocialMedia: models.LL2SocialMedia{Name: "Homepage"},
				URL:         "https://homepage.example.com",
			}},
		},
	)

	assert.Equal(t, "https://www.example.com", item.Website)
}

func TestBuildPublicCompanyListItemFallsBackToInfoURLBeforeWiki(t *testing.T) {
	item := buildPublicCompanyListItem(
		models.Agency{
			ID:        4002,
			SocialUrl: []models.SocialUrl{{Name: "Twitter", URL: "https://twitter.com/example"}},
		},
		models.LL2AgencyDetailed{
			LL2AgencyNormal: models.LL2AgencyNormal{
				LL2AgencyMini: models.LL2AgencyMini{Name: "Fallback Agency"},
			},
			InfoURL: "https://ll2.example.com/info",
			WikiURL: "https://ll2.example.com/wiki",
		},
	)

	assert.Equal(t, "https://ll2.example.com/info", item.Website)
}
