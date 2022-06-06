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

package rest

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"rewards/core"
	"rewards/core/model"
)

// InternalApisHandler handles the rest internal APIs implementation
type InternalApisHandler struct {
	app *core.Application
}

// createRewardHistoryEntryBody wrapper
type createRewardHistoryEntryBody struct {
	OrgID         string `json:"org_id"`
	UserID        string `json:"user_id"`
	RewardType    string `json:"reward_type"`
	RewardCode    string `json:"code"`
	BuildingBlock string `json:"building_block"`
	Description   string `json:"description"`
} //@name createRewardHistoryEntryBody

// CreateReward Create a new reward history entry from another BB
// @Description Create a new reward history entry from another BB
// @Tags Internal
// @ID InternalCreateReward
// @Accept json
// @Success 200 {object} model.Reward
// @Security InternalApiAuth
// @Router /int/reward_history [post]
func (h InternalApisHandler) CreateReward(w http.ResponseWriter, r *http.Request) {

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on internalapis.CreateReward: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var item createRewardHistoryEntryBody
	err = json.Unmarshal(data, &item)
	if err != nil {
		log.Printf("Error on internalapis.CreateReward: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	operation, err := h.app.Services.GetRewardOperationByCode(item.OrgID, item.RewardCode)
	if err != nil {
		log.Printf("Error on internalapis.CreateReward: Reward operation not found. Error: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if operation != nil && item.BuildingBlock == operation.BuildingBlock && item.RewardCode == operation.Code && operation.Amount > 0 {
		createdItem, err := h.app.Services.CreateReward(item.OrgID, model.Reward{
			UserID:        item.UserID,
			RewardType:    operation.RewardType,
			Code:          operation.Code,
			BuildingBlock: operation.BuildingBlock,
			Description:   item.Description,
			Amount:        operation.Amount,
		})
		if err != nil {
			log.Printf("Error on internalapis.CreateReward: %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonData, err := json.Marshal(createdItem)
		if err != nil {
			log.Printf("Error on internalapis.CreateReward: %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
		return
	}

	log.Printf("Error on internalapis.CreateReward: Unable to find reward operation for the described code, type and building block or the amount of the operation is zero")
	http.Error(w, "Error on internalapis.CreateReward: Unable to find reward operation for the described code, type and building block or the amount of the operation is zero", http.StatusInternalServerError)
}

// getRewardStatsBody wrapper
type getRewardStatsBody struct {
	OrgID string `json:"org_id"`
} //@name getRewardStatsBody

// GetRewardStats Gets reward quantity stats for the current moment
// @Description Gets reward quantity stats for the current moment
// @Tags Internal
// @ID InternalGetRewardStats
// @Accept json
// @Success 200 {array} model.RewardQuantityState
// @Security InternalApiAuth
// @Router /int/reward_history [post]
func (h InternalApisHandler) GetRewardStats(w http.ResponseWriter, r *http.Request) {

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on internalapis.GetRewardStats: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var item getRewardStatsBody
	err = json.Unmarshal(data, &item)
	if err != nil {
		log.Printf("Error on internalapis.GetRewardStats: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	types, err := h.app.Services.GetRewardTypes(item.OrgID)
	if err != nil {
		log.Printf("Error on internalapis.GetRewardStats: Reward types not found. Error: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := []model.RewardQuantityState{}
	if len(types) > 0 {
		for _, rewardType := range types {
			quantity, err := h.app.Services.GetRewardQuantity(item.OrgID, rewardType.RewardType)
			if err != nil {
				log.Printf("Error on internalapis.GetRewardStats: %s", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if quantity != nil {
				result = append(result, *quantity)
			}
		}
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		log.Printf("Error on internalapis.GetRewardStats: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
	return
}
