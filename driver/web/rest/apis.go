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

package rest

import (
	"encoding/json"
	"github.com/rokwire/core-auth-library-go/tokenauth"
	"io/ioutil"
	"log"
	"net/http"
	"rewards/core"
	"rewards/core/model"
)

const maxUploadSize = 15 * 1024 * 1024 // 15 mb

//ApisHandler handles the rest APIs implementation
type ApisHandler struct {
	app *core.Application
}

//Version gives the service version
// @Description Gives the service version.
// @Tags Client
// @ID Version
// @Produce plain
// @Success 200
// @Router /version [get]
func (h ApisHandler) Version(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(h.app.Services.GetVersion()))
}

// NewApisHandler creates new rest Handler instance
func NewApisHandler(app *core.Application) ApisHandler {
	return ApisHandler{app: app}
}

// NewAdminApisHandler creates new rest Handler instance
func NewAdminApisHandler(app *core.Application) AdminApisHandler {
	return AdminApisHandler{app: app}
}

// NewInternalApisHandler creates new rest Handler instance
func NewInternalApisHandler(app *core.Application) InternalApisHandler {
	return InternalApisHandler{app: app}
}

// GetUserBalance Retrieves balance for each user's wallet
// @Description Retrieves balance for each user's wallet
// @Tags Client
// @ID GetUserBalance
// @Success 200
// @Security UserAuth
// @Router /user/balance [get]
func (h *ApisHandler) GetUserBalance(userClaims *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {
	resData, err := h.app.Services.GetUserBalance(userClaims.OrgID, userClaims.Subject)
	if err != nil {
		log.Printf("Error on apis.GetUserRewardsAmount(%s): %s", userClaims.Subject, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if resData == nil {
		resData = []model.RewardTypeAmount{}
	}

	data, err := json.Marshal(resData)
	if err != nil {
		log.Printf("Error on apis.GetUserRewardsAmount(%s): %s", userClaims.Subject, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// GetUserRewardsHistory Retrieves the wallet history
// @Description Retrieves the wallet history
// @Tags Client
// @ID GetUserRewardsHistory
// @Param reward_type query string false "reward_type - filter by reward_type"
// @Param code  query string false "code - filter by code"
// @Param building_block query string false "filter by building_block"
// @Param limit query integer false "limit - limit the result"
// @Param offset query integer false "offset"
// @Success 200
// @Security UserAuth
// @Router /user/history [get]
func (h *ApisHandler) GetUserRewardsHistory(userClaims *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {
	rewardType := getStringQueryParam(r, "reward_type")
	code := getStringQueryParam(r, "code")
	buildingBlock := getStringQueryParam(r, "building_block")
	limitFilter := getInt64QueryParam(r, "limit")
	offsetFilter := getInt64QueryParam(r, "offset")

	resData, err := h.app.Services.GetUserRewardsHistory(userClaims.OrgID, userClaims.Subject, rewardType, code, buildingBlock, limitFilter, offsetFilter)
	if err != nil {
		log.Printf("Error on apis.getUserRewardsHistory(%s): %s", userClaims.Subject, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(resData)
	if err != nil {
		log.Printf("Error on apis.getUserRewardsHistory(%s): %s", userClaims.Subject, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// CreateUserRewardClaim Create a new user claim
// @Description Create a new claim user claim
// @Tags Client
// @ID CreateUserRewardClaim
// @Accept json
// @Success 200 {object} model.RewardClaim
// @Security AdminUserAuth
// @Router /user/claim [post]
func (h ApisHandler) CreateUserRewardClaim(userClaims *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on apis.CreateUserRewardClaim: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var item model.RewardClaim
	err = json.Unmarshal(data, &item)
	if err != nil {
		log.Printf("Error on apis.CreateUserRewardClaim: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	item.UserID = userClaims.Subject
	createdItem, err := h.app.Services.CreateRewardClaim(userClaims.OrgID, item)
	if err != nil {
		log.Printf("Error on apis.CreateUserRewardClaim: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(createdItem)
	if err != nil {
		log.Printf("Error on apis.CreateUserRewardClaim: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
