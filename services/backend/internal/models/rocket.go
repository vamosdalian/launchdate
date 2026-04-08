package models

// for DB
type Rocket struct {
	ID          int64    `bson:"id" json:"id,string"`
	ExternalID  int64    `bson:"external_id" json:"external_id"` // ll2 rocket id
	LaunchImage string   `bson:"launch_image" json:"launch_image"`
	MainImage   string   `bson:"main_image" json:"main_image"`
	ThumbImage  string   `bson:"thumb_image" json:"thumb_image"`
	ImageList   []string `bson:"image_list" json:"image_list"`
}

type RocketSerializer struct {
	Rocket
	Data LL2LauncherConfigNormal `json:"data"`
}

type RocketList struct {
	Count   int                `json:"count"`
	Rockets []RocketSerializer `json:"rockets"`
}

type PublicCompactRocket struct {
	ID         int64  `json:"id,string"`
	Name       string `json:"name"`
	ThumbImage string `bson:"thumb_image" json:"thumb_image"`
}

type PublicRocketList struct {
	Count   int                   `json:"count"`
	Rockets []PublicCompactRocket `json:"rockets"`
}

type PublicRocketDetail struct {
	ID              int64                 `json:"id,string"`
	Name            string                `json:"name"`
	Description     string                `json:"description"`
	Active          bool                  `json:"active"`
	Reusable        bool                  `json:"reusable"`
	LaunchImage     string                `json:"launch_image"`
	MainImage       string                `json:"main_image"`
	ImageList       []string              `json:"image_list"`
	AgencyInfo      PublicCompactAgency   `json:"agency_info"`
	Launches        []PublicCompactLaunch `json:"launches"`
	LaunchCost      float64               `json:"launch_cost"`    // in USD
	Diameter        float64               `json:"diameter"`       // in meters
	Length          float64               `json:"length"`         // in meters
	LiftoffThrust   float64               `json:"liftoff_thrust"` // in kN
	LaunchMass      float64               `json:"launch_mass"`    // in kg
	LeoCapacity     float64               `json:"leo_capacity"`   // in kg
	GtoCapacity     float64               `json:"gto_capacity"`   // in kg
	GeoCapacity     float64               `json:"geo_capacity"`   // in kg
	SsoCapacity     float64               `json:"sso_capacity"`   // in kg
	TotalLaunches   int                   `json:"total_launches"`
	SuccessLaunches int                   `json:"success_launches"`
	FailureLaunches int                   `json:"failure_launches"`
	TotalLandings   int                   `json:"total_landings"`
	SuccessLandings int                   `json:"success_landings"`
	FailureLandings int                   `json:"failure_landings"`
}
