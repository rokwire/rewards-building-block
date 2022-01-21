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
)

// Services exposes APIs for the driver adapters
type Services interface {
	GetVersion() string

	GetRewardTypes(ids []string) ([]model.RewardType, error)
	GetRewardType(id string) (*model.RewardType, error)
	CreateRewardType(item model.RewardType) (*model.RewardType, error)
	UpdateRewardType(id string, item model.RewardType) (*model.RewardType, error)
	DeleteRewardType(id string) error

	GetRewardPools(ids []string) ([]model.RewardPool, error)
	GetRewardPool(id string) (*model.RewardPool, error)
	CreateRewardPool(item model.RewardPool) (*model.RewardPool, error)
	UpdateRewardPool(id string, item model.RewardPool) (*model.RewardPool, error)
	DeleteRewardPool(id string) error
}

type servicesImpl struct {
	app *Application
}

func (s *servicesImpl) GetVersion() string {
	return s.app.getVersion()
}

func (s *servicesImpl) GetRewardTypes(ids []string) ([]model.RewardType, error) {
	return s.app.getRewardTypes(ids)
}

func (s *servicesImpl) GetRewardType(id string) (*model.RewardType, error) {
	return s.app.getRewardType(id)
}

func (s *servicesImpl) CreateRewardType(item model.RewardType) (*model.RewardType, error) {
	return s.app.createRewardType(item)
}

func (s *servicesImpl) UpdateRewardType(id string, item model.RewardType) (*model.RewardType, error) {
	return s.app.updateRewardType(id, item)
}

func (s *servicesImpl) DeleteRewardType(id string) error {
	return s.app.deleteGetRewardTypes(id)
}

func (s *servicesImpl) GetRewardPools(ids []string) ([]model.RewardPool, error) {
	return s.app.getRewardPools(ids)
}

func (s *servicesImpl) GetRewardPool(id string) (*model.RewardPool, error) {
	return s.app.getRewardPool(id)
}

func (s *servicesImpl) CreateRewardPool(item model.RewardPool) (*model.RewardPool, error) {
	return s.app.createRewardPool(item)
}

func (s *servicesImpl) UpdateRewardPool(id string, item model.RewardPool) (*model.RewardPool, error) {
	return s.app.updateRewardPool(id, item)
}

func (s *servicesImpl) DeleteRewardPool(id string) error {
	return s.app.deleteGetRewardTypes(id)
}

// Storage is used by core to storage data - DB storage adapter, file storage adapter etc
type Storage interface {
	GetRewardTypes(ids []string) ([]model.RewardType, error)
	GetRewardType(id string) (*model.RewardType, error)
	CreateRewardType(item model.RewardType) (*model.RewardType, error)
	UpdateRewardType(id string, item model.RewardType) (*model.RewardType, error)
	DeleteRewardType(id string) error

	GetRewardPools(ids []string) ([]model.RewardPool, error)
	GetRewardPool(id string) (*model.RewardPool, error)
	CreateRewardPool(item model.RewardPool) (*model.RewardPool, error)
	UpdateRewardPool(id string, item model.RewardPool) (*model.RewardPool, error)
	DeleteRewardPool(id string) error
}
