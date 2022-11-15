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
	"fmt"
	"log"
	"rewards/core/model"
)

func (app *Application) getVersion() string {
	return app.version
}

func (app *Application) getRewardTypes(allApps bool, appID *string, orgID string) ([]model.RewardType, error) {
	types := app.cacheAdapter.GetRewardTypes()
	if types != nil {
		return types, nil
	}

	var appIDParam *string
	if !allApps {
		appIDParam = appID //associated with current app
	}

	storedTypes, err := app.storage.GetRewardTypes(appIDParam, orgID)
	if err == nil && storedTypes != nil {
		app.cacheAdapter.SetRewardTypes(storedTypes)
	}
	return storedTypes, err
}

func (app *Application) getRewardType(allApps bool, appID *string, orgID string, id string) (*model.RewardType, error) {
	var appIDParam *string
	if !allApps {
		appIDParam = appID //associated with current app
	}
	return app.storage.GetRewardType(appIDParam, orgID, id)
}

func (app *Application) createRewardType(allApps bool, appID *string, orgID string, item model.RewardType) (*model.RewardType, error) {
	var appIDParam *string
	if !allApps {
		appIDParam = appID //associated with current app
	}
	return app.storage.CreateRewardType(appIDParam, orgID, item)
}

func (app *Application) updateRewardType(allApps bool, appID *string, orgID string, id string, item model.RewardType) (*model.RewardType, error) {
	var appIDParam *string
	if !allApps {
		appIDParam = appID //associated with current app
	}
	return app.storage.UpdateRewardType(appIDParam, orgID, id, item)
}

func (app *Application) deleteRewardTypes(allApps bool, appID *string, orgID string, id string) error {
	var appIDParam *string
	if !allApps {
		appIDParam = appID //associated with current app
	}
	return app.storage.DeleteRewardType(appIDParam, orgID, id)
}

func (app *Application) getRewardOperations(allApps bool, appID *string, orgID string) ([]model.RewardOperation, error) {
	var appIDParam *string
	if !allApps {
		appIDParam = appID //associated with current app
	}
	return app.storage.GetRewardOperations(appIDParam, orgID)
}

func (app *Application) getRewardOperationByID(allApps bool, appID *string, orgID string, id string) (*model.RewardOperation, error) {
	var appIDParam *string
	if !allApps {
		appIDParam = appID //associated with current app
	}
	return app.storage.GetRewardOperationByID(appIDParam, orgID, id)
}

func (app *Application) getRewardOperationByCode(appID *string, orgID string, code string) (*model.RewardOperation, error) {
	return app.storage.GetRewardOperationByCode(appID, orgID, code)
}

func (app *Application) createRewardOperation(allApps bool, appID *string, orgID string, item model.RewardOperation) (*model.RewardOperation, error) {
	var appIDParam *string
	if !allApps {
		appIDParam = appID //associated with current app
	}
	return app.storage.CreateRewardOperation(appIDParam, orgID, item)
}

func (app *Application) updateRewardOperation(allApps bool, appID *string, orgID string, id string, item model.RewardOperation) (*model.RewardOperation, error) {
	var appIDParam *string
	if !allApps {
		appIDParam = appID //associated with current app
	}
	return app.storage.UpdateRewardOperation(appIDParam, orgID, id, item)
}

func (app *Application) deleteRewardOperation(allApps bool, appID *string, orgID string, id string) error {
	var appIDParam *string
	if !allApps {
		appIDParam = appID //associated with current app
	}
	return app.storage.DeleteRewardOperation(appIDParam, orgID, id)
}

func (app *Application) createReward(appID *string, orgID string, item model.Reward) (*model.Reward, error) {
	if item.RewardType != "" && item.UserID != "" {
		rewardType, err := app.storage.GetRewardTypeByType(appID, orgID, item.RewardType)
		if err != nil {
			log.Printf("Error Application.createReward(): %s", err)
			return nil, fmt.Errorf("Error Application.createReward(): %s", err)
		}

		if rewardType == nil {
			log.Printf("Error Application.createReward() unable to find reward type '%s'", item.RewardType)
			return nil, fmt.Errorf("Error Application.createReward() unable to find reward type '%s'", item.RewardType)
		}

		if item.Amount <= 0 {
			log.Printf("Error Application.createReward() amount is zero or a negative value")
			return nil, fmt.Errorf("Error Application.createReward() amount is zero or a negative value")
		}

		//TBD: Check for available quantity!!!
		quantity, err := app.storage.GetRewardQuantityState(appID, orgID, item.RewardType, nil)
		if err != nil {
			log.Printf("Error Application.createReward(): %s", err)
			return nil, fmt.Errorf("Error Application.createReward(): %s", err)
		}

		if quantity.GrantableQuantity >= item.Amount {
			return app.storage.CreateUserReward(appID, orgID, item)
		}
		return nil, fmt.Errorf("error Application.createReward(): not enough available quantity")
	}
	return nil, fmt.Errorf("Error Application.createReward(): missing data. data dump: %+v", item)
}

// Reward pools

func (app *Application) getRewardInventories(allApps bool, appID *string, orgID string, ids []string, rewardType *string, inStock *bool, grantDepleted *bool, claimDepleted *bool, limit *int64, offset *int64) ([]model.RewardInventory, error) {
	var appIDParam *string
	if !allApps {
		appIDParam = appID //associated with current app
	}
	return app.storage.GetRewardInventories(appIDParam, orgID, ids, rewardType, inStock, grantDepleted, claimDepleted, limit, offset)
}

func (app *Application) getRewardInventory(allApps bool, appID *string, orgID string, id string) (*model.RewardInventory, error) {
	var appIDParam *string
	if !allApps {
		appIDParam = appID //associated with current app
	}
	return app.storage.GetRewardInventory(appIDParam, orgID, id)
}

func (app *Application) createRewardInventory(allApps bool, appID *string, orgID string, item model.RewardInventory) (*model.RewardInventory, error) {
	var appIDParam *string
	if !allApps {
		appIDParam = appID //associated with current app
	}
	return app.storage.CreateRewardInventory(appIDParam, orgID, item)
}

func (app *Application) updateRewardInventory(allApps bool, appID *string, orgID string, id string, item model.RewardInventory) (*model.RewardInventory, error) {
	var appIDParam *string
	if !allApps {
		appIDParam = appID //associated with current app
	}
	return app.storage.UpdateRewardInventory(appIDParam, orgID, id, item)
}

func (app *Application) getRewardClaims(allApps bool, appID *string, orgID string, ids []string, userID *string, rewardType *string, status *string, limit *int64, offset *int64) ([]model.RewardClaim, error) {
	var appIDParam *string
	if !allApps {
		appIDParam = appID //associated with current app
	}
	return app.storage.GetRewardClaims(appIDParam, orgID, ids, userID, rewardType, status, limit, offset)
}

func (app *Application) getRewardClaim(allApps bool, appID *string, orgID string, id string) (*model.RewardClaim, error) {
	var appIDParam *string
	if !allApps {
		appIDParam = appID //associated with current app
	}
	return app.storage.GetRewardClaim(appIDParam, orgID, id)
}

func (app *Application) createRewardClaim(allApps bool, appID *string, orgID string, item model.RewardClaim) (*model.RewardClaim, error) {
	if len(item.Items) > 0 {
		balanceMapping, err := app.getUserBalanceMapping(appID, orgID, item.UserID)
		if err != nil {
			return nil, fmt.Errorf("Error on app.createRewardClaim() - %s", err)
		}

		for _, claimEntry := range item.Items {
			balance := balanceMapping[claimEntry.RewardType]
			if balance < claimEntry.Amount {
				return nil, fmt.Errorf("Error on app.createRewardClaim() - User(%s) not enough quantity for %s. Expected: %d, but have: %d", item.UserID, claimEntry.RewardType, claimEntry.Amount, balance)
			}

			inStock := true
			quantity, err := app.storage.GetRewardQuantityState(appID, orgID, claimEntry.RewardType, &inStock)
			if err != nil {
				return nil, fmt.Errorf("Error on app.createRewardClaim() - %s", err)
			}
			if quantity == nil || claimEntry.Amount > quantity.ClaimableQuantity {
				return nil, fmt.Errorf("Error on app.createRewardClaim() - not enough quantity for %s. Expected: %d", claimEntry.RewardType, claimEntry.Amount)
			}
		}
		var appIDParam *string
		if !allApps {
			appIDParam = appID //associated with current app
		}
		return app.storage.CreateRewardClaim(appIDParam, orgID, item)
	}
	return nil, fmt.Errorf("Error on app.createRewardClaim() - missing or zero quantity for reward items")
}

func (app *Application) updateRewardClaim(allApps bool, appID *string, orgID string, id string, item model.RewardClaim) (*model.RewardClaim, error) {
	var appIDParam *string
	if !allApps {
		appIDParam = appID //associated with current app
	}
	return app.storage.UpdateRewardClaim(appIDParam, orgID, id, item)
}

func (app *Application) getUserBalance(appID *string, orgID string, userID string) ([]model.RewardTypeAmount, error) {
	rewardsBalance, err := app.storage.GetUserRewardsAmount(appID, orgID, userID, nil)
	if err != nil {
		return nil, fmt.Errorf("Error app.getUserBalance() %s", err)
	}

	claimsBalance, err := app.storage.GetUserClaimsAmount(appID, orgID, userID, nil)
	if err != nil {
		return nil, fmt.Errorf("Error app.getUserBalance() %s", err)
	}
	claimsMapping := map[string]int{}
	if len(claimsBalance) > 0 {
		for _, claimBalance := range claimsBalance {
			claimsMapping[claimBalance.RewardType] = claimBalance.Amount
		}
	}

	if len(rewardsBalance) > 0 {
		for index, rewardBalance := range rewardsBalance {
			claimAmount := claimsMapping[rewardBalance.RewardType]
			rewardBalance.Amount -= claimAmount
			rewardsBalance[index] = rewardBalance
		}
	}

	return rewardsBalance, nil
}

func (app *Application) getUserBalanceMapping(appID *string, orgID string, userID string) (map[string]int, error) {
	rewardsBalance, err := app.storage.GetUserRewardsAmount(appID, orgID, userID, nil)
	if err != nil {
		return nil, fmt.Errorf("Error app.getUserBalanceMapping() %s", err)
	}

	claimsBalance, err := app.storage.GetUserClaimsAmount(appID, orgID, userID, nil)
	if err != nil {
		return nil, fmt.Errorf("Error app.getUserBalanceMapping() %s", err)
	}
	rewardsMapping := map[string]int{}
	if len(rewardsBalance) > 0 {
		for _, balance := range rewardsBalance {
			rewardsMapping[balance.RewardType] = balance.Amount
		}
	}

	if len(claimsBalance) > 0 {
		for _, claimdBalance := range claimsBalance {
			rewardsMapping[claimdBalance.RewardType] -= claimdBalance.Amount
		}
	}

	return rewardsMapping, nil
}

func (app *Application) getUserRewardsHistory(appID *string, orgID string, userID string, rewardType *string, code *string, buildingBlock *string, limit *int64, offset *int64) ([]model.Reward, error) {
	return app.storage.GetUserRewardsHistory(appID, orgID, userID, rewardType, code, buildingBlock, limit, offset)
}

func (app *Application) getRewardQuantity(appID *string, orgID string, rewardType string) (*model.RewardQuantityState, error) {
	return app.storage.GetRewardQuantityState(appID, orgID, rewardType, nil)
}

// OnRewardTypesChanged callback that indicates the reward types collection is changed
func (app *Application) OnRewardTypesChanged() {
	app.cacheAdapter.SetRewardTypes(nil) // invalidate
}
