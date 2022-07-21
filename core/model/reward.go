// Copyright 2022 Board of Trustees of the University of Illinois.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package model

import "time"

// RewardType wraps the reward type
type RewardType struct {
	ID          string    `json:"id" bson:"_id"`
	OrgID       string    `json:"org_id" bson:"org_id"`
	AppID       string    `json:"app_id" bson:"app_id"`
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
	AppID         string    `json:"org_id" bson:"org_id"`
	RewardType    string    `json:"reward_type" bson:"reward_type"` // tshirt
	Code          string    `json:"code" bson:"code"`               //
	BuildingBlock string    `json:"building_block" bson:"building_block"`
	Amount        int       `json:"amount" bson:"amount"`
	Description   string    `json:"description" bson:"description"`
	DateCreated   time.Time `json:"date_created" bson:"date_created"`
	DateUpdated   time.Time `json:"date_updated" bson:"date_updated"`
} // @name RewardOperation

// RewardInventory defines physical amount (availability) of a single award type
type RewardInventory struct {
	ID            string    `json:"id" bson:"_id"`
	OrgID         string    `json:"org_id" bson:"org_id"`
	AppID         string    `json:"app_id" bson:"app_id"`
	RewardType    string    `json:"reward_type" bson:"reward_type"` // t-shirt
	InStock       bool      `json:"in_stock" bson:"in_stock"`
	AmountTotal   int       `json:"amount_total" bson:"amount_total"`
	AmountGranted int       `json:"amount_granted" bson:"amount_granted"`
	AmountClaimed int       `json:"amount_claimed" bson:"amount_claimed"`
	GrantDepleted bool      `json:"grant_depleted" bson:"grant_depleted"`
	ClaimDepleted bool      `json:"claim_depleted" bson:"claim_depleted"`
	Description   string    `json:"description" bson:"description"`
	DateCreated   time.Time `json:"date_created" bson:"date_created"`
	DateUpdated   time.Time `json:"date_updated" bson:"date_updated"`
} // @name RewardInventory

// GetGrantableAmount Gets grantable amount
func (ri *RewardInventory) GetGrantableAmount() int {
	return ri.AmountTotal - ri.AmountGranted
}

// GetClaimableAmount Gets claimable amount
func (ri *RewardInventory) GetClaimableAmount() int {
	return ri.AmountTotal - ri.AmountClaimed
}

// Reward wraps the history entry
type Reward struct {
	ID            string    `json:"id" bson:"_id"`
	OrgID         string    `json:"org_id" bson:"org_id"`
	AppID         string    `json:"app_id" bson:"app_id"`
	UserID        string    `json:"user_id" bson:"user_id"`
	RewardType    string    `json:"reward_type" bson:"reward_type"`
	Code          string    `json:"code" bson:"code"`
	BuildingBlock string    `json:"building_block" bson:"building_block"`
	Amount        int       `json:"amount" bson:"amount"`
	Description   string    `json:"description" bson:"description"`
	DateCreated   time.Time `json:"date_created" bson:"date_created"`
	DateUpdated   time.Time `json:"date_updated" bson:"date_updated"`
} // @name Reward

// RewardQuantityState wraps current reward inventory state
type RewardQuantityState struct {
	RewardType        string `json:"reward_type" bson:"reward_type"`
	GrantableQuantity int    `json:"grantable_quantity" bson:"grantable_quantity"`
	ClaimableQuantity int    `json:"claimable_quantity" bson:"claimable_quantity"`
}

// RewardClaim wraps a claim that is made by a user
type RewardClaim struct {
	ID          string            `json:"id" bson:"_id"`
	OrgID       string            `json:"org_id" bson:"org_id"`
	AppID       string            `json:"app_id" bson:"app_id"`
	UserID      string            `json:"user_id" bson:"user_id"`
	Items       []RewardClaimItem `json:"items" bson:"items"`
	Status      string            `json:"status" bson:"status"`
	Description string            `json:"description" bson:"description"`
	DateCreated time.Time         `json:"date_created" bson:"date_created"`
	DateUpdated time.Time         `json:"date_updated" bson:"date_updated"`
} // @name RewardClaim

// RewardClaimItem wraps a claim  entry that consists reward type and amount
type RewardClaimItem struct {
	RewardType  string `json:"reward_type" bson:"reward_type"`
	InventoryID string `json:"inventory_id" bson:"inventory_id"`

	Amount int `json:"amount" bson:"amount"`
} // @name RewardClaimItem
