package model

import "time"

// RewardType wraps the reward type
type RewardType struct {
	ID            string    `json:"id" bson:"_id"`
	RewardType    string    `json:"reward_type" bson:"reward_type"`       // illini_cash
	Name          string    `json:"name" bson:"name"`                     // "win_five_point_by_five_readings"
	DisplayName   string    `json:"display_name" bson:"display_name"`     // Win five point by five readings
	BuildingBlock string    `json:"building_block" bson:"building_block"` // "content"
	Amount        int64     `json:"amount" bson:"amount"`                 // 5
	Active        bool      `json:"active" bson:"active"`
	DateCreated   time.Time `json:"date_created" bson:"date_created"`
	DateUpdated   time.Time `json:"date_updated" bson:"date_updated"`
} // @name RewardType

// RewardPool wraps reward agent (Amazon, Illini Cash etc) with credits
type RewardPool struct {
	ID          string    `json:"id" bson:"_id"`
	RewardCode  string    `json:"reward_code" bson:"reward_code"` // illini_cash
	Name        string    `json:"name" bson:"name"`
	Data        JSONData  `json:"data" bson:"data"`
	Amount      int64     `json:"amount" bson:"amount"`
	Active      bool      `json:"active" bson:"active"`
	DateCreated time.Time `json:"date_created" bson:"date_created"`
	DateUpdated time.Time `json:"date_updated" bson:"date_updated"`
} // @name RewardPool

// RewardHistoryEntry wraps the history entry
type RewardHistoryEntry struct {
	ID          string    `json:"id" bson:"_id"`
	UserID      string    `json:"user_id" bson:"user_id"`
	RewardType  string    `json:"reward_type" bson:"reward_type"`
	Amount      int64     `json:"amount" bson:"amount"`
	Description string    `json:"description" bson:"description"`
	DateCreated time.Time `json:"date_created" bson:"date_created"`
	DateUpdated time.Time `json:"date_updated" bson:"date_updated"`
} // @name RewardHistoryEntry
