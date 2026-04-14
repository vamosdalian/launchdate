package models

// for DB
type LaunchBase struct {
	ID         int64 `bson:"id" json:"id,string"`
	ExternalID int64 `bson:"external_id"` // ll2 launch id
}

type LaunchBaseSerializer struct {
	ID   int64                         `json:"id,string"`
	Data LL2LocationSerializerWithPads `json:"data"`
}

type LaunchBaseList struct {
	Count    int                    `json:"count"`
	Launches []LaunchBaseSerializer `json:"launches"`
}
