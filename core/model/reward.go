package model

import "time"

// RewardType wraps the reward type
type RewardType struct {
	ID          string    `json:"id" bson:"_id"`
	OrgID       string    `json:"org_id" bson:"org_id"`
	RewardType  string    `json:"reward_type" bson:"reward_type"`   // tshirt
	DisplayName string    `json:"display_name" bson:"display_name"` //
	Active      bool      `json:"active" bson:"active"`
	DateCreated time.Time `json:"date_created" bson:"date_created"`
	DateUpdated time.Time `json:"date_updated" bson:"date_updated"`
} // @name RewardType

// RewardInventory defines physical amount (availability) of a single award type
type RewardInventory struct {
	ID          string    `json:"id" bson:"_id"`
	OrgID       string    `json:"org_id" bson:"org_id"`
	RewardType  string    `json:"reward_type" bson:"reward_type"` // t-shirt
	Amount      int64     `json:"amount" bson:"amount"`
	InStock     bool      `json:"in_stock" bson:"in_stock"`
	Depleted    bool      `json:"depleted" bson:"depleted"`
	DateCreated time.Time `json:"date_created" bson:"date_created"`
	DateUpdated time.Time `json:"date_updated" bson:"date_updated"`
} // @name RewardInventory

// RewardHistoryEntry wraps the history entry
type RewardHistoryEntry struct {
	ID          string    `json:"id" bson:"_id"`
	OrgID       string    `json:"org_id" bson:"org_id"`
	UserID      string    `json:"user_id" bson:"user_id"`
	RewardType  string    `json:"reward_type" bson:"reward_type"`
	Amount      int64     `json:"amount" bson:"amount"`
	Description string    `json:"description" bson:"description"`
	DateCreated time.Time `json:"date_created" bson:"date_created"`
	DateUpdated time.Time `json:"date_updated" bson:"date_updated"`
	DisplayName *string   `json:"display_name,omitempty" bson:"display_name,omitempty"`
} // @name RewardHistoryEntry
