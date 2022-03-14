/*
 *   Copyright (c) 2020 Board of Trustees of the University of Illinois.
 *   All rights reserved.

 *   Licensed under the Apache License, Version 2.0 (the "License");
 *   you may not use this file except in compliance with the License.
 *   You may obtain a copy of the License at

 *   http://www.apache.org/licenses/LICENSE-2.0

 *   Unless required by applicable law or agreed to in writing, software
 *   distributed under the License is distributed on an "AS IS" BASIS,
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *   See the License for the specific language governing permissions and
 *   limitations under the License.
 */

package storage

import (
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"rewards/core/model"
	"strconv"
	"time"
)

// Adapter implements the Storage interface
type Adapter struct {
	db *database
}

// Start starts the storage
func (sa *Adapter) Start() error {
	err := sa.db.start()
	return err
}

// NewStorageAdapter creates a new storage adapter instance
func NewStorageAdapter(mongoDBAuth string, mongoDBName string, mongoTimeout string) *Adapter {
	timeout, err := strconv.Atoi(mongoTimeout)
	if err != nil {
		log.Println("Set default timeout - 500")
		timeout = 500
	}
	timeoutMS := time.Millisecond * time.Duration(timeout)

	db := &database{mongoDBAuth: mongoDBAuth, mongoDBName: mongoDBName, mongoTimeout: timeoutMS}
	return &Adapter{db: db}
}

// GetRewardTypes Gets all reward types
func (sa *Adapter) GetRewardTypes(orgID string) ([]model.RewardType, error) {
	filter := bson.D{
		primitive.E{Key: "org_id", Value: orgID},
	}
	var result []model.RewardType
	err := sa.db.rewardTypes.Find(filter, &result, nil)
	if err != nil {
		log.Printf("storage.GetRewardTypes error: %s", err)
		return nil, fmt.Errorf("storage.GetRewardTypes error: %s", err)
	}
	if result == nil {
		result = []model.RewardType{}
	}
	return result, nil
}

// GetRewardType Gets a reward type by id
func (sa *Adapter) GetRewardType(orgID string, id string) (*model.RewardType, error) {
	filter := bson.D{
		primitive.E{Key: "org_id", Value: orgID},
		primitive.E{Key: "_id", Value: id},
	}
	var result []model.RewardType
	err := sa.db.rewardTypes.Find(filter, &result, nil)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result) == 0 {
		log.Printf("storage.UpdateRewardType error: %s", err)
		return nil, fmt.Errorf("storage.UpdateRewardType error: %s", err)
	}
	return &result[0], nil
}

// GetRewardTypeByType Gets a reward type by type
func (sa *Adapter) GetRewardTypeByType(orgID string, rewardType string) (*model.RewardType, error) {
	filter := bson.D{
		primitive.E{Key: "org_id", Value: orgID},
		primitive.E{Key: "reward_type", Value: rewardType},
	}
	var result []model.RewardType
	err := sa.db.rewardTypes.Find(filter, &result, nil)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result) == 0 {
		log.Printf("storage.GetRewardTypeByType error: %s", err)
		return nil, fmt.Errorf("storage.GetRewardTypeByType error: %s", err)
	}
	return &result[0], nil
}

// CreateRewardType creates a new reward type
func (sa *Adapter) CreateRewardType(orgID string, item model.RewardType) (*model.RewardType, error) {
	now := time.Now().UTC()
	item.ID = uuid.NewString()
	item.OrgID = orgID
	item.DateCreated = now
	item.DateUpdated = now
	_, err := sa.db.rewardTypes.InsertOne(&item)
	if err != nil {
		log.Printf("storage.GetRewardTypes error: %s", err)
		return nil, fmt.Errorf("storage.GetRewardTypes error: %s", err)
	}
	return &item, nil
}

// UpdateRewardType updates a reward type
func (sa *Adapter) UpdateRewardType(orgID string, id string, item model.RewardType) (*model.RewardType, error) {
	jsonID := item.ID
	if jsonID != id {
		return nil, fmt.Errorf("storage.UpdateRewardType attempt to override another object")
	}

	now := time.Now().UTC()
	filter := bson.D{
		primitive.E{Key: "_id", Value: id},
		primitive.E{Key: "org_id", Value: orgID},
	}
	update := bson.D{
		primitive.E{Key: "$set", Value: bson.D{
			primitive.E{Key: "display_name", Value: item.DisplayName},
			primitive.E{Key: "active", Value: item.Active},
			primitive.E{Key: "date_updated", Value: now},
		},
		},
	}
	_, err := sa.db.rewardInventories.UpdateOne(filter, update, nil)
	if err != nil {
		log.Printf("storage.UpdateRewardType error: %s", err)
		return nil, fmt.Errorf("storage.UpdateRewardType error: %s", err)
	}

	item.DateUpdated = now

	return &item, nil
}

// DeleteRewardType deletes a reward type
func (sa *Adapter) DeleteRewardType(orgID string, id string) error {
	// TBD check and deny if the reward type is in use!!!

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	_, err := sa.db.rewardInventories.DeleteOne(filter, nil)
	if err != nil {
		log.Printf("storage.DeleteRewardType error: %s", err)
		return fmt.Errorf("storage.DeleteRewardType error: %s", err)
	}

	return nil
}

// GetRewardInventories Gets all reward inventories
func (sa *Adapter) GetRewardInventories(orgID string, ids []string, rewardType *string) ([]model.RewardInventory, error) {
	filter := bson.D{}
	if len(ids) > 0 {
		filter = bson.D{
			primitive.E{Key: "_id", Value: bson.M{"$in": ids}},
		}
	}

	if rewardType != nil {
		filter = bson.D{
			primitive.E{Key: "reward_type", Value: *rewardType},
		}
	}

	var result []model.RewardInventory
	err := sa.db.rewardInventories.Find(filter, &result, &options.FindOptions{
		Sort: bson.D{{"date_created", 1}},
	})
	if err != nil {
		log.Printf("storage.GetRewardInventories error: %s", err)
		return nil, fmt.Errorf("storage.GetRewardInventories error: %s", err)
	}
	return result, nil
}

// GetRewardInventory Gets a reward inventory by id
func (sa *Adapter) GetRewardInventory(orgID string, id string) (*model.RewardInventory, error) {
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	var result []model.RewardInventory
	err := sa.db.rewardInventories.Find(filter, &result, nil)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result) == 0 {
		log.Printf("storage.GetRewardInventory error: %s", err)
		return nil, fmt.Errorf("storage.GetRewardInventory error: %s", err)
	}
	return &result[0], nil
}

// CreateRewardInventory creates a new reward inventory
func (sa *Adapter) CreateRewardInventory(orgID string, item model.RewardInventory) (*model.RewardInventory, error) {
	now := time.Now().UTC()
	item.ID = uuid.NewString()
	item.DateCreated = now
	item.DateUpdated = now
	_, err := sa.db.rewardInventories.InsertOne(&item)
	if err != nil {
		log.Printf("storage.CreateRewardInventory error: %s", err)
		return nil, fmt.Errorf("storage.CreateRewardInventory error: %s", err)
	}
	return &item, nil
}

// UpdateRewardInventory updates a reward pool
func (sa *Adapter) UpdateRewardInventory(orgID string, id string, item model.RewardInventory) (*model.RewardInventory, error) {
	jsonID := item.ID
	if jsonID != id {
		return nil, fmt.Errorf("storage.UpdateRewardInventory attempt to override another object")
	}

	now := time.Now().UTC()
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{
		primitive.E{Key: "$set", Value: bson.D{
			primitive.E{Key: "date_updated", Value: now},
			primitive.E{Key: "amount", Value: item.Amount},
			primitive.E{Key: "in_stock", Value: item.InStock},
			primitive.E{Key: "depleted", Value: item.Depleted},
			primitive.E{Key: "description", Value: item.Description},
		},
		},
	}
	_, err := sa.db.rewardInventories.UpdateOne(filter, update, nil)
	if err != nil {
		log.Printf("storage.UpdateRewardInventory error: %s", err)
		return nil, fmt.Errorf("storage.UpdateRewardInventory error: %s", err)
	}

	item.DateUpdated = now

	return &item, nil
}

// DeleteRewardInventory deletes a reward pool. Don't delete if it's in use!
func (sa *Adapter) DeleteRewardInventory(orgID string, id string) error {
	// TBD check and deny if the reward type is in use!!!

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	_, err := sa.db.rewardInventories.DeleteOne(filter, nil)
	if err != nil {
		log.Printf("storage.DeleteRewardInventory error: %s", err)
		return fmt.Errorf("storage.DeleteRewardInventory error: %s", err)
	}

	return nil
}

// GetRewardHistoryEntries Gets all reward history entries
func (sa *Adapter) GetRewardHistoryEntries(orgID string, userID string) ([]model.Reward, error) {
	filter := bson.D{
		primitive.E{Key: "user_id", Value: userID},
	}

	var result []model.Reward
	err := sa.db.rewardHistory.Find(filter, &result, &options.FindOptions{
		Sort: bson.D{{"date_created", -1}},
	})
	if err != nil {
		log.Printf("storage.GetRewardHistoryEntries error: %s", err)
		return nil, fmt.Errorf("storage.GetRewardHistoryEntries error: %s", err)
	}
	if result == nil {
		result = []model.Reward{}
	}
	return result, nil
}

// GetRewardHistoryEntry Gets a reward history entry by id
func (sa *Adapter) GetRewardHistoryEntry(orgID string, userID, id string) (*model.Reward, error) {
	filter := bson.D{
		primitive.E{Key: "_id", Value: id},
		primitive.E{Key: "user_id", Value: userID},
	}
	var result []model.Reward
	err := sa.db.rewardHistory.Find(filter, &result, nil)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result) == 0 {
		log.Printf("storage.GetRewardHistoryEntry error: %s", err)
		return nil, fmt.Errorf("storage.GetRewardHistoryEntry error: %s", err)
	}
	return &result[0], nil
}

// CreateReward creates a new reward history entry
func (sa *Adapter) CreateReward(orgID string, item model.Reward) (*model.Reward, error) {
	now := time.Now().UTC()
	item.ID = uuid.NewString()
	item.DateCreated = now
	item.DateUpdated = now
	_, err := sa.db.rewardHistory.InsertOne(&item)
	if err != nil {
		log.Printf("storage.CreateReward error: %s", err)
		return nil, fmt.Errorf("storage.CreateReward error: %s", err)
	}
	return &item, nil
}

// GetUserBalance Gets all balances for the user id
func (sa *Adapter) GetUserBalance(orgID string, userID string) (*model.WalletBalance, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"user_id": userID}},
		{"$group": bson.M{"_id": "$code", "amount": bson.M{"$sum": "$amount"}}},
	}

	var result []model.WalletBalance
	err := sa.db.rewardHistory.Aggregate(pipeline, &result, nil)
	if err != nil {
		log.Printf("storage.GetRewardHistoryEntries error: %s", err)
		return nil, fmt.Errorf("storage.GetRewardHistoryEntries error: %s", err)
	}
	if len(result) > 0 {
		return &result[0], nil
	}
	return &model.WalletBalance{
		Amount: 0,
	}, nil
}

// GetWalletBalance gets wallet balance by user id and wallet code
func (sa *Adapter) GetWalletBalance(orgID string, userID string, code string) (*model.WalletBalance, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"user_id": userID, "code": code}},
		{"$group": bson.M{"_id": "$code", "amount": bson.M{"$sum": "$amount"}}},
	}

	var result []model.WalletBalance
	err := sa.db.rewardHistory.Aggregate(pipeline, &result, nil)
	if err != nil {
		log.Printf("storage.GetRewardHistoryEntries error: %s", err)
		return nil, fmt.Errorf("storage.GetRewardHistoryEntries error: %s", err)
	}

	if len(result) > 0 {
		return &result[0], nil
	}

	return nil, nil
}

// SetListener sets the upper layer storage listener for sending collection changed callbacks
func (sa *Adapter) SetListener(listener Listener) {
	sa.db.listener = listener
}

// Event

func (m *database) onDataChanged(changeDoc map[string]interface{}) {
	if changeDoc == nil {
		return
	}
	log.Printf("onDataChanged: %+v\n", changeDoc)
	ns := changeDoc["ns"]
	if ns == nil {
		return
	}
	nsMap := ns.(map[string]interface{})
	coll := nsMap["coll"]

	if "reward_types" == coll {
		if m.listener != nil {
			m.listener.OnRewardTypesChanged()
		}
	}
}
