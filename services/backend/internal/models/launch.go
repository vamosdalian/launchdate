package models

// for DB
type Launch struct {
	ID              int64    `bson:"id" json:"id,string"`
	ExternalID      string   `bson:"external_id" json:"external_id"` // ll2 launch id
	BackgroundImage string   `bson:"background_image" json:"background_image"`
	ImageList       []string `bson:"image_list" json:"image_list"`
	ThumbImage      string   `bson:"thumb_image" json:"thumb_image"`
}

type LaunchSerializer struct {
	Launch
	Data LL2LaunchNormal `json:"data"`
}

type LaunchList struct {
	Count    int                `json:"count"`
	Launches []LaunchSerializer `json:"launches"`
}

type PublicCompactLaunch struct {
	ID         int64  `json:"id,string"`
	Name       string `json:"name"`
	LaunchTime string `json:"launch_time"` // ISO 8601 format
	Status     int    `json:"status"`
	ThumbImage string `json:"thumb_image"`
	RocketName string `json:"rocket_name"`
	AgencyName string `json:"agency_name"`
	Location   string `json:"location"`
}

type PublicLaunchList struct {
	Count    int                   `json:"count"`
	Launches []PublicCompactLaunch `json:"launches"`
}

type PublicLaunchDetail struct {
	ID              int64                 `json:"id,string"`
	Name            string                `json:"name"`
	LaunchTime      string                `json:"launch_time"` // ISO 8601 format
	Status          int                   `json:"status"`
	BackgroundImage string                `json:"background_image"`
	ImageList       []string              `json:"image_list"`
	RocketInfo      PublicCompactRocket   `json:"rocket_info"`
	AgencyInfo      PublicCompactAgency   `json:"agency_info"`
	LocationInfo    PublicCompactLocation `json:"location_info"`
	MissionInfo     []Mission             `json:"mission_info"`
	TimelineEvent   []TimelineEvent       `json:"timeline_event"`
}

type Mission struct {
	ID          int64  `json:"id,string"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type TimelineEvent struct {
	RelativeTime string `json:"relative_time"`
	Abbrev       string `json:"abbrev"`
	Description  string `json:"description"`
}
