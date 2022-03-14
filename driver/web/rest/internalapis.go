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
	OrgID       string `json:"org_id"`
	UserID      string `json:"user_id"`
	RewardType  string `json:"reward_type"`
	Description string `json:"description"`
} //@name createRewardHistoryEntryBody

// CreateRewardHistoryEntry Create a new reward history entry from another BB
// @Description Create a new reward history entry from another BB
// @Tags Internal
// @ID InternalCreateRewardHistoryEntry
// @Accept json
// @Success 200 {object} model.RewardHistoryEntry
// @Security InternalApiAuth
// @Router /int/reward_history [post]
func (h InternalApisHandler) CreateRewardHistoryEntry(w http.ResponseWriter, r *http.Request) {

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on adminapis.CreateRewardHistoryEntry: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var item createRewardHistoryEntryBody
	err = json.Unmarshal(data, &item)
	if err != nil {
		log.Printf("Error on adminapis.CreateRewardHistoryEntry: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdItem, err := h.app.Services.CreateRewardHistoryEntry(item.OrgID, model.RewardHistoryEntry{
		UserID:      item.UserID,
		RewardType:  item.RewardType,
		Description: item.Description,
	})
	if err != nil {
		log.Printf("Error on adminapis.CreateRewardHistoryEntry: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(createdItem)
	if err != nil {
		log.Printf("Error on adminapis.CreateRewardHistoryEntry: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
