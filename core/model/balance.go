package model

// WalletBalance wraps the balance aggregation response
type WalletBalance struct {
	Code   string `json:"code" bson:"_id"`
	Amount int64  `json:"amount" bson:"amount"`
} // @name WalletBalance
