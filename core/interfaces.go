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

package core

import (
	"rewards/core/model"
	"rewards/driven/storage"
)

// Services exposes APIs for the driver adapters
type Services interface {
	GetVersion() string

	GetRewardTypes(allApps bool, appID *string, orgID string) ([]model.RewardType, error)
	GetRewardType(allApps bool, appID *string, orgID string, id string) (*model.RewardType, error)
	CreateRewardType(appID *string, orgID string, item model.RewardType) (*model.RewardType, error)
	UpdateRewardType(appID *string, orgID string, id string, item model.RewardType) (*model.RewardType, error)
	DeleteRewardType(allApps bool, appID *string, orgID string, id string) error

	GetRewardOperations(allApps bool, appID *string, orgID string) ([]model.RewardOperation, error)
	GetRewardOperationByID(allApps bool, appID *string, orgID string, id string) (*model.RewardOperation, error)
	GetRewardOperationByCode(appID *string, orgID string, code string) (*model.RewardOperation, error)
	CreateRewardOperation(appID *string, orgID string, item model.RewardOperation) (*model.RewardOperation, error)
	UpdateRewardOperation(appID *string, orgID string, id string, item model.RewardOperation) (*model.RewardOperation, error)
	DeleteRewardOperation(allApps bool, appID *string, orgID string, id string) error

	GetRewardInventories(allApps bool, appID *string, orgID string, ids []string, rewardType *string, inStock *bool, grantDepleted *bool, claimDepleted *bool, limit *int64, offset *int64) ([]model.RewardInventory, error)
	GetRewardInventory(allApps bool, appID *string, orgID string, id string) (*model.RewardInventory, error)
	CreateRewardInventory(appID *string, orgID string, item model.RewardInventory) (*model.RewardInventory, error)
	UpdateRewardInventory(appID *string, orgID string, id string, item model.RewardInventory) (*model.RewardInventory, error)

	GetRewardClaims(allApps bool, appID *string, orgID string, ids []string, userID *string, rewardType *string, status *string, limit *int64, offset *int64) ([]model.RewardClaim, error)
	GetRewardClaim(allApps bool, appID *string, orgID string, id string) (*model.RewardClaim, error)
	CreateRewardClaim(appID *string, orgID string, item model.RewardClaim) (*model.RewardClaim, error)
	UpdateRewardClaim(appID *string, orgID string, id string, item model.RewardClaim) (*model.RewardClaim, error)

	CreateReward(appID *string, orgID string, item model.Reward) (*model.Reward, error)

	GetUserBalance(appID *string, orgID string, userID string) ([]model.RewardTypeAmount, error)
	GetUserRewardsHistory(appID *string, orgID string, userID string, rewardType *string, code *string, buildingBlock *string, limit *int64, offset *int64) ([]model.Reward, error)

	GetRewardQuantity(appID *string, orgID string, rewardType string) (*model.RewardQuantityState, error)
}

type servicesImpl struct {
	app *Application
}

func (s *servicesImpl) GetVersion() string {
	return s.app.getVersion()
}

func (s *servicesImpl) GetRewardTypes(allApps bool, appID *string, orgID string) ([]model.RewardType, error) {
	return s.app.getRewardTypes(allApps, appID, orgID)
}

func (s *servicesImpl) GetRewardType(allApps bool, appID *string, orgID string, id string) (*model.RewardType, error) {
	return s.app.getRewardType(allApps, appID, orgID, id)
}

func (s *servicesImpl) CreateRewardType(appID *string, orgID string, item model.RewardType) (*model.RewardType, error) {
	return s.app.createRewardType(appID, orgID, item)
}

func (s *servicesImpl) UpdateRewardType(appID *string, orgID string, id string, item model.RewardType) (*model.RewardType, error) {
	return s.app.updateRewardType(appID, orgID, id, item)
}

func (s *servicesImpl) DeleteRewardType(allApps bool, appID *string, orgID string, id string) error {
	return s.app.deleteRewardTypes(allApps, appID, orgID, id)
}

func (s *servicesImpl) GetRewardOperations(allApps bool, appID *string, orgID string) ([]model.RewardOperation, error) {
	return s.app.getRewardOperations(allApps, appID, orgID)
}

func (s *servicesImpl) GetRewardOperationByID(allApps bool, appID *string, orgID string, id string) (*model.RewardOperation, error) {
	return s.app.getRewardOperationByID(allApps, appID, orgID, id)
}

func (s *servicesImpl) GetRewardOperationByCode(appID *string, orgID string, code string) (*model.RewardOperation, error) {
	return s.app.getRewardOperationByCode(appID, orgID, code)
}

func (s *servicesImpl) CreateRewardOperation(appID *string, orgID string, item model.RewardOperation) (*model.RewardOperation, error) {
	return s.app.createRewardOperation(appID, orgID, item)
}

func (s *servicesImpl) UpdateRewardOperation(appID *string, orgID string, id string, item model.RewardOperation) (*model.RewardOperation, error) {
	return s.app.updateRewardOperation(appID, orgID, id, item)
}

func (s *servicesImpl) DeleteRewardOperation(allApps bool, appID *string, orgID string, id string) error {
	return s.app.deleteRewardOperation(allApps, appID, orgID, id)
}

func (s *servicesImpl) GetRewardInventories(allApps bool, appID *string, orgID string, ids []string, rewardType *string, inStock *bool, grantDepleted *bool, claimDepleted *bool, limit *int64, offset *int64) ([]model.RewardInventory, error) {
	return s.app.getRewardInventories(allApps, appID, orgID, ids, rewardType, inStock, grantDepleted, claimDepleted, limit, offset)
}

func (s *servicesImpl) GetRewardInventory(allApps bool, appID *string, orgID string, id string) (*model.RewardInventory, error) {
	return s.app.getRewardInventory(allApps, appID, orgID, id)
}

func (s *servicesImpl) CreateRewardInventory(appID *string, orgID string, item model.RewardInventory) (*model.RewardInventory, error) {
	return s.app.createRewardInventory(appID, orgID, item)
}

func (s *servicesImpl) UpdateRewardInventory(appID *string, orgID string, id string, item model.RewardInventory) (*model.RewardInventory, error) {
	return s.app.updateRewardInventory(appID, orgID, id, item)
}

func (s *servicesImpl) DeleteRewardInventory(allApps bool, appID *string, orgID string, id string) error {
	return s.app.deleteRewardTypes(allApps, appID, orgID, id)
}

func (s *servicesImpl) CreateReward(appID *string, orgID string, item model.Reward) (*model.Reward, error) {
	return s.app.createReward(appID, orgID, item)
}

func (s *servicesImpl) GetRewardClaims(allApps bool, appID *string, orgID string, ids []string, userID *string, rewardType *string, status *string, limit *int64, offset *int64) ([]model.RewardClaim, error) {
	return s.app.getRewardClaims(allApps, appID, orgID, ids, userID, rewardType, status, limit, offset)
}

func (s *servicesImpl) GetRewardClaim(allApps bool, appID *string, orgID string, id string) (*model.RewardClaim, error) {
	return s.app.getRewardClaim(allApps, appID, orgID, id)
}

func (s *servicesImpl) CreateRewardClaim(appID *string, orgID string, item model.RewardClaim) (*model.RewardClaim, error) {
	return s.app.createRewardClaim(appID, orgID, item)
}

func (s *servicesImpl) UpdateRewardClaim(appID *string, orgID string, id string, item model.RewardClaim) (*model.RewardClaim, error) {
	return s.app.updateRewardClaim(appID, orgID, id, item)
}

func (s *servicesImpl) GetUserBalance(appID *string, orgID string, userID string) ([]model.RewardTypeAmount, error) {
	return s.app.getUserBalance(appID, orgID, userID)
}

func (s *servicesImpl) GetUserRewardsHistory(appID *string, orgID string, userID string, rewardType *string, code *string, buildingBlock *string, limit *int64, offset *int64) ([]model.Reward, error) {
	return s.app.getUserRewardsHistory(appID, orgID, userID, rewardType, code, buildingBlock, limit, offset)
}

func (s *servicesImpl) GetRewardQuantity(appID *string, orgID string, rewardType string) (*model.RewardQuantityState, error) {
	return s.app.getRewardQuantity(appID, orgID, rewardType)
}

// Storage is used by core to storage data - DB storage adapter, file storage adapter etc
type Storage interface {
	PerformTransaction(func(context storage.TransactionContext) error) error

	GetRewardTypes(appID *string, orgID string) ([]model.RewardType, error)
	GetRewardType(appID *string, orgID string, id string) (*model.RewardType, error)
	GetRewardTypeByType(appID *string, orgID string, rewardType string) (*model.RewardType, error)
	CreateRewardType(appID *string, orgID string, item model.RewardType) (*model.RewardType, error)
	UpdateRewardType(appID *string, orgID string, id string, item model.RewardType) (*model.RewardType, error)
	DeleteRewardType(appID *string, orgID string, id string) error

	GetRewardOperations(appID *string, orgID string) ([]model.RewardOperation, error)
	GetRewardOperationByID(appID *string, orgID string, id string) (*model.RewardOperation, error)
	GetRewardOperationByCode(appID *string, orgID string, code string) (*model.RewardOperation, error)
	CreateRewardOperation(appID *string, orgID string, item model.RewardOperation) (*model.RewardOperation, error)
	UpdateRewardOperation(appID *string, orgID string, id string, item model.RewardOperation) (*model.RewardOperation, error)
	DeleteRewardOperation(appID *string, orgID string, id string) error

	GetRewardInventories(appID *string, orgID string, ids []string, rewardType *string, inStock *bool, grantDepleted *bool, claimDepleted *bool, limit *int64, offset *int64) ([]model.RewardInventory, error)
	GetRewardInventory(appID *string, orgID string, id string) (*model.RewardInventory, error)
	CreateRewardInventory(appID *string, orgID string, item model.RewardInventory) (*model.RewardInventory, error)
	UpdateRewardInventory(appID *string, orgID string, id string, item model.RewardInventory) (*model.RewardInventory, error)

	GetRewardClaims(appID *string, orgID string, ids []string, userID *string, rewardType *string, status *string, limit *int64, offset *int64) ([]model.RewardClaim, error)
	GetRewardClaim(appID *string, orgID string, id string) (*model.RewardClaim, error)
	CreateRewardClaim(appID *string, orgID string, item model.RewardClaim) (*model.RewardClaim, error)
	UpdateRewardClaim(appID *string, orgID string, id string, item model.RewardClaim) (*model.RewardClaim, error)

	GetUserRewardsHistory(appID *string, orgID string, userID string, rewardType *string, code *string, buildingBlock *string, limit *int64, offset *int64) ([]model.Reward, error)
	GetUserRewardByID(appID *string, orgID string, userID, id string) (*model.Reward, error)
	CreateUserReward(appID *string, orgID string, item model.Reward) (*model.Reward, error)

	// Quantities
	GetRewardQuantityState(appID *string, orgID string, rewardType string, inStock *bool) (*model.RewardQuantityState, error)

	// User APIs
	GetUserRewardsAmount(appID *string, orgID string, userID string, rewardType *string) ([]model.RewardTypeAmount, error)
	GetUserClaimsAmount(appID *string, orgID string, userID string, rewardType *string) ([]model.RewardTypeAmount, error)

	FindAllRewardTypeItems(context storage.TransactionContext) ([]model.RewardType, error)
	StoreMultiTenancyData(context storage.TransactionContext, appID string, orgID string) error

	SetListener(listener storage.Listener)
}
