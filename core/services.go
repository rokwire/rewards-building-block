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

func (app *Application) getRewardOperations(orgID string) ([]model.RewardOperation, error) {
	return app.storage.GetRewardOperations(orgID)
}

func (app *Application) getRewardOperationByID(orgID string, id string) (*model.RewardOperation, error) {
	return app.storage.GetRewardOperationByID(orgID, id)
}

func (app *Application) getRewardOperationByCode(orgID string, code string) (*model.RewardOperation, error) {
	return app.storage.GetRewardOperationByCode(orgID, code)
}

func (app *Application) createRewardOperation(orgID string, item model.RewardOperation) (*model.RewardOperation, error) {
	return app.storage.CreateRewardOperation(orgID, item)
}

func (app *Application) updateRewardOperation(orgID string, id string, item model.RewardOperation) (*model.RewardOperation, error) {
	return app.storage.UpdateRewardOperation(orgID, id, item)
}

func (app *Application) deleteRewardOperation(orgID string, id string) error {
	return app.storage.DeleteRewardOperation(orgID, id)
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

		if item.Amount <= 0 {
			log.Printf("Error Application.createReward() amount is zero or a negative value")
			return nil, fmt.Errorf("Error Application.createReward() amount is zero or a negative value")
		}

		//TBD: Check for available quantity!!!
		quantity, err := app.storage.GetRewardQuantityState(orgID, item.RewardType, nil)
		if err != nil {
			log.Printf("Error Application.createReward(): %s", err)
			return nil, fmt.Errorf("Error Application.createReward(): %s", err)
		}

		if quantity.GrantableQuantity >= item.Amount {
			return app.storage.CreateUserReward(orgID, item)
		}
		return nil, fmt.Errorf("error Application.createReward(): not enough available quantity")
	}
	return nil, fmt.Errorf("Error Application.createReward(): missing data. data dump: %+v", item)
}

// Reward pools

func (app *Application) getRewardInventories(orgID string, ids []string, rewardType *string, inStock *bool, grantDepleted *bool, claimDepleted *bool, limit *int64, offset *int64) ([]model.RewardInventory, error) {
	return app.storage.GetRewardInventories(orgID, ids, rewardType, inStock, grantDepleted, claimDepleted, limit, offset)
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

func (app *Application) getRewardClaims(orgID string, ids []string, userID *string, rewardType *string, status *string, limit *int64, offset *int64) ([]model.RewardClaim, error) {
	return app.storage.GetRewardClaims(orgID, ids, userID, rewardType, status, limit, offset)
}

func (app *Application) getRewardClaim(orgID string, id string) (*model.RewardClaim, error) {
	return app.storage.GetRewardClaim(orgID, id)
}

func (app *Application) createRewardClaim(orgID string, item model.RewardClaim) (*model.RewardClaim, error) {
	if len(item.Items) > 0 {
		balanceMapping, err := app.getUserBalanceMapping(orgID, item.UserID)
		if err != nil {
			return nil, fmt.Errorf("Error on app.createRewardClaim() - %s", err)
		}

		for _, claimEntry := range item.Items {
			balance := balanceMapping[claimEntry.RewardType]
			if balance < claimEntry.Amount {
				return nil, fmt.Errorf("Error on app.createRewardClaim() - User(%s) not enough quantity for %s. Expected: %d, but have: %d", item.UserID, claimEntry.RewardType, claimEntry.Amount, balance)
			}

			inStock := true
			quantity, err := app.storage.GetRewardQuantityState(orgID, claimEntry.RewardType, &inStock)
			if err != nil {
				return nil, fmt.Errorf("Error on app.createRewardClaim() - %s", err)
			}
			if quantity == nil || claimEntry.Amount > quantity.ClaimableQuantity {
				return nil, fmt.Errorf("Error on app.createRewardClaim() - not enough quantity for %s. Expected: %d", claimEntry.RewardType, claimEntry.Amount)
			}
		}
		return app.storage.CreateRewardClaim(orgID, item)
	}
	return nil, fmt.Errorf("Error on app.createRewardClaim() - missing or zero quantity for reward items")
}

func (app *Application) updateRewardClaim(orgID string, id string, item model.RewardClaim) (*model.RewardClaim, error) {
	return app.storage.UpdateRewardClaim(orgID, id, item)
}

func (app *Application) getUserBalance(orgID string, userID string) ([]model.RewardTypeAmount, error) {
	rewardsBalance, err := app.storage.GetUserRewardsAmount(orgID, userID, nil)
	if err != nil {
		return nil, fmt.Errorf("Error app.getUserBalance() %s", err)
	}

	claimsBalance, err := app.storage.GetUserClaimsAmount(orgID, userID, nil)
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

func (app *Application) getUserBalanceMapping(orgID string, userID string) (map[string]int, error) {
	rewardsBalance, err := app.storage.GetUserRewardsAmount(orgID, userID, nil)
	if err != nil {
		return nil, fmt.Errorf("Error app.getUserBalanceMapping() %s", err)
	}

	claimsBalance, err := app.storage.GetUserClaimsAmount(orgID, userID, nil)
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

func (app *Application) getUserRewardsHistory(orgID string, userID string, rewardType *string, code *string, buildingBlock *string, limit *int64, offset *int64) ([]model.Reward, error) {
	return app.storage.GetUserRewardsHistory(orgID, userID, rewardType, code, buildingBlock, limit, offset)
}

func (app *Application) getRewardQuantity(orgID string, rewardType string) (*model.RewardQuantityState, error) {
	return app.storage.GetRewardQuantityState(orgID, rewardType, nil)
}

// OnRewardTypesChanged callback that indicates the reward types collection is changed
func (app *Application) OnRewardTypesChanged() {
	app.cacheAdapter.SetRewardTypes(nil) // invalidate
}
