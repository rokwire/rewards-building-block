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

func (app *Application) getRewardTypes() ([]model.RewardType, error) {
	types := app.cacheAdapter.GetRewardTypes()
	if types != nil {
		return types, nil
	}

	storedTypes, err := app.storage.GetRewardTypes()
	if err == nil && storedTypes != nil {
		app.cacheAdapter.SetRewardTypes(storedTypes)
	}
	return storedTypes, err
}

func (app *Application) getRewardType(id string) (*model.RewardType, error) {
	return app.storage.GetRewardType(id)
}

func (app *Application) createRewardType(item model.RewardType) (*model.RewardType, error) {
	return  app.storage.CreateRewardType(item)
}

func (app *Application) updateRewardType(id string, item model.RewardType) (*model.RewardType, error) {
	return app.storage.UpdateRewardType(id, item)
}

func (app *Application) deleteGetRewardTypes(id string) error {
	return app.storage.DeleteRewardType(id)
}

func (app *Application) createRewardHistoryEntry(item model.RewardHistoryEntry) (*model.RewardHistoryEntry, error) {
	if item.RewardType != "" && item.UserID != "" {
		rewardType, err := app.storage.GetRewardTypeByType(item.RewardType)
		if err != nil {
			log.Printf("Error Application.createRewardHistoryEntry(): %s", err)
			return nil, fmt.Errorf("Error Application.createRewardHistoryEntry(): %s", err)
		}

		if rewardType == nil {
			log.Printf("Error Application.createRewardHistoryEntry() unable to find reward type '%s'", item.RewardType)
			return nil, fmt.Errorf("Error Application.createRewardHistoryEntry() unable to find reward type '%s'", item.RewardType)
		}

		item.Amount = rewardType.Amount

		return app.storage.CreateRewardHistoryEntry(item)
	}
	return nil, fmt.Errorf("Error Application.createRewardHistoryEntry(): missing data. data dump: %+v", item)
}

// Reward pools

func (app *Application) getRewardPools(ids []string) ([]model.RewardPool, error) {
	return app.storage.GetRewardPools(ids)
}

func (app *Application) getRewardPool(id string) (*model.RewardPool, error) {
	return app.storage.GetRewardPool(id)
}

func (app *Application) createRewardPool(item model.RewardPool) (*model.RewardPool, error) {
	return app.storage.CreateRewardPool(item)
}

func (app *Application) updateRewardPool(id string, item model.RewardPool) (*model.RewardPool, error) {
	return app.storage.UpdateRewardPool(id, item)
}

func (app *Application) deleteGetRewardPool(id string) error {
	return app.storage.DeleteRewardPool(id)
}

func (app *Application) getUserBalance(userID string) (*model.WalletBalance, error) {
	return app.storage.GetUserBalance(userID)
}

func (app *Application) getWalletBalance(userID string, code string) (*model.WalletBalance, error) {
	return app.storage.GetWalletBalance(userID, code)
}

func (app *Application) getWalletHistoryEntries(userID string) ([]model.RewardHistoryEntry, error) {
	history, err := app.storage.GetRewardHistoryEntries(userID)
	if err != nil {
		return nil, err
	}

	rewardTypes, err := app.Services.GetRewardTypes()
	if err != nil {
		log.Printf("Error on apis.GetRewardTypes(): %s", err)
	} else {
		if len(rewardTypes) > 0 && len(history) > 0 {
			mapping := map[string]model.RewardType{}
			for _, rewardType := range rewardTypes {
				mapping[rewardType.RewardType] = rewardType
			}
			for i, historyItem := range history {
				displayName := mapping[historyItem.RewardType].DisplayName
				history[i].DisplayName = &displayName
			}
		}
	}

	return history, nil
}


func (app *Application) OnRewardTypesChanged(){
	app.cacheAdapter.SetRewardTypes(nil) // invalidate
}