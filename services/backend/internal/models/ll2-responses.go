package models

type LL2LaunchList struct {
	Count    int               `json:"count"`
	Launches []LL2LaunchNormal `json:"launches"`
}

type LL2AgencyList struct {
	Count    int                 `json:"count"`
	Agencies []LL2AgencyDetailed `json:"agencies"`
}

type LL2LauncherList struct {
	Count     int                       `json:"count"`
	Launchers []LL2LauncherConfigNormal `json:"launchers"`
}

type LL2LauncherFamilyList struct {
	Count    int                               `json:"count"`
	Families []LL2LauncherConfigFamilyDetailed `json:"families"`
}

type LL2LocationList struct {
	Count     int                             `json:"count"`
	Locations []LL2LocationSerializerWithPads `json:"locations"`
}

type LL2PadList struct {
	Count int      `json:"count"`
	Pads  []LL2Pad `json:"pads"`
}
