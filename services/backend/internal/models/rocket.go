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
