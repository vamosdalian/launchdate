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
