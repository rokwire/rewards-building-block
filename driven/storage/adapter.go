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
	"context"
	"fmt"
	"log"
	"rewards/core/model"
	"strconv"
	"time"

	"github.com/rokwire/logging-library-go/logs"

	"github.com/rokwire/logging-library-go/errors"
	"github.com/rokwire/logging-library-go/logutils"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Adapter implements the Storage interface
type Adapter struct {
	db     *database
	logger *logs.Logger
}

// Start starts the storage
func (sa *Adapter) Start() error {

	err := sa.db.start()
	return err
}

//PerformTransaction performs a transaction
func (sa *Adapter) PerformTransaction(transaction func(context TransactionContext) error) error {
	// transaction
	err := sa.db.dbClient.UseSession(context.Background(), func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			sa.abortTransaction(sessionContext)
			return errors.WrapErrorAction(logutils.ActionStart, logutils.TypeTransaction, nil, err)
		}
		err = transaction(sessionContext)
		if err != nil {
			sa.abortTransaction(sessionContext)
			return errors.WrapErrorAction("performing", logutils.TypeTransaction, nil, err)
		}

		err = sessionContext.CommitTransaction(sessionContext)
		if err != nil {
			sa.abortTransaction(sessionContext)
			return errors.WrapErrorAction(logutils.ActionCommit, logutils.TypeTransaction, nil, err)
		}
		return nil
	})

	return err
}

func (sa *Adapter) abortTransaction(sessionContext mongo.SessionContext) {
	err := sessionContext.AbortTransaction(sessionContext)
	if err != nil {
		sa.logger.Errorf("error aborting a transaction - %s", err)
	}
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

//StoreMultiTenancyData stores multi-tenancy to already exisiting data in the collections
func (sa *Adapter) StoreMultiTenancyData(context TransactionContext, appID string, orgID string) error {
	//TODO

	filter := bson.D{}
	update := bson.D{
		primitive.E{Key: "$set", Value: bson.D{
			primitive.E{Key: "app_id", Value: appID},
			primitive.E{Key: "org_id", Value: orgID},
		}},
	}
	//types
	_, err := sa.db.rewardTypes.UpdateManyWithContext(context, filter, update, nil)
	if err != nil {
		return err
	}
	//history
	_, err = sa.db.rewardHistory.UpdateManyWithContext(context, filter, update, nil)
	if err != nil {
		return err
	}
	//operations
	_, err = sa.db.rewardOperations.UpdateManyWithContext(context, filter, update, nil)
	if err != nil {
		return err
	}
	//inventories
	_, err = sa.db.rewardInventories.UpdateManyWithContext(context, filter, update, nil)
	if err != nil {
		return err
	}

	//claims
	_, err = sa.db.rewardClaims.UpdateManyWithContext(context, filter, update, nil)
	if err != nil {
		return err
	}

	return nil
}

// FindAllContentItems  finds all content items
func (sa *Adapter) FindAllRewardTypeItems(context TransactionContext) ([]model.RewardType, error) {
	filter := bson.D{}
	var result []model.RewardType
	err := sa.db.rewardTypes.Find(filter, &result, nil)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetRewardTypes Gets all reward types
func (sa *Adapter) GetRewardTypes(appID *string, orgID string) ([]model.RewardType, error) {
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
func (sa *Adapter) GetRewardType(appID *string, orgID string, id string) (*model.RewardType, error) {
	filter := bson.D{
		primitive.E{Key: "app_id", Value: appID},
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
func (sa *Adapter) GetRewardTypeByType(appID *string, orgID string, rewardType string) (*model.RewardType, error) {
	filter := bson.D{
		primitive.E{Key: "app_id", Value: appID},
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
func (sa *Adapter) CreateRewardType(appID *string, orgID string, item model.RewardType) (*model.RewardType, error) {
	now := time.Now().UTC()
	item.ID = uuid.NewString()
	item.OrgID = orgID
	item.AppID = *appID
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
func (sa *Adapter) UpdateRewardType(appID *string, orgID string, id string, item model.RewardType) (*model.RewardType, error) {
	jsonID := item.ID
	if jsonID != id {
		return nil, fmt.Errorf("storage.UpdateRewardType attempt to override another object")
	}

	now := time.Now().UTC()
	filter := bson.D{
		primitive.E{Key: "_id", Value: id},
		primitive.E{Key: "org_id", Value: orgID},
		primitive.E{Key: "app_id", Value: appID},
	}
	update := bson.D{
		primitive.E{Key: "$set", Value: bson.D{
			primitive.E{Key: "display_name", Value: item.DisplayName},
			primitive.E{Key: "active", Value: item.Active},
			primitive.E{Key: "description", Value: item.Description},
			primitive.E{Key: "date_updated", Value: now},
		}},
	}
	_, err := sa.db.rewardTypes.UpdateOne(filter, update, nil)
	if err != nil {
		log.Printf("storage.UpdateRewardType error: %s", err)
		return nil, fmt.Errorf("storage.UpdateRewardType error: %s", err)
	}

	item.DateUpdated = now

	return &item, nil
}

// DeleteRewardType deletes a reward type
func (sa *Adapter) DeleteRewardType(appID *string, orgID string, id string) error {
	// TBD check and deny if the reward type is in use!!!

	filter := bson.D{primitive.E{Key: "_id", Value: id}, primitive.E{Key: "org_id", Value: orgID},
		primitive.E{Key: "app_id", Value: appID}}
	_, err := sa.db.rewardTypes.DeleteOne(filter, nil)
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
		primitive.E{Key: "org_id", Value: orgID},
		primitive.E{Key: "_id", Value: id},
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

	filter := bson.D{
		primitive.E{Key: "org_id", Value: orgID},
		primitive.E{Key: "_id", Value: id},
	}
	_, err := sa.db.rewardOperations.DeleteOne(filter, nil)
	if err != nil {
		log.Printf("storage.DeleteRewardOperation error: %s", err)
		return fmt.Errorf("storage.DeleteRewardOperation error: %s", err)
	}

	return nil
}

// GetRewardInventories Gets all reward inventories
func (sa *Adapter) GetRewardInventories(orgID string, ids []string, rewardType *string, inStock *bool, grantDepleted *bool, claimDepleted *bool, limit *int64, offset *int64) ([]model.RewardInventory, error) {
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

	if grantDepleted != nil {
		filter = append(filter, primitive.E{Key: "grant_depleted", Value: *grantDepleted})
	}

	if claimDepleted != nil {
		filter = append(filter, primitive.E{Key: "claim_depleted", Value: *claimDepleted})
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
	filter := bson.D{
		primitive.E{Key: "org_id", Value: orgID},
		primitive.E{Key: "_id", Value: id},
	}
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
	item.OrgID = orgID

	if err := sa.validateInventoryCreateOrUpdate(item); err != nil {
		return nil, err
	}

	_, err := sa.db.rewardInventories.InsertOne(&item)
	if err != nil {
		log.Printf("storage.CreateRewardInventory error: %s", err)
		return nil, fmt.Errorf("storage.CreateRewardInventory error: %s", err)
	}
	return &item, nil
}

// UpdateRewardInventory updates a reward pool
func (sa *Adapter) UpdateRewardInventory(orgID string, id string, item model.RewardInventory) (*model.RewardInventory, error) {
	return sa.UpdateRewardInventoryWithContext(nil, orgID, id, item)
}

// UpdateRewardInventoryWithContext updates a reward inventory with a context
func (sa *Adapter) UpdateRewardInventoryWithContext(ctx context.Context, orgID string, id string, item model.RewardInventory) (*model.RewardInventory, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	jsonID := item.ID
	if jsonID != id || orgID != item.OrgID {
		return nil, fmt.Errorf("storage.UpdateRewardInventory attempt to override another object")
	}

	if err := sa.validateInventoryCreateOrUpdate(item); err != nil {
		return nil, err
	}

	item.GrantDepleted = item.AmountTotal <= item.AmountGranted
	item.ClaimDepleted = item.AmountTotal <= item.AmountClaimed

	now := time.Now().UTC()
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{
		primitive.E{Key: "$set", Value: bson.D{
			primitive.E{Key: "date_updated", Value: now},
			primitive.E{Key: "amount_total", Value: item.AmountTotal},
			primitive.E{Key: "amount_granted", Value: item.AmountGranted},
			primitive.E{Key: "amount_claimed", Value: item.AmountClaimed},
			primitive.E{Key: "grant_depleted", Value: item.GrantDepleted},
			primitive.E{Key: "claim_depleted", Value: item.ClaimDepleted},
			primitive.E{Key: "in_stock", Value: item.InStock},
			primitive.E{Key: "description", Value: item.Description},
		},
		},
	}
	_, err := sa.db.rewardInventories.UpdateOneWithContext(ctx, filter, update, nil)
	if err != nil {
		log.Printf("storage.UpdateRewardInventory error: %s", err)
		return nil, fmt.Errorf("storage.UpdateRewardInventory error: %s", err)
	}

	item.DateUpdated = now

	return &item, nil
}

func (sa *Adapter) validateInventoryCreateOrUpdate(item model.RewardInventory) error {
	if item.AmountTotal <= 0 {
		return fmt.Errorf("inventory amount is zero or negative")
	} else if item.AmountGranted > item.AmountTotal {
		return fmt.Errorf("quantity granted is greater than the total amount")
	} else if item.AmountClaimed > item.AmountTotal {
		return fmt.Errorf("quantity claimed is greater than the total amount")
	}

	return nil
}

// DeleteRewardInventory deletes a reward pool. Don't delete if it's in use!
func (sa *Adapter) DeleteRewardInventory(orgID string, id string) error {
	// TBD check and deny if the reward type is in use!!!

	filter := bson.D{
		primitive.E{Key: "org_id", Value: orgID},
		primitive.E{Key: "_id", Value: id},
	}
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

	err := sa.db.dbClient.UseSession(context.Background(), func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			log.Printf("error starting a transaction - %s", err)
			return err
		}

		grantDepleted := false
		inventories, err := sa.GetRewardInventories(orgID, nil, &item.RewardType, nil, &grantDepleted, nil, nil, nil)
		if err != nil {
			log.Printf("storage.CreateUserReward error: %s", err)
			return fmt.Errorf("storage.CreateUserReward error: %s", err)
		}

		if len(inventories) > 0 {
			accumulatedAmount := 0
			remaingAmnount := item.Amount
			for _, inventory := range inventories {
				grantableAmount := inventory.GetGrantableAmount()
				if grantableAmount > 0 {
					if grantableAmount < remaingAmnount {
						inventory.AmountGranted += grantableAmount
						accumulatedAmount += grantableAmount
						remaingAmnount -= grantableAmount
					} else {
						inventory.AmountGranted += remaingAmnount
						accumulatedAmount += remaingAmnount
						remaingAmnount -= remaingAmnount
					}
					_, err = sa.UpdateRewardInventoryWithContext(sessionContext, orgID, inventory.ID, inventory)
					if err != nil {
						abortTransaction(sessionContext)
						log.Printf("storage.CreateUserReward error: %s", err)
						return fmt.Errorf("storage.CreateUserReward error: %s", err)
					}
					if accumulatedAmount == item.Amount {
						break
					}
				}
			}

			if accumulatedAmount < item.Amount {
				abortTransaction(sessionContext)
				log.Printf("storage.CreateUserReward insuficient amount in the inventory for: %s", item.RewardType)
				return fmt.Errorf("storage.CreateUserReward insuficient amount in the inventory for: %s", item.RewardType)
			}
		}

		_, err = sa.db.rewardHistory.InsertOneWithContext(sessionContext, &item)
		if err != nil {
			abortTransaction(sessionContext)
			log.Printf("storage.CreateUserReward error: %s", err)
			return fmt.Errorf("storage.CreateUserReward error: %s", err)
		}

		//commit the transaction
		err = sessionContext.CommitTransaction(sessionContext)
		if err != nil {
			abortTransaction(sessionContext)
			fmt.Println(err)
			return err
		}
		return nil
	})

	if err != nil {
		log.Printf("storage.CreateUserReward transaction error: %s", err)
		return nil, fmt.Errorf("storage.CreateUserReward transaction error: %s", err)
	}

	return &item, nil
}

// GetUserRewardsAmount Gets user's rewards amount
func (sa *Adapter) GetUserRewardsAmount(orgID string, userID string, rewardType *string) ([]model.RewardTypeAmount, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"org_id": orgID, "user_id": userID}},
	}
	if rewardType != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"reward_type": *rewardType}})
	}
	pipeline = append(pipeline, bson.M{"$group": bson.M{"_id": "$reward_type", "amount": bson.M{"$sum": "$amount"}}})

	var result []model.RewardTypeAmount
	err := sa.db.rewardHistory.Aggregate(pipeline, &result, nil)
	if err != nil {
		log.Printf("storage.GetUserRewardsAmount error: %s", err)
		return nil, fmt.Errorf("storage.GetUserRewardsAmount error: %s", err)
	}

	return result, nil
}

// GetUserClaimsAmount Gets user's claims amount
func (sa *Adapter) GetUserClaimsAmount(orgID string, userID string, rewardType *string) ([]model.RewardTypeAmount, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"org_id": orgID, "user_id": userID}},
		{"$unwind": bson.M{"path": "$items"}},
	}

	if rewardType != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"items.reward_type": *rewardType}})
	}
	pipeline = append(pipeline, bson.M{"$group": bson.M{"_id": "$items.reward_type", "amount": bson.M{"$sum": "$items.amount"}}})

	var result []model.RewardTypeAmount
	err := sa.db.rewardClaims.Aggregate(pipeline, &result, nil)
	if err != nil {
		log.Printf("storage.GetUserClaimsAmount error: %s", err)
		return nil, fmt.Errorf("storage.GetUserClaimsAmount error: %s", err)
	}

	return result, nil
}

// GetRewardQuantityState Gets reward quantities state for the current moment
func (sa *Adapter) GetRewardQuantityState(orgID string, rewardType string, inStock *bool) (*model.RewardQuantityState, error) {

	inventories, err := sa.GetRewardInventories(orgID, nil, &rewardType, inStock, nil, nil, nil, nil)
	if err != nil {
		log.Printf("storage.GetRewardQuantityState error: %s", err)
		return nil, fmt.Errorf("storage.GetRewardQuantityState error: %s", err)
	}

	var totalQuantity int = 0
	var grantedQuantity int = 0
	var grantableQuantity int = 0
	var claimableQuantity int = 0
	if len(inventories) > 0 {
		for _, inventory := range inventories {
			totalQuantity += inventory.AmountTotal
			grantedQuantity += inventory.AmountGranted
			if inventory.InStock {
				claimableQuantity += inventory.AmountTotal - inventory.AmountClaimed
			}
		}
		grantableQuantity = totalQuantity - grantedQuantity

		return &model.RewardQuantityState{
			RewardType:        rewardType,
			GrantableQuantity: grantableQuantity,
			ClaimableQuantity: claimableQuantity,
		}, nil
	}

	return nil, nil
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

	err := sa.db.dbClient.UseSession(context.Background(), func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			log.Printf("error starting a transaction - %s", err)
			return err
		}

		for _, claimEntry := range item.Items {
			claimDepleted := false
			inventories, err := sa.GetRewardInventories(orgID, nil, &claimEntry.RewardType, nil, nil, &claimDepleted, nil, nil)
			if err != nil {
				log.Printf("storage.CreateUserReward error: %s", err)
				return fmt.Errorf("storage.CreateRewardClaim error: %s", err)
			}

			if len(inventories) > 0 {
				accumulatedAmount := 0
				remaingAmnount := claimEntry.Amount
				for _, inventory := range inventories {
					claimableAmount := inventory.GetClaimableAmount()
					if claimableAmount > 0 {
						if claimableAmount < remaingAmnount {
							inventory.AmountClaimed += claimableAmount
							accumulatedAmount += claimableAmount
							remaingAmnount -= claimableAmount
						} else {
							inventory.AmountClaimed += remaingAmnount
							accumulatedAmount += remaingAmnount
							remaingAmnount -= remaingAmnount
						}
						_, err = sa.UpdateRewardInventoryWithContext(sessionContext, orgID, inventory.ID, inventory)
						if err != nil {
							abortTransaction(sessionContext)
							log.Printf("storage.CreateUserReward error: %s", err)
							return fmt.Errorf("storage.CreateUserReward error: %s", err)
						}
						if accumulatedAmount == claimEntry.Amount {
							break
						}
					}
				}

				if accumulatedAmount < claimEntry.Amount {
					abortTransaction(sessionContext)
					log.Printf("storage.CreateRewardClaim insuficient amount in the inventory for: %s", claimEntry.RewardType)
					return fmt.Errorf("storage.CreateRewardClaim insuficient amount in the inventory for: %s", claimEntry.RewardType)
				}
			}
		}

		_, err = sa.db.rewardClaims.InsertOneWithContext(sessionContext, &item)
		if err != nil {
			log.Printf("storage.CreateRewardClaim error: %s", err)
			return fmt.Errorf("storage.CreateRewardClaim error: %s", err)
		}

		//commit the transaction
		err = sessionContext.CommitTransaction(sessionContext)
		if err != nil {
			abortTransaction(sessionContext)
			fmt.Println(err)
			return err
		}

		return nil
	})

	if err != nil {
		log.Printf("storage.CreateRewardClaim transaction error: %s", err)
		return nil, fmt.Errorf("storage.CreateRewardClaim transaction error: %s", err)
	}

	return &item, nil
}

// UpdateRewardClaim updates a reward claim
func (sa *Adapter) UpdateRewardClaim(orgID string, id string, item model.RewardClaim) (*model.RewardClaim, error) {
	return sa.UpdateRewardClaimWithContext(nil, orgID, id, item)
}

// UpdateRewardClaimWithContext updates a reward claim with a context
func (sa *Adapter) UpdateRewardClaimWithContext(ctx context.Context, orgID string, id string, item model.RewardClaim) (*model.RewardClaim, error) {
	if ctx == nil {
		ctx = context.Background()
	}

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
			primitive.E{Key: "description", Value: item.Description},
			primitive.E{Key: "status", Value: item.Status},
			primitive.E{Key: "date_updated", Value: now},
		}},
	}
	_, err := sa.db.rewardClaims.UpdateOneWithContext(ctx, filter, update, nil)
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

func abortTransaction(sessionContext mongo.SessionContext) {
	err := sessionContext.AbortTransaction(sessionContext)
	if err != nil {
		log.Printf("error on aborting a transaction - %s", err)
	}
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

//TransactionContext wraps mongo.SessionContext for use by external packages
type TransactionContext interface {
	mongo.SessionContext
}
