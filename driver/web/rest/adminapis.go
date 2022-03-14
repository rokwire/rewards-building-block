package rest

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/rokwire/core-auth-library-go/tokenauth"
	"io/ioutil"
	"log"
	"net/http"
	"rewards/core"
	"rewards/core/model"
	"strings"
)

// AdminApisHandler handles the rest Admin APIs implementation
type AdminApisHandler struct {
	app *core.Application
}

// GetRewardTypes Retrieves  all reward types
// @Description Retrieves  all reward types
// @Tags Admin
// @ID AdminGetRewardTypes
// @Success 200 {array} model.RewardType
// @Security AdminUserAuth
// @Router /admin/reward_types [get]
func (h AdminApisHandler) GetRewardTypes(claims *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {
	resData, err := h.app.Services.GetRewardTypes(claims.OrgID)
	if err != nil {
		log.Printf("Error on adminapis.GetRewardTypes(): %s", err)
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
func (h AdminApisHandler) GetRewardType(claims *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	resData, err := h.app.Services.GetRewardType(claims.OrgID, id)
	if err != nil {
		log.Printf("Error on adminapis.GetRewardType(%s): %s", id, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(resData)
	if err != nil {
		log.Printf("Error on adminapis.GetRewardType(%s): %s", id, err)
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
func (h AdminApisHandler) UpdateRewardType(claims *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on adminapis.UpdateRewardType(%s): %s", id, err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var item model.RewardType
	err = json.Unmarshal(data, &item)
	if err != nil {
		log.Printf("Error on adminapis.UpdateRewardType(%s): %s", id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resData, err := h.app.Services.UpdateRewardType(claims.OrgID, id, item)
	if err != nil {
		log.Printf("Error on adminapis.UpdateRewardType(%s): %s", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(resData)
	if err != nil {
		log.Printf("Error on adminapis.UpdateRewardType(%s): %s", id, err)
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
func (h AdminApisHandler) CreateRewardType(claims *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on adminapis.CreateRewardType: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var item model.RewardType
	err = json.Unmarshal(data, &item)
	if err != nil {
		log.Printf("Error on adminapis.CreateRewardType: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdItem, err := h.app.Services.CreateRewardType(claims.OrgID, item)
	if err != nil {
		log.Printf("Error on adminapis.CreateRewardType: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(createdItem)
	if err != nil {
		log.Printf("Error on adminapis.CreateRewardType: %s", err)
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
func (h AdminApisHandler) DeleteRewardType(claims *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.app.Services.DeleteRewardType(claims.OrgID, id)
	if err != nil {
		log.Printf("Error on adminapis.DeleteRewardType(%s): %s", id, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

// GetRewardInventories Retrieves  all reward inventories
// @Description Retrieves  all reward types
// @Param ids query string false "Coma separated IDs of the desired records"
// @Tags Admin
// @ID AdminGetRewardInventories
// @Success 200 {array} model.RewardInventory
// @Security AdminUserAuth
// @Router /admin/reward_pools [get]
func (h AdminApisHandler) GetRewardInventories(claims *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {

	IDs := []string{}
	IDskeys, ok := r.URL.Query()["ids"]
	if ok && len(IDskeys[0]) > 0 {
		extIDs := IDskeys[0]
		IDs = strings.Split(extIDs, ",")
	}

	resData, err := h.app.Services.GetRewardInventories(claims.OrgID, IDs, nil)
	if err != nil {
		log.Printf("Error on adminapis.GetRewardInventories: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if resData == nil {
		resData = []model.RewardInventory{}
	}

	data, err := json.Marshal(resData)
	if err != nil {
		log.Printf("Error on adminapis.GetRewardInventories: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// GetRewardInventory Retrieves a reward inventory by id
// @Description Retrieves a reward inventory by id
// @Tags Admin
// @ID AdminGetRewardInventory
// @Accept json
// @Produce json
// @Success 200 {object} model.RewardInventory
// @Security AdminUserAuth
// @Router /admin/reward_pools/{id} [get]
func (h AdminApisHandler) GetRewardInventory(claims *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	resData, err := h.app.Services.GetRewardInventory(claims.OrgID, id)
	if err != nil {
		log.Printf("Error on adminapis.GetRewardInventory(%s): %s", id, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(resData)
	if err != nil {
		log.Printf("Error on adminapis.GetRewardInventory(%s): %s", id, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// UpdateRewardInventory Updates a reward inventory with the specified id
// @Description Updates a reward inventory with the specified id
// @Tags Admin
// @ID AdminUpdateRewardInventory
// @Accept json
// @Produce json
// @Success 200 {object} model.RewardInventory
// @Security AdminUserAuth
// @Router /admin/reward_pool/{id} [put]
func (h AdminApisHandler) UpdateRewardInventory(claims *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on adminapis.UpdateRewardInventory(%s): %s", id, err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var item model.RewardInventory
	err = json.Unmarshal(data, &item)
	if err != nil {
		log.Printf("Error on adminapis.UpdateRewardInventory(%s): %s", id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resData, err := h.app.Services.UpdateRewardInventory(claims.OrgID, id, item)
	if err != nil {
		log.Printf("Error on adminapis.UpdateRewardInventory(%s): %s", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(resData)
	if err != nil {
		log.Printf("Error on adminapis.UpdateRewardInventory(%s): %s", id, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// CreateRewardInventory Create a new reward inventory
// @Description Create a new reward inventory
// @Tags Admin
// @ID AdminCreateRewardInventory
// @Accept json
// @Success 200 {object} model.RewardInventory
// @Security AdminUserAuth
// @Router /admin/reward_pool [post]
func (h AdminApisHandler) CreateRewardInventory(claims *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on adminapis.CreateRewardInventory: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var item model.RewardInventory
	err = json.Unmarshal(data, &item)
	if err != nil {
		log.Printf("Error on adminapis.CreateRewardInventory: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdItem, err := h.app.Services.CreateRewardInventory(claims.OrgID, item)
	if err != nil {
		log.Printf("Error on adminapis.CreateRewardInventory: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(createdItem)
	if err != nil {
		log.Printf("Error on adminapis.CreateRewardInventory: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// DeleteRewardInventory Deletes a reward inventory with the specified id
// @Description Deletes a reward inventory with the specified id
// @Tags Admin
// @ID AdminDeleteRewardInventory
// @Success 200
// @Security AdminUserAuth
// @Router /admin/reward_pool/{id} [delete]
func (h AdminApisHandler) DeleteRewardInventory(claims *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.app.Services.DeleteRewardInventory(claims.OrgID, id)
	if err != nil {
		log.Printf("Error on adminapis.DeleteRewardInventory(%s): %s", id, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}
