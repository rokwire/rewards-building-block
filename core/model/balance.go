package model

type WalletBalance struct {
	Code string `json:"code" bson:"_id"`
	Amount int64 `json:"amount" bson:"amount"`
}
