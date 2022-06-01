package rest

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"rewards/core"
	"rewards/core/model"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rokwire/core-auth-library-go/tokenauth"
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
// @Router /admin/types [get]
func (h AdminApisHandler) GetRewardTypes(claims *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {
	resData, err := h.app.Services.GetRewardTypes(&claims.AppID, claims.OrgID)
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
// @Router /admin/types/{id} [get]
func (h AdminApisHandler) GetRewardType(claims *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	resData, err := h.app.Services.GetRewardType(&claims.AppID, claims.OrgID, id)
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
// @Param data body model.RewardType true "body json"
// @Accept json
// @Produce json
// @Success 200 {object} model.RewardType
// @Security AdminUserAuth
// @Router /admin/types/{id} [put]
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

	resData, err := h.app.Services.UpdateRewardType(&claims.AppID, claims.OrgID, id, item)
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
// @Param data body model.RewardType true "body json"
// @Accept json
// @Success 200 {object} model.RewardType
// @Security AdminUserAuth
// @Router /admin/types [post]
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

	createdItem, err := h.app.Services.CreateRewardType(&claims.AppID, claims.OrgID, item)
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
// @Router /admin/types/{id} [delete]
func (h AdminApisHandler) DeleteRewardType(claims *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.app.Services.DeleteRewardType(&claims.AppID, claims.OrgID, id)
	if err != nil {
		log.Printf("Error on adminapis.DeleteRewardType(%s): %s", id, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

// GetRewardOperations Retrieves  all reward operations
// @Description Retrieves  all reward types
// @Tags Admin
// @ID AdminGetRewardOperations
// @Success 200 {array} model.RewardOperation
// @Security AdminUserAuth
// @Router /admin/operations [get]
func (h AdminApisHandler) GetRewardOperations(claims *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {
	resData, err := h.app.Services.GetRewardTypes(&claims.AppID, claims.OrgID)
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

// GetRewardOperation Retrieves a reward operation by id
// @Description Retrieves a reward operation by id
// @Tags Admin
// @ID AdminGetRewardOperation
// @Accept json
// @Produce json
// @Success 200 {object} model.RewardOperation
// @Security AdminUserAuth
// @Router /admin/operations/{id} [get]
func (h AdminApisHandler) GetRewardOperation(claims *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	resData, err := h.app.Services.GetRewardType(&claims.AppID, claims.OrgID, id)
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

// UpdateRewardOperation Updates a reward operation with the specified id
// @Description Updates a reward operation with the specified id
// @Tags Admin
// @ID AdminUpdateRewardOperation
// @Param data body model.RewardOperation true "body json"
// @Accept json
// @Produce json
// @Success 200 {object} model.RewardOperation
// @Security AdminUserAuth
// @Router /admin/operations/{id} [put]
func (h AdminApisHandler) UpdateRewardOperation(claims *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {
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

	resData, err := h.app.Services.UpdateRewardType(&claims.AppID, claims.OrgID, id, item)
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

// CreateRewardOperation Create a new operation
// @Description Create a new operation type
// @Tags Admin
// @ID AdminCreateRewardOperation
// @Param data body model.RewardOperation true "body json"
// @Accept json
// @Success 200 {object} model.RewardOperation
// @Security AdminUserAuth
// @Router /admin/operations [post]
func (h AdminApisHandler) CreateRewardOperation(claims *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {

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

	createdItem, err := h.app.Services.CreateRewardType(&claims.AppID, claims.OrgID, item)
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

// DeleteRewardOperation Deletes a reward operation with the specified id
// @Description Deletes a reward operation with the specified id
// @Tags Admin
// @ID AdminDeleteRewardOperation
// @Success 200
// @Security AdminUserAuth
// @Router /admin/operations/{id} [delete]
func (h AdminApisHandler) DeleteRewardOperation(claims *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.app.Services.DeleteRewardType(&claims.AppID, claims.OrgID, id)
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
// @Param in_stock query string false "in_stock - possible values: missing (e.g no filter), 0- false, 1- true"
// @Param grant_depleted query string false "grant_depleted - possible values: missing (e.g no filter), 0- false, 1- true"
// @Param claim_depleted query string false "claim_depleted - possible values: missing (e.g no filter), 0- false, 1- true"
// @Param limit query string false "limit - limit the result"
// @Param offset query string false "offset"
// @Tags Admin
// @ID AdminGetRewardInventories
// @Success 200 {array} model.RewardInventory
// @Security AdminUserAuth
// @Router /admin/inventories [get]
func (h AdminApisHandler) GetRewardInventories(claims *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {

	rewardType := getStringQueryParam(r, "reward_type")
	inStock := getBoolQueryParam(r, "in_stock", nil)
	grantDepleted := getBoolQueryParam(r, "grant_depleted", nil)
	claimDepleted := getBoolQueryParam(r, "claim_depleted", nil)
	limitFilter := getInt64QueryParam(r, "limit")
	offsetFilter := getInt64QueryParam(r, "offset")

	IDs := []string{}
	IDskeys, ok := r.URL.Query()["ids"]
	if ok && len(IDskeys[0]) > 0 {
		extIDs := IDskeys[0]
		IDs = strings.Split(extIDs, ",")
	}

	resData, err := h.app.Services.GetRewardInventories(claims.OrgID, IDs, rewardType, inStock, grantDepleted, claimDepleted, limitFilter, offsetFilter)
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
// @Router /admin/inventories/{id} [get]
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
// @Param data body model.RewardInventory true "body json"
// @Accept json
// @Produce json
// @Success 200 {object} model.RewardInventory
// @Security AdminUserAuth
// @Router /admin/inventories/{id} [put]
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
// @Param data body model.RewardInventory true "body json"
// @Accept json
// @Success 200 {object} model.RewardInventory
// @Security AdminUserAuth
// @Router /admin/inventories [post]
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

// GetRewardClaims Retrieves  all reward claims
// @Description Retrieves  all reward claims
// @Param ids query string false "Coma separated IDs of the desired records"
// @Param user_id query string false "user_id"
// @Param status query string false "status"
// @Param limit query string false "limit - limit the result"
// @Param offset query string false "offset"
// @Tags Admin
// @ID AdminGetRewardClaims
// @Success 200 {array} model.RewardClaim
// @Security AdminUserAuth
// @Router /admin/claims [get]
func (h AdminApisHandler) GetRewardClaims(claims *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {

	rewardType := getStringQueryParam(r, "reward_type")
	userID := getStringQueryParam(r, "user_id")
	status := getStringQueryParam(r, "status")
	limitFilter := getInt64QueryParam(r, "limit")
	offsetFilter := getInt64QueryParam(r, "offset")

	IDs := []string{}
	IDskeys, ok := r.URL.Query()["ids"]
	if ok && len(IDskeys[0]) > 0 {
		extIDs := IDskeys[0]
		IDs = strings.Split(extIDs, ",")
	}

	resData, err := h.app.Services.GetRewardClaims(claims.OrgID, IDs, userID, rewardType, status, limitFilter, offsetFilter)
	if err != nil {
		log.Printf("Error on adminapis.getRewardClaims: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if resData == nil {
		resData = []model.RewardClaim{}
	}

	data, err := json.Marshal(resData)
	if err != nil {
		log.Printf("Error on adminapis.getRewardClaims: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// GetRewardClaim Retrieves a reward claim by id
// @Description Retrieves a claim inventory by id
// @Tags Admin
// @ID AdminGetRewardClaim
// @Accept json
// @Produce json
// @Success 200 {object} model.RewardClaim
// @Security AdminUserAuth
// @Router /admin/claims/{id} [get]
func (h AdminApisHandler) GetRewardClaim(claims *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	resData, err := h.app.Services.GetRewardClaim(claims.OrgID, id)
	if err != nil {
		log.Printf("Error on adminapis.getRewardClaim(%s): %s", id, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(resData)
	if err != nil {
		log.Printf("Error on adminapis.getRewardClaim(%s): %s", id, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// UpdateRewardClaim Updates a reward claim with the specified id
// @Description Updates a reward claim with the specified id
// @Tags Admin
// @ID AdminUpdateRewardClaim
// @Param data body model.RewardClaim true "body json"
// @Accept json
// @Produce json
// @Success 200 {object} model.RewardClaim
// @Security AdminUserAuth
// @Router /admin/claims/{id} [put]
func (h AdminApisHandler) UpdateRewardClaim(claims *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on adminapis.updateRewardClaim(%s): %s", id, err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var item model.RewardClaim
	err = json.Unmarshal(data, &item)
	if err != nil {
		log.Printf("Error on adminapis.updateRewardClaim(%s): %s", id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resData, err := h.app.Services.UpdateRewardClaim(claims.OrgID, id, item)
	if err != nil {
		log.Printf("Error on adminapis.updateRewardClaim(%s): %s", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(resData)
	if err != nil {
		log.Printf("Error on adminapis.updateRewardClaim(%s): %s", id, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// CreateRewardClaim Create a new claim inventory
// @Description Create a new claim inventory
// @Tags Admin
// @ID AdminCreateRewardClaim
// @Param data body model.RewardClaim true "body json"
// @Accept json
// @Success 200 {object} model.RewardInventory
// @Security AdminUserAuth
// @Router /admin/claims [post]
func (h AdminApisHandler) CreateRewardClaim(claims *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on adminapis.createRewardClaim: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var item model.RewardClaim
	err = json.Unmarshal(data, &item)
	if err != nil {
		log.Printf("Error on adminapis.createRewardClaim: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdItem, err := h.app.Services.CreateRewardClaim(claims.OrgID, item)
	if err != nil {
		log.Printf("Error on adminapis.createRewardClaim: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(createdItem)
	if err != nil {
		log.Printf("Error on adminapis.createRewardClaim: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
