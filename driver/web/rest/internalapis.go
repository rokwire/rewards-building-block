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
	RewardCode    string `json:"reward_code"`
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
		log.Printf("Error on adminapis.CreateUserReward: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var item createRewardHistoryEntryBody
	err = json.Unmarshal(data, &item)
	if err != nil {
		log.Printf("Error on adminapis.CreateUserReward: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	operation, err := h.app.Services.GetRewardOperationByCode(item.OrgID, item.RewardCode)
	if err != nil {
		log.Printf("Error on adminapis.GetRewardOperationByCode: Reward operation not found. Error: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if operation != nil && item.BuildingBlock == operation.BuildingBlock && item.RewardType == operation.RewardType && operation.Amount > 0 {
		createdItem, err := h.app.Services.CreateReward(item.OrgID, model.Reward{
			UserID:      item.UserID,
			RewardType:  item.RewardType,
			Description: item.Description,
			Amount:      operation.Amount,
		})
		if err != nil {
			log.Printf("Error on adminapis.CreateUserReward: %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonData, err := json.Marshal(createdItem)
		if err != nil {
			log.Printf("Error on adminapis.CreateUserReward: %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
		return
	}

	log.Printf("Error on adminapis.CreateUserReward: Unable to find reward operation for the described code, type and building block or the amount of the operation is zero")
	http.Error(w, "Error on adminapis.CreateUserReward: Unable to find reward operation for the described code, type and building block or the amount of the operation is zero", http.StatusInternalServerError)
}
