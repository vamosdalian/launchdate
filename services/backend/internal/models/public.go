package models

type PublicLaunchStatus string

const (
	PublicLaunchStatusScheduled PublicLaunchStatus = "scheduled"
	PublicLaunchStatusSuccess   PublicLaunchStatus = "success"
	PublicLaunchStatusFailure   PublicLaunchStatus = "failure"
	PublicLaunchStatusCancelled PublicLaunchStatus = "cancelled"
	PublicLaunchStatusDelayed   PublicLaunchStatus = "delayed"
	PublicLaunchStatusInFlight  PublicLaunchStatus = "in_flight"
	PublicLaunchStatusUnknown   PublicLaunchStatus = "unknown"
)

type PublicRocketRef struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	ImageURL   string `json:"image_url"`
	ThumbImage string `json:"thumb_image"`
}

type PublicCompanyRef struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

type PublicLaunchBaseRef struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Location  string  `json:"location"`
	Country   string  `json:"country"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type PublicLaunchSummary struct {
	ID              string              `json:"id"`
	Name            string              `json:"name"`
	LaunchTime      string              `json:"launch_time"`
	Status          PublicLaunchStatus  `json:"status"`
	StatusLabel     string              `json:"status_label"`
	ThumbImage      string              `json:"thumb_image"`
	BackgroundImage string              `json:"background_image"`
	Rocket          PublicRocketRef     `json:"rocket"`
	Company         PublicCompanyRef    `json:"company"`
	LaunchBase      PublicLaunchBaseRef `json:"launch_base"`
}

type PublicLaunchPage struct {
	Count    int                   `json:"count"`
	Launches []PublicLaunchSummary `json:"launches"`
}

type PublicMissionSummary struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type PublicTimelineEntry struct {
	RelativeTime string `json:"relative_time"`
	Abbrev       string `json:"abbrev"`
	Description  string `json:"description"`
}

type PublicLaunchView struct {
	ID              string                 `json:"id"`
	Name            string                 `json:"name"`
	LaunchTime      string                 `json:"launch_time"`
	Status          PublicLaunchStatus     `json:"status"`
	StatusLabel     string                 `json:"status_label"`
	BackgroundImage string                 `json:"background_image"`
	ImageList       []string               `json:"image_list"`
	Rocket          PublicRocketRef        `json:"rocket"`
	Company         PublicCompanyRef       `json:"company"`
	LaunchBase      PublicLaunchBaseRef    `json:"launch_base"`
	Missions        []PublicMissionSummary `json:"missions"`
	Timeline        []PublicTimelineEntry  `json:"timeline"`
}

type PublicRocketListItem struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	ThumbImage string `json:"thumb_image"`
}

type PublicRocketPage struct {
	Count   int                    `json:"count"`
	Rockets []PublicRocketListItem `json:"rockets"`
}

type PublicRocketView struct {
	ID              string                `json:"id"`
	Name            string                `json:"name"`
	Description     string                `json:"description"`
	Active          bool                  `json:"active"`
	Reusable        bool                  `json:"reusable"`
	LaunchImage     string                `json:"launch_image"`
	MainImage       string                `json:"main_image"`
	ImageList       []string              `json:"image_list"`
	Company         PublicCompanyRef      `json:"company"`
	Launches        []PublicLaunchSummary `json:"launches"`
	LaunchCost      float64               `json:"launch_cost"`
	Diameter        float64               `json:"diameter"`
	Length          float64               `json:"length"`
	LiftoffThrust   float64               `json:"liftoff_thrust"`
	LaunchMass      float64               `json:"launch_mass"`
	LeoCapacity     float64               `json:"leo_capacity"`
	GtoCapacity     float64               `json:"gto_capacity"`
	GeoCapacity     float64               `json:"geo_capacity"`
	SsoCapacity     float64               `json:"sso_capacity"`
	TotalLaunches   int                   `json:"total_launches"`
	SuccessLaunches int                   `json:"success_launches"`
	FailureLaunches int                   `json:"failure_launches"`
	TotalLandings   int                   `json:"total_landings"`
	SuccessLandings int                   `json:"success_landings"`
	FailureLandings int                   `json:"failure_landings"`
}

type PublicCompanyStats struct {
	RocketCount        int `json:"rocket_count"`
	LaunchCount        int `json:"launch_count"`
	SuccessfulLaunches int `json:"successful_launches"`
	FailedLaunches     int `json:"failed_launches"`
	PendingLaunches    int `json:"pending_launches"`
}

type PublicCompanyListItem struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Founded      int    `json:"founded"`
	Founder      string `json:"founder"`
	Headquarters string `json:"headquarters"`
	Employees    int    `json:"employees"`
	Website      string `json:"website"`
	ImageURL     string `json:"image_url"`
}

type PublicCompanyPage struct {
	Count     int                     `json:"count"`
	Companies []PublicCompanyListItem `json:"companies"`
}

type PublicCompanyView struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	Founded      int                    `json:"founded"`
	Founder      string                 `json:"founder"`
	Headquarters string                 `json:"headquarters"`
	Employees    int                    `json:"employees"`
	Website      string                 `json:"website"`
	ImageURL     string                 `json:"image_url"`
	Rockets      []PublicRocketListItem `json:"rockets"`
	Launches     []PublicLaunchSummary  `json:"launches"`
	Stats        PublicCompanyStats     `json:"stats"`
}

type PublicLaunchBaseStats struct {
	LaunchCount         int `json:"launch_count"`
	UpcomingLaunchCount int `json:"upcoming_launch_count"`
	SuccessfulLaunches  int `json:"successful_launches"`
	FailedLaunches      int `json:"failed_launches"`
	SuccessRate         int `json:"success_rate"`
}

type PublicLaunchBaseListItem struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Location    string  `json:"location"`
	Country     string  `json:"country"`
	Description string  `json:"description"`
	ImageURL    string  `json:"image_url"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}

type PublicLaunchBasePage struct {
	Count       int                        `json:"count"`
	LaunchBases []PublicLaunchBaseListItem `json:"launch_bases"`
}

type PublicLaunchBaseView struct {
	ID          string                `json:"id"`
	Name        string                `json:"name"`
	Location    string                `json:"location"`
	Country     string                `json:"country"`
	Description string                `json:"description"`
	ImageURL    string                `json:"image_url"`
	Latitude    float64               `json:"latitude"`
	Longitude   float64               `json:"longitude"`
	Launches    []PublicLaunchSummary `json:"launches"`
	Stats       PublicLaunchBaseStats `json:"stats"`
}
