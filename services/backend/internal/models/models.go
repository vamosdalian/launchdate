package models

type LL2Stats struct {
	Launches         int `json:"launches"`
	Agencies         int `json:"agencies"`
	Launchers        int `json:"launchers"`
	LauncherFamilies int `json:"launcher_families"`
	Locations        int `json:"locations"`
	Pads             int `json:"pads"`
}

type Stats struct {
	Rockets     int      `json:"rocket"`
	Launches    int      `json:"launch"`
	Agencies    int      `json:"agency"`
	LaunchBases int      `json:"launch_base"`
	LL2         LL2Stats `json:"ll2"`
}
