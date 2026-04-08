package models

type LL2LauncherResponse struct {
	Count    int                         `json:"count" bson:"count"`
	Next     string                      `json:"next" bson:"next"`
	Previous string                      `json:"previous" bson:"previous"`
	Results  []LL2LauncherConfigDetailed `json:"results" bson:"results"`
}

type LL2LauncherConfigList struct {
	ResponseMode string                        `json:"response_mode" bson:"response_mode"`
	ID           int                           `json:"id" bson:"id"`
	URL          string                        `json:"url" bson:"url"`
	Name         string                        `json:"name" bson:"name"`
	Families     []LL2LauncherConfigFamilyMini `json:"families" bson:"families"`
	FullName     string                        `json:"full_name" bson:"full_name"`
	Variant      string                        `json:"variant" bson:"variant"`
}

type LL2LauncherConfigNormal struct {
	LL2LauncherConfigList `bson:",inline"`
	Active                bool               `json:"active" bson:"active"`
	IsPlaceholder         bool               `json:"is_placeholder" bson:"is_placeholder"`
	Manufacturer          LL2AgencyNormal    `json:"manufacturer" bson:"manufacturer"`
	Program               []LL2ProgramNormal `json:"program" bson:"program"`
	Reusable              bool               `json:"reusable" bson:"reusable"`
	Image                 LL2Image           `json:"image" bson:"image"`
	InfoURL               string             `json:"info_url" bson:"info_url"`
	WikiURL               string             `json:"wiki_url" bson:"wiki_url"`
}

type LL2LauncherConfigDetailed struct {
	LL2LauncherConfigNormal       `bson:",inline"`
	Description                   string  `json:"description" bson:"description"`
	Alias                         string  `json:"alias" bson:"alias"`
	MinStage                      int     `json:"min_stage" bson:"min_stage"`
	MaxStage                      int     `json:"max_stage" bson:"max_stage"`
	Length                        float64 `json:"length" bson:"length"`
	Diameter                      float64 `json:"diameter" bson:"diameter"`
	MaidenFlight                  string  `json:"maiden_flight" bson:"maiden_flight"`
	LaunchCost                    int     `json:"launch_cost" bson:"launch_cost"`
	LaunchMass                    float64 `json:"launch_mass" bson:"launch_mass"`
	LeoCapacity                   float64 `json:"leo_capacity" bson:"leo_capacity"`
	GtoCapacity                   float64 `json:"gto_capacity" bson:"gto_capacity"`
	GeoCapacity                   float64 `json:"geo_capacity" bson:"geo_capacity"`
	SsoCapacity                   float64 `json:"sso_capacity" bson:"sso_capacity"`
	ToThrust                      float64 `json:"to_thrust" bson:"to_thrust"`
	Apogee                        float64 `json:"apogee" bson:"apogee"`
	TotalLaunchCount              int     `json:"total_launch_count" bson:"total_launch_count"`
	ConsecutiveSuccessfulLaunches int     `json:"consecutive_successful_launches" bson:"consecutive_successful_launches"`
	SuccessfulLaunches            int     `json:"successful_launches" bson:"successful_launches"`
	FailedLaunches                int     `json:"failed_launches" bson:"failed_launches"`
	PendingLaunches               int     `json:"pending_launches" bson:"pending_launches"`
	AttemptedLandings             int     `json:"attempted_landings" bson:"attempted_landings"`
	SuccessfulLandings            int     `json:"successful_landings" bson:"successful_landings"`
	FailedLandings                int     `json:"failed_landings" bson:"failed_landings"`
	ConsecutiveSuccessfulLandings int     `json:"consecutive_successful_landings" bson:"consecutive_successful_landings"`
	FastestTurnaround             string  `json:"fastest_turnaround" bson:"fastest_turnaround"`
}

type LL2LauncherFamilyResponse struct {
	Count    int                               `json:"count" bson:"count"`
	Next     string                            `json:"next" bson:"next"`
	Previous string                            `json:"previous" bson:"previous"`
	Results  []LL2LauncherConfigFamilyDetailed `json:"results" bson:"results"`
}

type LL2LauncherConfigFamilyMini struct {
	ResponseMode string `json:"response_mode" bson:"response_mode"`
	ID           int    `json:"id" bson:"id"`
	Name         string `json:"name" bson:"name"`
}

type LL2LauncherConfigFamilyNormal struct {
	LL2LauncherConfigFamilyMini `bson:",inline"`
	Manufacturer                []LL2AgencyNormal           `json:"manufacturer" bson:"manufacturer"`
	Parent                      LL2LauncherConfigFamilyMini `json:"parent" bson:"parent"`
}

type LL2LauncherConfigFamilyDetailed struct {
	LL2LauncherConfigFamilyNormal `bson:",inline"`
	Description                   string `json:"description" bson:"description"`
	Active                        bool   `json:"active" bson:"active"`
	MaidenFlight                  string `json:"maiden_flight" bson:"maiden_flight"`
	TotalLaunchCount              int    `json:"total_launch_count" bson:"total_launch_count"`
	ConsecutiveSuccessfulLaunches int    `json:"consecutive_successful_launches" bson:"consecutive_successful_launches"`
	SuccessfulLaunches            int    `json:"successful_launches" bson:"successful_launches"`
	FailedLaunches                int    `json:"failed_launches" bson:"failed_launches"`
	PendingLaunches               int    `json:"pending_launches" bson:"pending_launches"`
	AttemptedLandings             int    `json:"attempted_landings" bson:"attempted_landings"`
	SuccessfulLandings            int    `json:"successful_landings" bson:"successful_landings"`
	FailedLandings                int    `json:"failed_landings" bson:"failed_landings"`
	ConsecutiveSuccessfulLandings int    `json:"consecutive_successful_landings" bson:"consecutive_successful_landings"`
}
