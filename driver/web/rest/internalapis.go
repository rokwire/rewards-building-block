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

// CreateRewardHistoryEntry Create a new reward history entry from another BB
// @Description Create a new reward history entry from another BB
// @Tags Admin
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

	var item model.RewardHistoryEntry
	err = json.Unmarshal(data, &item)
	if err != nil {
		log.Printf("Error on adminapis.CreateRewardHistoryEntry: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdItem, err := h.app.Services.CreateRewardHistoryEntry(item)
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
