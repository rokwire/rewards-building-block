package model

// RewardTypeAmount wraps the balance aggregation response
type RewardTypeAmount struct {
	RewardType string `json:"reward_type" bson:"_id"`
	Amount     int64  `json:"amount" bson:"amount"`
} // @name RewardTypeAmount
