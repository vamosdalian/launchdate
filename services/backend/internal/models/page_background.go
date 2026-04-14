package models

import "time"

type PageBackgroundDefinition struct {
	Key         string
	DisplayName string
}

const (
	PageBackgroundKeyHome        = "home"
	PageBackgroundKeyLaunches    = "launches"
	PageBackgroundKeyRockets     = "rockets"
	PageBackgroundKeyLaunchBases = "launch-bases"
	PageBackgroundKeyCompanies   = "companies"
)

var PageBackgroundDefinitions = []PageBackgroundDefinition{
	{Key: PageBackgroundKeyHome, DisplayName: "Home"},
	{Key: PageBackgroundKeyLaunches, DisplayName: "Launches"},
	{Key: PageBackgroundKeyRockets, DisplayName: "Rockets"},
	{Key: PageBackgroundKeyLaunchBases, DisplayName: "Launch Bases"},
	{Key: PageBackgroundKeyCompanies, DisplayName: "Companies"},
}

type PageBackground struct {
	ID              int64      `bson:"id,omitempty" json:"id,omitempty,string"`
	PageKey         string     `bson:"page_key" json:"page_key"`
	BackgroundImage string     `bson:"background_image" json:"background_image"`
	DisplayName     string     `bson:"-" json:"display_name"`
	Configured      bool       `bson:"-" json:"configured"`
	CreatedAt       *time.Time `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt       *time.Time `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

type PageBackgroundUpdateRequest struct {
	BackgroundImage string `json:"background_image"`
}

func IsValidPageBackgroundKey(key string) bool {
	for _, definition := range PageBackgroundDefinitions {
		if definition.Key == key {
			return true
		}
	}

	return false
}

func PageBackgroundDisplayName(key string) string {
	for _, definition := range PageBackgroundDefinitions {
		if definition.Key == key {
			return definition.DisplayName
		}
	}

	return key
}
