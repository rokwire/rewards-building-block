package model

import "time"

// RewardType wraps the reward type
type RewardType struct {
	ID          string    `json:"id" bson:"_id"`
	OrgID       string    `json:"org_id" bson:"org_id"`
	RewardType  string    `json:"reward_type" bson:"reward_type"`   // tshirt
	DisplayName string    `json:"display_name" bson:"display_name"` //
	Active      bool      `json:"active" bson:"active"`
	Description string    `json:"description" bson:"description"`
	DateCreated time.Time `json:"date_created" bson:"date_created"`
	DateUpdated time.Time `json:"date_updated" bson:"date_updated"`
} // @name RewardType

// RewardOperation wraps reward operation (defines amount of reward, BB and the type)
type RewardOperation struct {
	ID            string    `json:"id" bson:"_id"`
	OrgID         string    `json:"org_id" bson:"org_id"`
	RewardType    string    `json:"reward_type" bson:"reward_type"` // tshirt
	Code          string    `json:"code" bson:"code"`               //
	BuildingBlock string    `json:"building_block" bson:"building_block"`
	Amount        int64     `json:"amount" bson:"amount"`
	Description   string    `json:"description" bson:"description"`
	DateCreated   time.Time `json:"date_created" bson:"date_created"`
	DateUpdated   time.Time `json:"date_updated" bson:"date_updated"`
} // @name RewardOperation

// RewardInventory defines physical amount (availability) of a single award type
type RewardInventory struct {
	ID          string    `json:"id" bson:"_id"`
	OrgID       string    `json:"org_id" bson:"org_id"`
	RewardType  string    `json:"reward_type" bson:"reward_type"` // t-shirt
	Amount      int64     `json:"amount" bson:"amount"`
	InStock     bool      `json:"in_stock" bson:"in_stock"`
	Depleted    bool      `json:"depleted" bson:"depleted"`
	Description string    `json:"description" bson:"description"`
	DateCreated time.Time `json:"date_created" bson:"date_created"`
	DateUpdated time.Time `json:"date_updated" bson:"date_updated"`
} // @name RewardInventory

// Reward wraps the history entry
type Reward struct {
	ID            string    `json:"id" bson:"_id"`
	OrgID         string    `json:"org_id" bson:"org_id"`
	UserID        string    `json:"user_id" bson:"user_id"`
	RewardType    string    `json:"reward_type" bson:"reward_type"`
	Code          string    `json:"code" bson:"code"`
	BuildingBlock string    `json:"building_block" bson:"building_block"`
	Amount        int64     `json:"amount" bson:"amount"`
	Description   string    `json:"description" bson:"description"`
	DateCreated   time.Time `json:"date_created" bson:"date_created"`
	DateUpdated   time.Time `json:"date_updated" bson:"date_updated"`
} // @name Reward

// RewardQuantityState wraps current reward inventory state
type RewardQuantityState struct {
	RewardType         string `json:"reward_type" bson:"reward_type"`
	RewardableQuantity int64  `json:"rewardable_quantity" bson:"rewardable_quantity"`
	ClaimableQuantity  int64  `json:"claimable_quantity" bson:"claimable_quantity"`
}

// RewardClaim wraps a claim that is made by a user
type RewardClaim struct {
	ID          string            `json:"id" bson:"_id"`
	OrgID       string            `json:"org_id" bson:"org_id"`
	UserID      string            `json:"user_id" bson:"user_id"`
	Items       []RewardClaimItem `json:"items" bson:"items"`
	Status      string            `json:"status" bson:"status"`
	Description string            `json:"description" bson:"description"`
	DateCreated time.Time         `json:"date_created" bson:"date_created"`
	DateUpdated time.Time         `json:"date_updated" bson:"date_updated"`
} // @name RewardClaim

// RewardClaimItem wraps a claim  entry that consists reward type and amount
type RewardClaimItem struct {
	RewardType string `json:"reward_type" bson:"reward_type"`
	Amount     int64  `json:"amount" bson:"amount"`
} // @name RewardClaimItem
