package models

// for DB
type LaunchBase struct {
	ID         int64 `bson:"id"`
	ExternalID int64 `bson:"external_id"` // ll2 launch id
}

type LaunchBaseSerializer struct {
	ID   int64                         `json:"id"`
	Data LL2LocationSerializerWithPads `json:"data"`
}

type LaunchBaseList struct {
	Count    int                    `json:"count"`
	Launches []LaunchBaseSerializer `json:"launches"`
}

type PublicCompactLocation struct {
	ID   int64   `json:"id,string"`
	Name string  `json:"name"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
}
