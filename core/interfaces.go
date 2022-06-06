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

package core

import (
	"rewards/core/model"
	"rewards/driven/storage"
)

// Services exposes APIs for the driver adapters
type Services interface {
	GetVersion() string

	GetRewardTypes(orgID string) ([]model.RewardType, error)
	GetRewardType(orgID string, id string) (*model.RewardType, error)
	CreateRewardType(orgID string, item model.RewardType) (*model.RewardType, error)
	UpdateRewardType(orgID string, id string, item model.RewardType) (*model.RewardType, error)
	DeleteRewardType(orgID string, id string) error

	GetRewardOperations(orgID string) ([]model.RewardOperation, error)
	GetRewardOperationByID(orgID string, id string) (*model.RewardOperation, error)
	GetRewardOperationByCode(orgID string, code string) (*model.RewardOperation, error)
	CreateRewardOperation(orgID string, item model.RewardOperation) (*model.RewardOperation, error)
	UpdateRewardOperation(orgID string, id string, item model.RewardOperation) (*model.RewardOperation, error)
	DeleteRewardOperation(orgID string, id string) error

	GetRewardInventories(orgID string, ids []string, rewardType *string, inStock *bool, grantDepleted *bool, claimDepleted *bool, limit *int64, offset *int64) ([]model.RewardInventory, error)
	GetRewardInventory(orgID string, id string) (*model.RewardInventory, error)
	CreateRewardInventory(orgID string, item model.RewardInventory) (*model.RewardInventory, error)
	UpdateRewardInventory(orgID string, id string, item model.RewardInventory) (*model.RewardInventory, error)

	GetRewardClaims(orgID string, ids []string, userID *string, rewardType *string, status *string, limit *int64, offset *int64) ([]model.RewardClaim, error)
	GetRewardClaim(orgID string, id string) (*model.RewardClaim, error)
	CreateRewardClaim(orgID string, item model.RewardClaim) (*model.RewardClaim, error)
	UpdateRewardClaim(orgID string, id string, item model.RewardClaim) (*model.RewardClaim, error)

	CreateReward(orgID string, item model.Reward) (*model.Reward, error)

	GetUserBalance(orgID string, userID string) ([]model.RewardTypeAmount, error)
	GetUserRewardsHistory(orgID string, userID string, rewardType *string, code *string, buildingBlock *string, limit *int64, offset *int64) ([]model.Reward, error)

	GetRewardQuantity(orgID string, rewardType string) (*model.RewardQuantityState, error)
}

type servicesImpl struct {
	app *Application
}

func (s *servicesImpl) GetVersion() string {
	return s.app.getVersion()
}

func (s *servicesImpl) GetRewardTypes(orgID string) ([]model.RewardType, error) {
	return s.app.getRewardTypes(orgID)
}

func (s *servicesImpl) GetRewardType(orgID string, id string) (*model.RewardType, error) {
	return s.app.getRewardType(orgID, id)
}

func (s *servicesImpl) CreateRewardType(orgID string, item model.RewardType) (*model.RewardType, error) {
	return s.app.createRewardType(orgID, item)
}

func (s *servicesImpl) UpdateRewardType(orgID string, id string, item model.RewardType) (*model.RewardType, error) {
	return s.app.updateRewardType(orgID, id, item)
}

func (s *servicesImpl) DeleteRewardType(orgID string, id string) error {
	return s.app.deleteRewardTypes(orgID, id)
}

func (s *servicesImpl) GetRewardOperations(orgID string) ([]model.RewardOperation, error) {
	return s.app.getRewardOperations(orgID)
}

func (s *servicesImpl) GetRewardOperationByID(orgID string, id string) (*model.RewardOperation, error) {
	return s.app.getRewardOperationByID(orgID, id)
}

func (s *servicesImpl) GetRewardOperationByCode(orgID string, code string) (*model.RewardOperation, error) {
	return s.app.getRewardOperationByCode(orgID, code)
}

func (s *servicesImpl) CreateRewardOperation(orgID string, item model.RewardOperation) (*model.RewardOperation, error) {
	return s.app.createRewardOperation(orgID, item)
}

func (s *servicesImpl) UpdateRewardOperation(orgID string, id string, item model.RewardOperation) (*model.RewardOperation, error) {
	return s.app.updateRewardOperation(orgID, id, item)
}

func (s *servicesImpl) DeleteRewardOperation(orgID string, id string) error {
	return s.app.deleteRewardOperation(orgID, id)
}

func (s *servicesImpl) GetRewardInventories(orgID string, ids []string, rewardType *string, inStock *bool, grantDepleted *bool, claimDepleted *bool, limit *int64, offset *int64) ([]model.RewardInventory, error) {
	return s.app.getRewardInventories(orgID, ids, rewardType, inStock, grantDepleted, claimDepleted, limit, offset)
}

func (s *servicesImpl) GetRewardInventory(orgID string, id string) (*model.RewardInventory, error) {
	return s.app.getRewardInventory(orgID, id)
}

func (s *servicesImpl) CreateRewardInventory(orgID string, item model.RewardInventory) (*model.RewardInventory, error) {
	return s.app.createRewardInventory(orgID, item)
}

func (s *servicesImpl) UpdateRewardInventory(orgID string, id string, item model.RewardInventory) (*model.RewardInventory, error) {
	return s.app.updateRewardInventory(orgID, id, item)
}

func (s *servicesImpl) DeleteRewardInventory(orgID string, id string) error {
	return s.app.deleteRewardTypes(orgID, id)
}

func (s *servicesImpl) CreateReward(orgID string, item model.Reward) (*model.Reward, error) {
	return s.app.createReward(orgID, item)
}

func (s *servicesImpl) GetRewardClaims(orgID string, ids []string, userID *string, rewardType *string, status *string, limit *int64, offset *int64) ([]model.RewardClaim, error) {
	return s.app.getRewardClaims(orgID, ids, userID, rewardType, status, limit, offset)
}

func (s *servicesImpl) GetRewardClaim(orgID string, id string) (*model.RewardClaim, error) {
	return s.app.getRewardClaim(orgID, id)
}

func (s *servicesImpl) CreateRewardClaim(orgID string, item model.RewardClaim) (*model.RewardClaim, error) {
	return s.app.createRewardClaim(orgID, item)
}

func (s *servicesImpl) UpdateRewardClaim(orgID string, id string, item model.RewardClaim) (*model.RewardClaim, error) {
	return s.app.updateRewardClaim(orgID, id, item)
}

func (s *servicesImpl) GetUserBalance(orgID string, userID string) ([]model.RewardTypeAmount, error) {
	return s.app.getUserBalance(orgID, userID)
}

func (s *servicesImpl) GetUserRewardsHistory(orgID string, userID string, rewardType *string, code *string, buildingBlock *string, limit *int64, offset *int64) ([]model.Reward, error) {
	return s.app.getUserRewardsHistory(orgID, userID, rewardType, code, buildingBlock, limit, offset)
}

func (s *servicesImpl) GetRewardQuantity(orgID string, rewardType string) (*model.RewardQuantityState, error) {
	return s.app.getRewardQuantity(orgID, rewardType)
}

// Storage is used by core to storage data - DB storage adapter, file storage adapter etc
type Storage interface {
	GetRewardTypes(orgID string) ([]model.RewardType, error)
	GetRewardType(orgID string, id string) (*model.RewardType, error)
	GetRewardTypeByType(orgID string, rewardType string) (*model.RewardType, error)
	CreateRewardType(orgID string, item model.RewardType) (*model.RewardType, error)
	UpdateRewardType(orgID string, id string, item model.RewardType) (*model.RewardType, error)
	DeleteRewardType(orgID string, id string) error

	GetRewardOperations(orgID string) ([]model.RewardOperation, error)
	GetRewardOperationByID(orgID string, id string) (*model.RewardOperation, error)
	GetRewardOperationByCode(orgID string, code string) (*model.RewardOperation, error)
	CreateRewardOperation(orgID string, item model.RewardOperation) (*model.RewardOperation, error)
	UpdateRewardOperation(orgID string, id string, item model.RewardOperation) (*model.RewardOperation, error)
	DeleteRewardOperation(orgID string, id string) error

	GetRewardInventories(orgID string, ids []string, rewardType *string, inStock *bool, grantDepleted *bool, claimDepleted *bool, limit *int64, offset *int64) ([]model.RewardInventory, error)
	GetRewardInventory(orgID string, id string) (*model.RewardInventory, error)
	CreateRewardInventory(orgID string, item model.RewardInventory) (*model.RewardInventory, error)
	UpdateRewardInventory(orgID string, id string, item model.RewardInventory) (*model.RewardInventory, error)

	GetRewardClaims(orgID string, ids []string, userID *string, rewardType *string, status *string, limit *int64, offset *int64) ([]model.RewardClaim, error)
	GetRewardClaim(orgID string, id string) (*model.RewardClaim, error)
	CreateRewardClaim(orgID string, item model.RewardClaim) (*model.RewardClaim, error)
	UpdateRewardClaim(orgID string, id string, item model.RewardClaim) (*model.RewardClaim, error)

	GetUserRewardsHistory(orgID string, userID string, rewardType *string, code *string, buildingBlock *string, limit *int64, offset *int64) ([]model.Reward, error)
	GetUserRewardByID(orgID string, userID, id string) (*model.Reward, error)
	CreateUserReward(orgID string, item model.Reward) (*model.Reward, error)

	// Quantities
	GetRewardQuantityState(orgID string, rewardType string, inStock *bool) (*model.RewardQuantityState, error)

	// User APIs
	GetUserRewardsAmount(orgID string, userID string, rewardType *string) ([]model.RewardTypeAmount, error)
	GetUserClaimsAmount(orgID string, userID string, rewardType *string) ([]model.RewardTypeAmount, error)

	SetListener(listener storage.Listener)
}
