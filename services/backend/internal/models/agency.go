package models

// for DB
type Agency struct {
	ID         int64       `json:"id,string" bson:"id"`
	ExternalID int64       `json:"external_id" bson:"external_id"` //ll2 agency id
	ThumbImage string      `json:"thumb_image" bson:"thumb_image"`
	Images     []string    `json:"images" bson:"images"`
	SocialUrl  []SocialUrl `json:"social_url" bson:"social_url"`
	ShowOnHome bool        `json:"show_on_home" bson:"show_on_home"`
}

type SocialUrl struct {
	Name string `json:"name" bson:"name"`
	URL  string `json:"url" bson:"url"`
}

type AgencySerializer struct {
	Agency
	Data LL2AgencyNormal `json:"data"`
}

type AgencyList struct {
	Count    int                `json:"count"`
	Agencies []AgencySerializer `json:"agencies"`
}
