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
	"fmt"
	"log"
	"rewards/core/model"
)

func (app *Application) getVersion() string {
	return app.version
}

func (app *Application) getRewardTypes(orgID string) ([]model.RewardType, error) {
	types := app.cacheAdapter.GetRewardTypes()
	if types != nil {
		return types, nil
	}

	storedTypes, err := app.storage.GetRewardTypes(orgID)
	if err == nil && storedTypes != nil {
		app.cacheAdapter.SetRewardTypes(storedTypes)
	}
	return storedTypes, err
}

func (app *Application) getRewardType(orgID string, id string) (*model.RewardType, error) {
	return app.storage.GetRewardType(orgID, id)
}

func (app *Application) createRewardType(orgID string, item model.RewardType) (*model.RewardType, error) {
	return app.storage.CreateRewardType(orgID, item)
}

func (app *Application) updateRewardType(orgID string, id string, item model.RewardType) (*model.RewardType, error) {
	return app.storage.UpdateRewardType(orgID, id, item)
}

func (app *Application) deleteRewardTypes(orgID string, id string) error {
	return app.storage.DeleteRewardType(orgID, id)
}

func (app *Application) createReward(orgID string, item model.Reward) (*model.Reward, error) {
	if item.RewardType != "" && item.UserID != "" {
		rewardType, err := app.storage.GetRewardTypeByType(orgID, item.RewardType)
		if err != nil {
			log.Printf("Error Application.createReward(): %s", err)
			return nil, fmt.Errorf("Error Application.createReward(): %s", err)
		}

		if rewardType == nil {
			log.Printf("Error Application.createReward() unable to find reward type '%s'", item.RewardType)
			return nil, fmt.Errorf("Error Application.createReward() unable to find reward type '%s'", item.RewardType)
		}

		//item.Amount = rewardType.Amount

		return app.storage.CreateReward(orgID, item)
	}
	return nil, fmt.Errorf("Error Application.createReward(): missing data. data dump: %+v", item)
}

// Reward pools

func (app *Application) getRewardInventories(orgID string, ids []string, rewardType *string) ([]model.RewardInventory, error) {
	return app.storage.GetRewardInventories(orgID, ids, rewardType)
}

func (app *Application) getRewardInventory(orgID string, id string) (*model.RewardInventory, error) {
	return app.storage.GetRewardInventory(orgID, id)
}

func (app *Application) createRewardInventory(orgID string, item model.RewardInventory) (*model.RewardInventory, error) {
	return app.storage.CreateRewardInventory(orgID, item)
}

func (app *Application) updateRewardInventory(orgID string, id string, item model.RewardInventory) (*model.RewardInventory, error) {
	return app.storage.UpdateRewardInventory(orgID, id, item)
}

func (app *Application) deleteGetRewardInventory(orgID string, id string) error {
	return app.storage.DeleteRewardInventory(orgID, id)
}

func (app *Application) getUserBalance(orgID string, userID string) (*model.WalletBalance, error) {
	return app.storage.GetUserBalance(orgID, userID)
}

func (app *Application) getWalletBalance(orgID string, userID string, code string) (*model.WalletBalance, error) {
	return app.storage.GetWalletBalance(orgID, userID, code)
}

func (app *Application) getWalletHistoryEntries(orgID string, userID string) ([]model.Reward, error) {
	history, err := app.storage.GetRewardHistoryEntries(orgID, userID)
	if err != nil {
		return nil, err
	}

	rewardTypes, err := app.Services.GetRewardTypes(orgID)
	if err != nil {
		log.Printf("Error on apis.GetRewardTypes(): %s", err)
	} else {
		if len(rewardTypes) > 0 && len(history) > 0 {
			mapping := map[string]model.RewardType{}
			for _, rewardType := range rewardTypes {
				mapping[rewardType.RewardType] = rewardType
			}
		}
	}

	return history, nil
}

// OnRewardTypesChanged callback that indicates the reward types collection is changed
func (app *Application) OnRewardTypesChanged() {
	app.cacheAdapter.SetRewardTypes(nil) // invalidate
}
