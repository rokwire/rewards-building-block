package model

import "time"

// RewardType wraps the reward type
type RewardType struct {
	ID            string    `json:"id" bson:"_id"`
	Name          string    `json:"name" bson:"name"`
	DisplayName   string    `json:"display_name" bson:"display_name"`
	BuildingBlock string    `json:"building_block" bson:"building_block"`
	Amount        int       `json:"amount" bson:"amount"`
	Active        bool      `json:"active" bson:"active"`
	DateCreated   time.Time `json:"date_created" bson:"date_created"`
	DateUpdated   time.Time `json:"date_updated" bson:"date_updated"`
} //@name RewardType

// RewardHistoryEntry wraps the history entry
type RewardHistoryEntry struct {
	ID            string    `json:"id" bson:"_id"`
	Name          string    `json:"name" bson:"name"`
	BuildingBlock string    `json:"building_block" bson:"building_block"`
	Amount        int       `json:"amount" bson:"amount"`
	DateCreated   time.Time `json:"date_created" bson:"date_created"`
} //@name RewardHistoryEntry
