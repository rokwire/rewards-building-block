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
		log.Printf("storage.GetRewardType error: %s", err)
		return nil, fmt.Errorf("storage.GetRewardType error: %s", err)
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
		log.Printf("storage.CreateRewardType error: %s", err)
		return nil, fmt.Errorf("storage.CreateRewardType error: %s", err)
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
			primitive.E{Key: "description", Value: item.Description},
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

// GetRewardOperations Gets all reward operations
func (sa *Adapter) GetRewardOperations(orgID string) ([]model.RewardOperation, error) {
	filter := bson.D{
		primitive.E{Key: "org_id", Value: orgID},
	}
	var result []model.RewardOperation
	err := sa.db.rewardOperations.Find(filter, &result, nil)
	if err != nil {
		log.Printf("storage.GetRewardOperations error: %s", err)
		return nil, fmt.Errorf("storage.GetRewardOperations error: %s", err)
	}
	if result == nil {
		result = []model.RewardOperation{}
	}
	return result, nil
}

// GetRewardOperationByID Gets a reward operation by id
func (sa *Adapter) GetRewardOperationByID(orgID string, id string) (*model.RewardOperation, error) {
	filter := bson.D{
		primitive.E{Key: "org_id", Value: orgID},
		primitive.E{Key: "_id", Value: id},
	}
	var result []model.RewardOperation
	err := sa.db.rewardOperations.Find(filter, &result, nil)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result) == 0 {
		log.Printf("storage.GetRewardOperationByID error: %s", err)
		return nil, fmt.Errorf("storage.GetRewardOperationByID error: %s", err)
	}
	return &result[0], nil
}

// GetRewardOperationByCode Gets a reward type by code
func (sa *Adapter) GetRewardOperationByCode(orgID string, code string) (*model.RewardOperation, error) {
	filter := bson.D{
		primitive.E{Key: "org_id", Value: orgID},
		primitive.E{Key: "code", Value: code},
	}
	var result []model.RewardOperation
	err := sa.db.rewardOperations.Find(filter, &result, nil)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result) == 0 {
		log.Printf("storage.GetRewardOperationByCode error: unable to find reward operation with code: %s", code)
		return nil, fmt.Errorf("storage.GetRewardOperationByCode error: unable to find reward operation with code: %s", code)
	}
	return &result[0], nil
}

// CreateRewardOperation creates a new reward operation
func (sa *Adapter) CreateRewardOperation(orgID string, item model.RewardOperation) (*model.RewardOperation, error) {
	now := time.Now().UTC()
	item.ID = uuid.NewString()
	item.OrgID = orgID
	item.DateCreated = now
	item.DateUpdated = now
	_, err := sa.db.rewardOperations.InsertOne(&item)
	if err != nil {
		log.Printf("storage.CreateRewardOperation error: %s", err)
		return nil, fmt.Errorf("storage.CreateRewardOperation error: %s", err)
	}
	return &item, nil
}

// UpdateRewardOperation updates a reward operation
func (sa *Adapter) UpdateRewardOperation(orgID string, id string, item model.RewardOperation) (*model.RewardOperation, error) {
	jsonID := item.ID
	if jsonID != id {
		return nil, fmt.Errorf("storage.UpdateRewardOperation attempt to override another object")
	}

	now := time.Now().UTC()
	filter := bson.D{
		primitive.E{Key: "_id", Value: id},
		primitive.E{Key: "org_id", Value: orgID},
	}
	update := bson.D{
		primitive.E{Key: "$set", Value: bson.D{
			primitive.E{Key: "amount", Value: item.Amount},
			primitive.E{Key: "description", Value: item.Description},
			primitive.E{Key: "date_updated", Value: now},
		},
		},
	}
	_, err := sa.db.rewardOperations.UpdateOne(filter, update, nil)
	if err != nil {
		log.Printf("storage.UpdateRewardOperation error: %s", err)
		return nil, fmt.Errorf("storage.UpdateRewardOperation error: %s", err)
	}

	item.DateUpdated = now

	return &item, nil
}

// DeleteRewardOperation deletes a reward operation
func (sa *Adapter) DeleteRewardOperation(orgID string, id string) error {
	// TBD check and deny if the reward type is in use!!!

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	_, err := sa.db.rewardOperations.DeleteOne(filter, nil)
	if err != nil {
		log.Printf("storage.DeleteRewardOperation error: %s", err)
		return fmt.Errorf("storage.DeleteRewardOperation error: %s", err)
	}

	return nil
}

// GetRewardInventories Gets all reward inventories
func (sa *Adapter) GetRewardInventories(orgID string, ids []string, rewardType *string, inStock *bool, depleted *bool, limit *int64, offset *int64) ([]model.RewardInventory, error) {
	filter := bson.D{
		primitive.E{Key: "org_id", Value: orgID},
	}

	if len(ids) > 0 {
		filter = append(filter, primitive.E{Key: "_id", Value: bson.M{"$in": ids}})
	}

	if rewardType != nil {
		filter = append(filter, primitive.E{Key: "reward_type", Value: *rewardType})
	}

	if inStock != nil {
		filter = append(filter, primitive.E{Key: "in_stock", Value: *inStock})
	}

	if depleted != nil {
		filter = append(filter, primitive.E{Key: "depleted", Value: *depleted})
	}

	findOptions := options.FindOptions{
		Sort: bson.D{{"date_created", -1}},
	}
	if limit != nil {
		findOptions.SetLimit(*limit)
	}
	if offset != nil {
		findOptions.SetSkip(*offset)
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

// GetUserRewardsHistory Gets all reward history entries
func (sa *Adapter) GetUserRewardsHistory(orgID string, userID string, rewardType *string, code *string, buildingBlock *string, limit *int64, offset *int64) ([]model.Reward, error) {
	filter := bson.D{
		primitive.E{Key: "org_id", Value: orgID},
		primitive.E{Key: "user_id", Value: userID},
	}

	if rewardType != nil {
		filter = append(filter, primitive.E{Key: "reward_type", Value: *rewardType})
	}

	if code != nil {
		filter = append(filter, primitive.E{Key: "code", Value: *code})
	}

	if rewardType != nil {
		filter = append(filter, primitive.E{Key: "building_block", Value: *buildingBlock})
	}

	findOptions := options.FindOptions{
		Sort: bson.D{{"date_created", -1}},
	}
	if limit != nil {
		findOptions.SetLimit(*limit)
	}
	if offset != nil {
		findOptions.SetSkip(*offset)
	}

	var result []model.Reward
	err := sa.db.rewardHistory.Find(filter, &result, &findOptions)
	if err != nil {
		log.Printf("storage.getUserRewardsHistory error: %s", err)
		return nil, fmt.Errorf("storage.getUserRewardsHistory error: %s", err)
	}
	if result == nil {
		result = []model.Reward{}
	}
	return result, nil
}

// GetUserRewardByID Gets a reward history entry by id
func (sa *Adapter) GetUserRewardByID(orgID string, userID, id string) (*model.Reward, error) {
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
		log.Printf("storage.GetUserRewardByID error: %s", err)
		return nil, fmt.Errorf("storage.GetUserRewardByID error: %s", err)
	}
	return &result[0], nil
}

// CreateUserReward creates a new reward history entry
func (sa *Adapter) CreateUserReward(orgID string, item model.Reward) (*model.Reward, error) {
	now := time.Now().UTC()
	item.ID = uuid.NewString()
	item.DateCreated = now
	item.DateUpdated = now
	item.OrgID = orgID
	_, err := sa.db.rewardHistory.InsertOne(&item)
	if err != nil {
		log.Printf("storage.CreateUserReward error: %s", err)
		return nil, fmt.Errorf("storage.CreateUserReward error: %s", err)
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
		log.Printf("storage.GetUserBalance error: %s", err)
		return nil, fmt.Errorf("storage.GetUserBalance error: %s", err)
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
		log.Printf("storage.GetWalletBalance error: %s", err)
		return nil, fmt.Errorf("storage.GetWalletBalance error: %s", err)
	}

	if len(result) > 0 {
		return &result[0], nil
	}

	return nil, nil
}

// GetRewardQuantity Gets reward quantities state for the current moment
func (sa *Adapter) GetRewardQuantity(orgID string, rewardType string) (*model.RewardQuantity, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"org_id": orgID, "reward_type": rewardType}},
		{"$group": bson.M{"_id": "$code", "amount": bson.M{"$sum": "$amount"}}},
	}

	var inventoryAmount int64
	var inventoryResult []struct {
		Amount int64 `bson:"amount"`
	}
	err := sa.db.rewardInventories.Aggregate(pipeline, &inventoryResult, nil)
	if err != nil {
		log.Printf("storage.GetRewardQuantity error: %s", err)
		return nil, fmt.Errorf("storage.GetRewardQuantity error: %s", err)
	}
	if len(inventoryResult) > 0 {
		inventoryAmount = inventoryResult[0].Amount
	}

	pipeline = []bson.M{
		{"$match": bson.M{"org_id": orgID, "reward_type": rewardType}},
		{"$group": bson.M{"_id": "$code", "amount": bson.M{"$sum": "$amount"}}},
	}

	var inventoryInStockAmount int64
	var inventoryInStockResult []struct {
		Amount int64 `bson:"amount"`
	}
	err = sa.db.rewardInventories.Aggregate(pipeline, &inventoryInStockResult, nil)
	if err != nil {
		log.Printf("storage.GetRewardQuantity error: %s", err)
		return nil, fmt.Errorf("storage.GetRewardQuantity error: %s", err)
	}
	if len(inventoryInStockResult) > 0 {
		inventoryInStockAmount = inventoryInStockResult[0].Amount
	}

	pipeline = []bson.M{
		{"$match": bson.M{"org_id": orgID, "reward_type": rewardType}},
		{"$group": bson.M{"_id": "$code", "amount": bson.M{"$sum": "$amount"}}},
	}

	var rewardsAmount int64
	var rewardsResult []struct {
		Amount int64 `bson:"amount"`
	}
	err = sa.db.rewardHistory.Aggregate(pipeline, &rewardsResult, nil)
	if err != nil {
		log.Printf("storage.GetRewardQuantity error: %s", err)
		return nil, fmt.Errorf("storage.GetRewardQuantity error: %s", err)
	}
	if len(rewardsResult) > 0 {
		rewardsAmount = rewardsResult[0].Amount
	}

	pipeline = []bson.M{
		{"$unwind": bson.M{"path": "$items"}},
		{"$match": bson.M{"org_id": orgID, "items.reward_type": rewardType}},
		{"$group": bson.M{"_id": rewardType, "amount": bson.M{"$sum": "$items.amount"}}},
	}

	var claimsAmount int64
	var claimsResult []struct {
		Amount int64 `bson:"amount"`
	}
	err = sa.db.rewardClaims.Aggregate(pipeline, &claimsResult, nil)
	if err != nil {
		log.Printf("storage.GetRewardQuantity error: %s", err)
		return nil, fmt.Errorf("storage.GetRewardQuantity error: %s", err)
	}
	if len(claimsResult) > 0 {
		claimsAmount = claimsResult[0].Amount
	}

	rewardableQuantity := inventoryAmount - rewardsAmount
	claimableQuantity := inventoryInStockAmount - claimsAmount
	return &model.RewardQuantity{
		RewardType:         rewardType,
		RewardableQuantity: rewardableQuantity,
		ClaimableQuantity:  claimableQuantity,
	}, nil
}

// GetRewardClaims Gets all reward claims
func (sa *Adapter) GetRewardClaims(orgID string, ids []string, userID *string, rewardType *string, status *string, limit *int64, offset *int64) ([]model.RewardClaim, error) {
	filter := bson.D{
		primitive.E{Key: "org_id", Value: orgID},
	}
	var result []model.RewardClaim
	err := sa.db.rewardClaims.Find(filter, &result, nil)
	if err != nil {
		log.Printf("storage.getRewardClaims error: %s", err)
		return nil, fmt.Errorf("storage.getRewardClaims error: %s", err)
	}
	if result == nil {
		result = []model.RewardClaim{}
	}
	return result, nil
}

// GetRewardClaim Gets a reward claim by id
func (sa *Adapter) GetRewardClaim(orgID string, id string) (*model.RewardClaim, error) {
	filter := bson.D{
		primitive.E{Key: "org_id", Value: orgID},
		primitive.E{Key: "_id", Value: id},
	}
	var result []model.RewardClaim
	err := sa.db.rewardClaims.Find(filter, &result, nil)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result) == 0 {
		log.Printf("storage.getRewardClaim error: %s", err)
		return nil, fmt.Errorf("storage.getRewardClaim error: %s", err)
	}
	return &result[0], nil
}

// CreateRewardClaim creates a new reward claim
func (sa *Adapter) CreateRewardClaim(orgID string, item model.RewardClaim) (*model.RewardClaim, error) {
	now := time.Now().UTC()
	item.ID = uuid.NewString()
	item.OrgID = orgID
	item.DateCreated = now
	item.DateUpdated = now
	_, err := sa.db.rewardClaims.InsertOne(&item)
	if err != nil {
		log.Printf("storage.createRewardClaim error: %s", err)
		return nil, fmt.Errorf("storage.createRewardClaim error: %s", err)
	}
	return &item, nil
}

// UpdateRewardClaim updates a reward claim
func (sa *Adapter) UpdateRewardClaim(orgID string, id string, item model.RewardClaim) (*model.RewardClaim, error) {
	jsonID := item.ID
	if jsonID != id {
		return nil, fmt.Errorf("storage.updateRewardClaim attempt to override another object")
	}

	now := time.Now().UTC()
	filter := bson.D{
		primitive.E{Key: "_id", Value: id},
		primitive.E{Key: "org_id", Value: orgID},
	}
	update := bson.D{
		primitive.E{Key: "$set", Value: bson.D{
			primitive.E{Key: "items", Value: item.Items},
			primitive.E{Key: "description", Value: item.Description},
			primitive.E{Key: "status", Value: item.Status},
			primitive.E{Key: "date_updated", Value: now},
		},
		},
	}
	_, err := sa.db.rewardClaims.UpdateOne(filter, update, nil)
	if err != nil {
		log.Printf("storage.updateRewardClaim error: %s", err)
		return nil, fmt.Errorf("storage.updateRewardClaim error: %s", err)
	}

	item.DateUpdated = now

	return &item, nil
}

// DeleteRewardClaim deletes a reward claim
func (sa *Adapter) DeleteRewardClaim(orgID string, id string) error {
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	_, err := sa.db.rewardClaims.DeleteOne(filter, nil)
	if err != nil {
		log.Printf("storage.deleteRewardClaim error: %s", err)
		return fmt.Errorf("storage.deleteRewardClaim error: %s", err)
	}

	return nil
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
