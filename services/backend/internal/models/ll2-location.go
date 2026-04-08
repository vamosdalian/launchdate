package models

type LL2PadResponse struct {
	Count    int      `json:"count" bson:"count"`
	Next     string   `json:"next" bson:"next"`
	Previous string   `json:"previous" bson:"previous"`
	Results  []LL2Pad `json:"results" bson:"results"`
}

type LL2PadSerializerNoLocation struct {
	Id                        int               `json:"id" bson:"id"`
	URL                       string            `json:"url" bson:"url"`
	Active                    bool              `json:"active" bson:"active"`
	Agencies                  []LL2AgencyNormal `json:"agencies" bson:"agencies"`
	Name                      string            `json:"name" bson:"name"`
	Image                     LL2Image          `json:"image" bson:"image"`
	Description               string            `json:"description" bson:"description"`
	InfoURL                   string            `json:"info_url" bson:"info_url"`
	WikiURL                   string            `json:"wiki_url" bson:"wiki_url"`
	MapURL                    string            `json:"map_url" bson:"map_url"`
	Latitude                  float64           `json:"latitude" bson:"latitude"`
	Longitude                 float64           `json:"longitude" bson:"longitude"`
	Country                   LL2Country        `json:"country" bson:"country"`
	MapImage                  string            `json:"map_image" bson:"map_image"`
	TotalLaunchCount          int               `json:"total_launch_count" bson:"total_launch_count"`
	OrbitalLaunchAttemptCount int               `json:"orbital_launch_attempt_count" bson:"orbital_launch_attempt_count"`
	FastestTurnaround         string            `json:"fastest_turnaround" bson:"fastest_turnaround"`
}

type LL2Pad struct {
	LL2PadSerializerNoLocation `bson:",inline"`
	Location                   LL2Location `json:"location" bson:"location"`
}

type LL2LocationResponse struct {
	Count    int                             `json:"count" bson:"count"`
	Next     string                          `json:"next" bson:"next"`
	Previous string                          `json:"previous" bson:"previous"`
	Results  []LL2LocationSerializerWithPads `json:"results" bson:"results"`
}

type LL2Location struct {
	ResponseMode      string                   `json:"response_mode" bson:"response_mode"`
	ID                int                      `json:"id" bson:"id"`
	Name              string                   `json:"name" bson:"name"`
	URL               string                   `json:"url" bson:"url"`
	CelestialBody     LL2CelestialBodyDetailed `json:"celestial_body" bson:"celestial_body"`
	Active            bool                     `json:"active" bson:"active"`
	Country           LL2Country               `json:"country" bson:"country"`
	Description       string                   `json:"description" bson:"description"`
	Image             LL2Image                 `json:"image" bson:"image"`
	MapImage          string                   `json:"map_image" bson:"map_image"`
	Latitude          float64                  `json:"latitude" bson:"latitude"`
	Longitude         float64                  `json:"longitude" bson:"longitude"`
	TimezoneName      string                   `json:"timezone_name" bson:"timezone_name"`
	TotalLaunchCount  int                      `json:"total_launch_count" bson:"total_launch_count"`
	TotalLandingCount int                      `json:"total_landing_count" bson:"total_landing_count"`
}

type LL2LocationSerializerWithPads struct {
	LL2Location `bson:",inline"`
	Pads        []LL2PadSerializerNoLocation `json:"pads" bson:"pads"`
}
