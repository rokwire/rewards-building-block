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
	item.ID = uuid.NewString()
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

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	err := sa.db.rewardTypes.ReplaceOne(filter, item, nil)
	if err != nil {
		log.Printf("storage.UpdateRewardType error: %s", err)
		return nil, fmt.Errorf("storage.UpdateRewardType error: %s", err)
	}
	return &item, nil
}

// DeleteRewardTypes deletes a reward type
func (sa *Adapter) DeleteRewardTypes(id string) error {
	// TBD check and deny if the reward type is in use!!!

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	_, err := sa.db.rewardTypes.DeleteOne(filter, nil)
	if err != nil {
		log.Printf("storage.DeleteRewardTypes error: %s", err)
		return fmt.Errorf("storage.DeleteRewardTypes error: %s", err)
	}

	return nil
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
