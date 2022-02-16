package model

// WalletBalance wraps the balance aggregation response
type WalletBalance struct {
	Amount int64  `json:"amount" bson:"amount"`
} // @name WalletBalance
