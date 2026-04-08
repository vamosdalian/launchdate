package models

import "time"

type ImageList struct {
	Count  int     `json:"count"`
	Images []Image `json:"images"`
}

type Image struct {
	ID          int64        `bson:"id" json:"id,string"`
	Key         string       `bson:"key" json:"key"`
	URL         string       `bson:"-" json:"url"`
	Name        string       `bson:"name" json:"name"`
	Width       int          `bson:"width" json:"width"`
	Height      int          `bson:"height" json:"height"`
	Size        int64        `bson:"size" json:"size"`
	ContentType string       `bson:"content_type" json:"content_type"`
	UploadTime  time.Time    `bson:"upload_time" json:"upload_time"`
	ThumbImages []ThumbImage `bson:"thumb_images" json:"thumb_images"`
}

type ThumbImage struct {
	Key         string    `bson:"key" json:"key"`
	URL         string    `bson:"-" json:"url"`
	Width       int       `bson:"width" json:"width"`
	Height      int       `bson:"height" json:"height"`
	Size        int64     `bson:"size" json:"size"`
	ContentType string    `bson:"content_type" json:"content_type"`
	UploadTime  time.Time `bson:"upload_time" json:"upload_time"`
}
