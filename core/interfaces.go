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

	GetRewardTypes(orgID string) ([]model.RewardType, error)
	GetRewardType(orgID string, id string) (*model.RewardType, error)
	CreateRewardType(orgID string, item model.RewardType) (*model.RewardType, error)
	UpdateRewardType(orgID string, id string, item model.RewardType) (*model.RewardType, error)
	DeleteRewardType(orgID string, id string) error

	GetRewardInventories(orgID string, ids []string, rewardType *string) ([]model.RewardInventory, error)
	GetRewardInventory(orgID string, id string) (*model.RewardInventory, error)
	CreateRewardInventory(orgID string, item model.RewardInventory) (*model.RewardInventory, error)
	UpdateRewardInventory(orgID string, id string, item model.RewardInventory) (*model.RewardInventory, error)
	DeleteRewardInventory(orgID string, id string) error

	CreateRewardHistoryEntry(orgID string, item model.RewardHistoryEntry) (*model.RewardHistoryEntry, error)

	GetUserBalance(orgID string, userID string) (*model.WalletBalance, error)
	GetWalletBalance(orgID string, userID string, code string) (*model.WalletBalance, error)
	GetWalletHistoryEntries(orgID string, userID string) ([]model.RewardHistoryEntry, error)
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

func (s *servicesImpl) GetRewardInventories(orgID string, ids []string, rewardType *string) ([]model.RewardInventory, error) {
	return s.app.getRewardInventories(orgID, ids, rewardType)
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

func (s *servicesImpl) CreateRewardHistoryEntry(orgID string, item model.RewardHistoryEntry) (*model.RewardHistoryEntry, error) {
	return s.app.createRewardHistoryEntry(orgID, item)
}

func (s *servicesImpl) GetUserBalance(orgID string, userID string) (*model.WalletBalance, error) {
	return s.app.getUserBalance(orgID, userID)
}

func (s *servicesImpl) GetWalletBalance(orgID string, userID string, code string) (*model.WalletBalance, error) {
	return s.app.getWalletBalance(orgID, userID, code)
}

func (s *servicesImpl) GetWalletHistoryEntries(orgID string, userID string) ([]model.RewardHistoryEntry, error) {
	return s.app.getWalletHistoryEntries(orgID, userID)
}

// Storage is used by core to storage data - DB storage adapter, file storage adapter etc
type Storage interface {
	GetRewardTypes(orgID string) ([]model.RewardType, error)
	GetRewardType(orgID string, id string) (*model.RewardType, error)
	GetRewardTypeByType(orgID string, rewardType string) (*model.RewardType, error)
	CreateRewardType(orgID string, item model.RewardType) (*model.RewardType, error)
	UpdateRewardType(orgID string, id string, item model.RewardType) (*model.RewardType, error)
	DeleteRewardType(orgID string, id string) error

	GetRewardInventories(orgID string, ids []string, rewardType *string) ([]model.RewardInventory, error)
	GetRewardInventory(orgID string, id string) (*model.RewardInventory, error)
	CreateRewardInventory(orgID string, item model.RewardInventory) (*model.RewardInventory, error)
	UpdateRewardInventory(orgID string, id string, item model.RewardInventory) (*model.RewardInventory, error)
	DeleteRewardInventory(orgID string, id string) error

	GetRewardHistoryEntries(orgID string, userID string) ([]model.RewardHistoryEntry, error)
	GetRewardHistoryEntry(orgID string, userID, id string) (*model.RewardHistoryEntry, error)
	CreateRewardHistoryEntry(orgID string, item model.RewardHistoryEntry) (*model.RewardHistoryEntry, error)

	// User APIs
	GetUserBalance(orgID string, userID string) (*model.WalletBalance, error)
	GetWalletBalance(orgID string, userID string, code string) (*model.WalletBalance, error)

	SetListener(listener storage.Listener)
}
