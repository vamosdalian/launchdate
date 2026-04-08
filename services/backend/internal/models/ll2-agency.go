package models

type LL2AngecyResponse struct {
	Count    int                  `json:"count"`
	Next     string               `json:"next"`
	Previous string               `json:"previous"`
	Results  []*LL2AgencyDetailed `json:"results"`
}

type LL2AgencyMini struct {
	ResponseMode string        `json:"response_mode" bson:"response_mode"`
	ID           int           `json:"id" bson:"id"`
	URL          string        `json:"url" bson:"url"`
	Name         string        `json:"name" bson:"name"`
	Abbrev       string        `json:"abbrev" bson:"abbrev"`
	Type         LL2AgencyType `json:"type" bson:"type"`
}

type LL2AgencyType struct {
	ID   int    `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
}

type LL2AgencyNormal struct {
	LL2AgencyMini `bson:",inline"`
	Featured      bool         `json:"featured" bson:"featured"`
	Country       []LL2Country `json:"country" bson:"country"`
	Description   string       `json:"description" bson:"description"`
	Administrator string       `json:"administrator" bson:"administrator"`
	FoundingYear  int          `json:"founding_year" bson:"founding_year"`
	Launchers     string       `json:"launchers" bson:"launchers"`
	Spacecraft    string       `json:"spacecraft" bson:"spacecraft"`
	Parent        string       `json:"parent" bson:"parent"`
	Image         LL2Image     `json:"image" bson:"image"`
	Logo          LL2Image     `json:"logo" bson:"logo"`
	SocialLogo    LL2Image     `json:"social_logo" bson:"social_logo"`
}

type LL2AgencyDetailed struct {
	LL2AgencyNormal               `bson:",inline"`
	TotalLaunchCount              int                  `json:"total_launch_count" bson:"total_launch_count"`
	ConsecutiveSuccessfulLaunches int                  `json:"consecutive_successful_launches" bson:"consecutive_successful_launches"`
	SuccessfulLaunches            int                  `json:"successful_launches" bson:"successful_launches"`
	FailedLaunches                int                  `json:"failed_launches" bson:"failed_launches"`
	PendingLaunches               int                  `json:"pending_launches" bson:"pending_launches"`
	ConsecutiveSuccessfulLandings int                  `json:"consecutive_successful_landings" bson:"consecutive_successful_landings"`
	SuccessfulLandings            int                  `json:"successful_landings" bson:"successful_landings"`
	FailedLandings                int                  `json:"failed_landings" bson:"failed_landings"`
	AttemptedLandings             int                  `json:"attempted_landings" bson:"attempted_landings"`
	SuccessfulLandingsSpacecraft  int                  `json:"successful_landings_spacecraft" bson:"successful_landings_spacecraft"`
	FailedLandingsSpacecraft      int                  `json:"failed_landings_spacecraft" bson:"failed_landings_spacecraft"`
	AttemptedLandingsSpacecraft   int                  `json:"attempted_landings_spacecraft" bson:"attempted_landings_spacecraft"`
	SuccessfulLandingsPayload     int                  `json:"successful_landings_payload" bson:"successful_landings_payload"`
	FailedLandingsPayload         int                  `json:"failed_landings_payload" bson:"failed_landings_payload"`
	AttemptedLandingsPayload      int                  `json:"attempted_landings_payload" bson:"attempted_landings_payload"`
	InfoURL                       string               `json:"info_url" bson:"info_url"`
	WikiURL                       string               `json:"wiki_url" bson:"wiki_url"`
	SocialMediaLinks              []LL2SocialMediaLink `json:"social_media_links" bson:"social_media_links"`
	// TODO launcher_list
	// TODO spacecraft_list
}
