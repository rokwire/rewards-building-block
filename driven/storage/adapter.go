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
func (sa *Adapter) GetRewardTypes(ids []string) ([]model.RewardType, error) {
	filter := bson.D{}
	if len(ids) > 0 {
		filter = bson.D{
			primitive.E{Key: "_id", Value: bson.M{"$in": ids}},
		}
	}

	var result []model.RewardType
	err := sa.db.rewardTypes.Find(filter, &result, nil)
	if err != nil {
		log.Printf("storage.GetRewardTypes error: %s", err)
		return nil, fmt.Errorf("storage.GetRewardTypes error: %s", err)
	}
	return result, nil
}

// GetRewardType Gets a reward type by id
func (sa *Adapter) GetRewardType(id string) (*model.RewardType, error) {
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
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

// CreateRewardType creates a new reward type
func (sa *Adapter) CreateRewardType(item model.RewardType) (*model.RewardType, error) {
	now := time.Now().UTC()
	item.ID = uuid.NewString()
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
func (sa *Adapter) UpdateRewardType(id string, item model.RewardType) (*model.RewardType, error) {
	jsonID := item.ID
	if jsonID != id {
		return nil, fmt.Errorf("storage.UpdateRewardType attempt to override another object")
	}

	now := time.Now().UTC()
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{
		primitive.E{Key: "$set", Value: bson.D{
			primitive.E{Key: "name", Value: item.Name},
			primitive.E{Key: "display_name", Value: item.DisplayName},
			primitive.E{Key: "building_block", Value: item.BuildingBlock},
			primitive.E{Key: "amount", Value: item.Amount},
			primitive.E{Key: "active", Value: item.Active},
			primitive.E{Key: "date_updated", Value: now},
		},
		},
	}
	_, err := sa.db.rewardPools.UpdateOne(filter, update, nil)
	if err != nil {
		log.Printf("storage.UpdateRewardType error: %s", err)
		return nil, fmt.Errorf("storage.UpdateRewardType error: %s", err)
	}

	item.DateUpdated = now

	return &item, nil
}

// DeleteRewardType deletes a reward type
func (sa *Adapter) DeleteRewardType(id string) error {
	// TBD check and deny if the reward type is in use!!!

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	_, err := sa.db.rewardPools.DeleteOne(filter, nil)
	if err != nil {
		log.Printf("storage.DeleteRewardType error: %s", err)
		return fmt.Errorf("storage.DeleteRewardType error: %s", err)
	}

	return nil
}

// GetRewardPools Gets all reward pools
func (sa *Adapter) GetRewardPools(ids []string) ([]model.RewardPool, error) {
	filter := bson.D{}
	if len(ids) > 0 {
		filter = bson.D{
			primitive.E{Key: "_id", Value: bson.M{"$in": ids}},
		}
	}

	var result []model.RewardPool
	err := sa.db.rewardPools.Find(filter, &result, nil)
	if err != nil {
		log.Printf("storage.GetRewardPools error: %s", err)
		return nil, fmt.Errorf("storage.GetRewardPools error: %s", err)
	}
	return result, nil
}

// GetRewardPool Gets a reward pool by id
func (sa *Adapter) GetRewardPool(id string) (*model.RewardPool, error) {
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	var result []model.RewardPool
	err := sa.db.rewardPools.Find(filter, &result, nil)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result) == 0 {
		log.Printf("storage.GetRewardPool error: %s", err)
		return nil, fmt.Errorf("storage.GetRewardPool error: %s", err)
	}
	return &result[0], nil
}

// CreateRewardPool creates a new reward pool
func (sa *Adapter) CreateRewardPool(item model.RewardPool) (*model.RewardPool, error) {
	now := time.Now().UTC()
	item.ID = uuid.NewString()
	item.DateCreated = now
	item.DateUpdated = now
	_, err := sa.db.rewardPools.InsertOne(&item)
	if err != nil {
		log.Printf("storage.CreateRewardPool error: %s", err)
		return nil, fmt.Errorf("storage.CreateRewardPool error: %s", err)
	}
	return &item, nil
}

// UpdateRewardPool updates a reward pool
func (sa *Adapter) UpdateRewardPool(id string, item model.RewardPool) (*model.RewardPool, error) {
	jsonID := item.ID
	if jsonID != id {
		return nil, fmt.Errorf("storage.UpdateRewardPool attempt to override another object")
	}

	now := time.Now().UTC()
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{
		primitive.E{Key: "$set", Value: bson.D{
			primitive.E{Key: "date_updated", Value: now},
			primitive.E{Key: "name", Value: item.Name},
			primitive.E{Key: "amount", Value: item.Amount},
			primitive.E{Key: "active", Value: item.Active},
			primitive.E{Key: "data", Value: item.Data},
		},
		},
	}
	_, err := sa.db.rewardPools.UpdateOne(filter, update, nil)
	if err != nil {
		log.Printf("storage.UpdateRewardPool error: %s", err)
		return nil, fmt.Errorf("storage.UpdateRewardPool error: %s", err)
	}

	item.DateUpdated = now

	return &item, nil
}

// DeleteRewardPool deletes a reward pool. Don't delete if it's in use!
func (sa *Adapter) DeleteRewardPool(id string) error {
	// TBD check and deny if the reward type is in use!!!

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	_, err := sa.db.rewardPools.DeleteOne(filter, nil)
	if err != nil {
		log.Printf("storage.DeleteRewardPool error: %s", err)
		return fmt.Errorf("storage.DeleteRewardPool error: %s", err)
	}

	return nil
}

// GetRewardHistoryEntries Gets all reward history entries
func (sa *Adapter) GetRewardHistoryEntries(userID string, rewardType string) ([]model.RewardHistoryEntry, error) {
	filter := bson.D{
		primitive.E{Key: "user_id", Value: userID},
		primitive.E{Key: "type", Value: rewardType},
	}

	var result []model.RewardHistoryEntry
	err := sa.db.rewardHistory.Find(filter, &result, &options.FindOptions{
		Sort: bson.D{{"date_created", -1}},
	})
	if err != nil {
		log.Printf("storage.GetRewardHistoryEntries error: %s", err)
		return nil, fmt.Errorf("storage.GetRewardHistoryEntries error: %s", err)
	}
	return result, nil
}

// GetRewardHistoryEntry Gets a reward history entry by id
func (sa *Adapter) GetRewardHistoryEntry(userID, id string) (*model.RewardHistoryEntry, error) {
	filter := bson.D{
		primitive.E{Key: "_id", Value: id},
		primitive.E{Key: "user_id", Value: userID},
	}
	var result []model.RewardHistoryEntry
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

// CreateRewardHistoryEntry creates a new reward history entry
func (sa *Adapter) CreateRewardHistoryEntry(item model.RewardHistoryEntry) (*model.RewardHistoryEntry, error) {
	now := time.Now().UTC()
	item.ID = uuid.NewString()
	item.DateCreated = now
	item.DateUpdated = now
	_, err := sa.db.rewardHistory.InsertOne(&item)
	if err != nil {
		log.Printf("storage.CreateRewardHistoryEntry error: %s", err)
		return nil, fmt.Errorf("storage.CreateRewardHistoryEntry error: %s", err)
	}
	return &item, nil
}

func (sa *Adapter) GetUserBalance(userID string) ([]model.WalletBalance, error){
	pipeline := []bson.M{
		{"$match": bson.M{"user_id": userID}},
		{"$group": bson.M{"_id": "$code", "amount": bson.M{"$sum":"$amount"},},},
	}

	var result []model.WalletBalance
	err := sa.db.rewardHistory.Aggregate(pipeline, &result, nil)
	if err != nil {
		log.Printf("storage.GetRewardHistoryEntries error: %s", err)
		return nil, fmt.Errorf("storage.GetRewardHistoryEntries error: %s", err)
	}
	return result, nil
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

	if "configs" == coll {
		log.Println("configs collection changed")
	} else {
		log.Println("other collection changed")
	}
}
