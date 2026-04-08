// LL2Launch represents the launch data structure from Launch Library 2 API
// https://ll.thespacedevs.com/2.3.0/json
package models

type LL2Response struct {
	Count    int                  `json:"count"`
	Next     string               `json:"next"`
	Previous string               `json:"previous"`
	Results  []*LL2LaunchDetailed `json:"results"`
}

type LL2LaunchBasic struct {
	ID               string          `json:"id" bson:"id"`
	URL              string          `json:"url" bson:"url"`
	Name             string          `json:"name" bson:"name"`
	ResponseMode     string          `json:"response_mode" bson:"response_mode"`
	Slug             string          `json:"slug" bson:"slug"`
	LaunchDesignator string          `json:"launch_designator" bson:"launch_designator"`
	Status           LL2Status       `json:"status" bson:"status"`
	LastUpdated      string          `json:"last_updated" bson:"last_updated"`
	Net              string          `json:"net" bson:"net"`
	NetPrecision     LL2NetPrecision `json:"net_precision" bson:"net_precision"`
	WindowEnd        string          `json:"window_end" bson:"window_end"`
	WindowStart      string          `json:"window_start" bson:"window_start"`
	Image            LL2Image        `json:"image" bson:"image"`
	Infographic      string          `json:"infographic" bson:"infographic"`
}

type LL2Image struct {
	ID           int               `json:"id" bson:"id"`
	Name         string            `json:"name" bson:"name"`
	ImageURL     string            `json:"image_url" bson:"image_url"`
	ThumbnailURL string            `json:"thumbnail_url" bson:"thumbnail_url"`
	Credit       string            `json:"credit" bson:"credit"`
	License      LL2ImageLicense   `json:"license" bson:"license"`
	SingleUse    bool              `json:"single_use" bson:"single_use"`
	Variants     []LL2ImageVariant `json:"variants" bson:"variants"`
}

type LL2ImageLicense struct {
	ID       int    `json:"id" bson:"id"`
	Name     string `json:"name" bson:"name"`
	Priority int    `json:"priority" bson:"priority"`
	Link     string `json:"link" bson:"link"`
}

type LL2ImageVariant struct {
	ID     int    `json:"id" bson:"id"`
	URL    string `json:"url" bson:"url"`
	Width  int    `json:"width" bson:"width"`
	Height int    `json:"height" bson:"height"`
}

type LL2Status struct {
	ID          int    `json:"id" bson:"id"`
	Name        string `json:"name" bson:"name"`
	Abbrev      string `json:"abbrev" bson:"abbrev"`
	Description string `json:"description" bson:"description"`
}

type LL2NetPrecision struct {
	ID          int    `json:"id" bson:"id"`
	Name        string `json:"name" bson:"name"`
	Abbrev      string `json:"abbrev" bson:"abbrev"`
	Description string `json:"description" bson:"description"`
}

type LL2LaunchNormal struct {
	LL2LaunchBasic                 `bson:",inline"`
	Probability                    int                `json:"probability" bson:"probability"`
	WeatherConcerns                string             `json:"weather_concerns" bson:"weather_concerns"`
	FailReason                     string             `json:"failreason" bson:"failreason"`
	HashTag                        string             `json:"hashtag" bson:"hashtag"`
	LaunchServiceProvider          LL2AgencyMini      `json:"launch_service_provider" bson:"launch_service_provider"`
	Rocket                         LL2RocketNormal    `json:"rocket" bson:"rocket"`
	Mission                        LL2Mission         `json:"mission" bson:"mission"`
	Pad                            LL2Pad             `json:"pad" bson:"pad"`
	WebcastLive                    bool               `json:"webcast_live" bson:"webcast_live"`
	Program                        []LL2ProgramNormal `json:"program" bson:"program"`
	OrbitalLaunchAttemptCount      int                `json:"orbital_launch_attempt_count" bson:"orbital_launch_attempt_count"`
	LocationLaunchAttemptCount     int                `json:"location_launch_attempt_count" bson:"location_launch_attempt_count"`
	PadLaunchAttemptCount          int                `json:"pad_launch_attempt_count" bson:"pad_launch_attempt_count"`
	AgencyLaunchAttemptCount       int                `json:"agency_launch_attempt_count" bson:"agency_launch_attempt_count"`
	OrbitalLaunchAttemptCountYear  int                `json:"orbital_launch_attempt_count_year" bson:"orbital_launch_attempt_count_year"`
	LocationLaunchAttemptCountYear int                `json:"location_launch_attempt_count_year" bson:"location_launch_attempt_count_year"`
	PadLaunchAttemptCountYear      int                `json:"pad_launch_attempt_count_year" bson:"pad_launch_attempt_count_year"`
	AgencyLaunchAttemptCountYear   int                `json:"agency_launch_attempt_count_year" bson:"agency_launch_attempt_count_year"`
}

type LL2RocketNormal struct {
	ID            int                   `json:"id" bson:"id"`
	Configuration LL2LauncherConfigList `json:"configuration" bson:"configuration"`
}

type LL2Mission struct {
	ID          int                 `json:"id" bson:"id"`
	Name        string              `json:"name" bson:"name"`
	Description string              `json:"description" bson:"description"`
	Type        string              `json:"type" bson:"type"`
	Image       LL2Image            `json:"image" bson:"image"`
	Orbit       LL2Orbit            `json:"orbit" bson:"orbit"`
	Agencies    []LL2AgencyDetailed `json:"agencies" bson:"agencies"`
	InfoURLs    []LL2InfoURL        `json:"info_urls" bson:"info_urls"`
	VidURLs     []LL2VidURL         `json:"vid_urls" bson:"vid_urls"`
}

type LL2Orbit struct {
	ID            int                  `json:"id" bson:"id"`
	Name          string               `json:"name" bson:"name"`
	Abbrev        string               `json:"abbrev" bson:"abbrev"`
	CelestialBody LL2CelestialBodyMini `json:"celestial_body" bson:"celestial_body"`
}

type LL2CelestialBodyMini struct {
	ResponseMode string `json:"response_mode" bson:"response_mode"`
	ID           int    `json:"id" bson:"id"`
	Name         string `json:"name" bson:"name"`
}

type LL2Country struct {
	ID                      int    `json:"id" bson:"id"`
	Name                    string `json:"name" bson:"name"`
	Alpha2Code              string `json:"alpha2_code" bson:"alpha2_code"`
	Alpha3Code              string `json:"alpha3_code" bson:"alpha3_code"`
	NationalityName         string `json:"nationality_name" bson:"nationality_name"`
	NationalityNameComposed string `json:"nationality_name_composed" bson:"nationality_name_composed"`
}

type LL2SocialMediaLink struct {
	ID          int            `json:"id" bson:"id"`
	SocialMedia LL2SocialMedia `json:"social_media" bson:"social_media"`
	URL         string         `json:"url" bson:"url"`
}

type LL2SocialMedia struct {
	Id   int      `json:"id" bson:"id"`
	Name string   `json:"name" bson:"name"`
	URL  string   `json:"url" bson:"url"`
	Logo LL2Image `json:"logo" bson:"logo"`
}

type LL2InfoURL struct {
	Priority     int            `json:"priority" bson:"priority"`
	Source       string         `json:"source" bson:"source"`
	Title        string         `json:"title" bson:"title"`
	Description  string         `json:"description" bson:"description"`
	FeatureImage string         `json:"feature_image" bson:"feature_image"`
	URL          string         `json:"url" bson:"url"`
	Type         LL2InfoURLType `json:"type" bson:"type"`
	Language     LL2Language    `json:"language" bson:"language"`
}

type LL2InfoURLType struct {
	Id   int    `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
}

type LL2Language struct {
	ID   int    `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
	Code string `json:"code" bson:"code"`
}

type LL2VidURL struct {
	Priority     int           `json:"priority" bson:"priority"`
	Source       string        `json:"source" bson:"source"`
	Publisher    string        `json:"publisher" bson:"publisher"`
	Title        string        `json:"title" bson:"title"`
	Description  string        `json:"description" bson:"description"`
	FeatureImage string        `json:"feature_image" bson:"feature_image"`
	URL          string        `json:"url" bson:"url"`
	Type         LL2VidURLType `json:"type" bson:"type"`
	Language     LL2Language   `json:"language" bson:"language"`
	StartTime    string        `json:"start_time" bson:"start_time"`
	EndTime      string        `json:"end_time" bson:"end_time"`
	Live         bool          `json:"live" bson:"live"`
}

type LL2VidURLType struct {
	Id   int    `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
}

type LL2CelestialBodyDetailed struct {
	LL2CelestialBodyMini   `bson:",inline"`
	Type                   LL2CelestialBodyType `json:"type" bson:"type"`
	Diameter               float64              `json:"diameter" bson:"diameter"`
	Mass                   float64              `json:"mass" bson:"mass"`
	Gravity                float64              `json:"gravity" bson:"gravity"`
	LengthOfDay            string               `json:"length_of_day" bson:"length_of_day"`
	Atmosphere             bool                 `json:"atmosphere" bson:"atmosphere"`
	Image                  LL2Image             `json:"image" bson:"image"`
	Description            string               `json:"description" bson:"description"`
	WikiURL                string               `json:"wiki_url" bson:"wiki_url"`
	TotalAttemptedLandes   int                  `json:"total_attempted_landes" bson:"total_attempted_landes"`
	SuccessfulLaunches     int                  `json:"successful_launches" bson:"successful_launches"`
	FailedLaunches         int                  `json:"failed_launches" bson:"failed_launches"`
	TotalAttemptedLandings int                  `json:"total_attempted_landings" bson:"total_attempted_landings"`
	SuccessfulLandings     int                  `json:"successful_landings" bson:"successful_landings"`
	FailedLandings         int                  `json:"failed_landings" bson:"failed_landings"`
}

type LL2CelestialBodyType struct {
	ID   int    `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
}

type LL2ProgramNormal struct {
	ResponseMode   string            `json:"response_mode" bson:"response_mode"`
	ID             int               `json:"id" bson:"id"`
	URL            string            `json:"url" bson:"url"`
	Name           string            `json:"name" bson:"name"`
	Image          LL2Image          `json:"image" bson:"image"`
	InfoUrl        string            `json:"info_url" bson:"info_url"`
	WikiUrl        string            `json:"wiki_url" bson:"wiki_url"`
	Description    string            `json:"description" bson:"description"`
	Agencies       []LL2AgencyMini   `json:"agencies" bson:"agencies"`
	StartDate      string            `json:"start_date" bson:"start_date"`
	EndDate        string            `json:"end_date" bson:"end_date"`
	MissionPatches []LL2MissionPatch `json:"mission_patches" bson:"mission_patches"`
	Type           LL2ProgramType    `json:"type" bson:"type"`
}
type LL2MissionPatch struct {
	ID           int           `json:"id" bson:"id"`
	Name         string        `json:"name" bson:"name"`
	Priority     int           `json:"priority" bson:"priority"`
	ImageUrl     string        `json:"image_url" bson:"image_url"`
	Agency       LL2AgencyMini `json:"agency" bson:"agency"`
	ResponseMode string        `json:"response_mode" bson:"response_mode"`
}

type LL2ProgramType struct {
	ID   int    `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
}

type LL2LaunchDetailed struct {
	LL2LaunchNormal `bson:",inline"`
	FlightclubUrl   string             `json:"flightclub_url" bson:"flightclub_url"`
	Updates         []LL2Update        `json:"updates" bson:"updates"`
	InfoURLs        []LL2InfoURL       `json:"info_urls" bson:"info_urls"`
	VidURLs         []LL2VidURL        `json:"vid_urls" bson:"vid_urls"`
	Timeline        []LL2TimelineEvent `json:"timeline" bson:"timeline"`
	PadTurnaround   string             `json:"pad_turnaround" bson:"pad_turnaround"`
	MissionPatches  []LL2MissionPatch  `json:"mission_patches" bson:"mission_patches"`
}

type LL2Update struct {
	ID           int    `json:"id" bson:"id"`
	ProfileImage string `json:"profile_image" bson:"profile_image"`
	Comment      string `json:"comment" bson:"comment"`
	InfoUrl      string `json:"info_url" bson:"info_url"`
	CreatedBy    string `json:"created_by" bson:"created_by"`
	CreatedOn    string `json:"created_on" bson:"created_on"`
}

type LL2TimelineEvent struct {
	RelativeTime string               `json:"relative_time" bson:"relative_time"`
	Type         LL2TimelineEventType `json:"type" bson:"type"`
}

type LL2TimelineEventType struct {
	ID          int    `json:"id" bson:"id"`
	Abbrev      string `json:"abbrev" bson:"abbrev"`
	Description string `json:"description" bson:"description"`
}
