package rest

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"rewards/core"
	"rewards/core/model"
	"strings"
)

//AdminApisHandler handles the rest Admin APIs implementation
type AdminApisHandler struct {
	app *core.Application
}

// GetRewardTypes Retrieves  all reward types
// @Description Retrieves  all reward types
// @Param ids query string false "Coma separated IDs of the desired records"
// @Tags Admin
// @ID AdminGetRewardTypes
// @Success 200 {array} model.RewardType
// @Security AdminUserAuth
// @Router /admin/reward_types [get]
func (h AdminApisHandler) GetRewardTypes(w http.ResponseWriter, r *http.Request) {

	IDs := []string{}
	IDskeys, ok := r.URL.Query()["ids"]
	if ok && len(IDskeys[0]) > 0 {
		extIDs := IDskeys[0]
		IDs = strings.Split(extIDs, ",")
	}

	resData, err := h.app.Services.GetRewardTypes(IDs)
	if err != nil {
		log.Printf("Error on getting reward types: %s\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if resData == nil {
		resData = []model.RewardType{}
	}

	data, err := json.Marshal(resData)
	if err != nil {
		log.Printf("Error on marshal reward types: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// GetRewardType Retrieves a reward type by id
// @Description Retrieves a reward type by id
// @Tags Admin
// @ID AdminRewardTypes
// @Accept json
// @Produce json
// @Success 200 {object} model.RewardType
// @Security AdminUserAuth
// @Router /admin/reward_types/{id} [get]
func (h AdminApisHandler) GetRewardType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	resData, err := h.app.Services.GetRewardType(id)
	if err != nil {
		log.Printf("Error on getting reward type id - %s\n %s", id, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(resData)
	if err != nil {
		log.Println("Error on marshal the reward type")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// UpdateRewardType Updates a reward type with the specified id
// @Description Updates a reward type with the specified id
// @Tags Admin
// @ID AdminUpdateRewardType
// @Accept json
// @Produce json
// @Success 200 {object} model.RewardType
// @Security AdminUserAuth
// @Router /admin/reward_types/{id} [put]
func (h AdminApisHandler) UpdateRewardType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal update reward type - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var item model.RewardType
	err = json.Unmarshal(data, &item)
	if err != nil {
		log.Printf("Error on unmarshal the update reward type request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resData, err := h.app.Services.UpdateRewardType(id, item)
	if err != nil {
		log.Printf("Error on update reward type with id - %s\n %s", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(resData)
	if err != nil {
		log.Printf("Error on marshal the updated  reward type: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// CreateRewardType Create a new reward type
// @Description Create a new reward type
// @Tags Admin
// @ID AdminCreateRewardType
// @Accept json
// @Success 200 {object} model.RewardType
// @Security AdminUserAuth
// @Router /admin/reward_types [post]
func (h AdminApisHandler) CreateRewardType(w http.ResponseWriter, r *http.Request) {

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal create a reward type - %s\n", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var item model.RewardType
	err = json.Unmarshal(data, &item)
	if err != nil {
		log.Printf("Error on unmarshal the create reward type request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdItem, err := h.app.Services.CreateRewardType(item)
	if err != nil {
		log.Printf("Error on creating reward type: %s\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(createdItem)
	if err != nil {
		log.Printf("Error on marshal creating reward type: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// DeleteRewardType Deletes a reward type with the specified id
// @Description Deletes a reward type with the specified id
// @Tags Admin
// @ID AdminDeleteRewardType
// @Success 200
// @Security AdminUserAuth
// @Router /admin/reward_types/{id} [delete]
func (h AdminApisHandler) DeleteRewardType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.app.Services.DeleteRewardTypes(id)
	if err != nil {
		log.Printf("Error on deleting reward type with id - %s\n %s", id, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}
