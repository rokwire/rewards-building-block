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
	resData, err := h.app.Services.GetUserBalance(userClaims.Subject)
	if err != nil {
		log.Printf("Error on apis.GetUserBalance(%s): %s", userClaims.Subject, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if resData == nil{
		resData = &model.WalletBalance{
			Amount: 0,
		}
	}

	data, err := json.Marshal(resData)
	if err != nil {
		log.Printf("Error on apis.GetUserBalance(%s): %s", userClaims.Subject, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// GetUserHistory Retrieves the wallet history
// @Description Retrieves the wallet history
// @Tags Client
// @ID GetUserHistory
// @Success 200
// @Security UserAuth
// @Router /user/history [get]
func (h *ApisHandler) GetUserHistory(userClaims *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {
	resData, err := h.app.Services.GetWalletHistoryEntries(userClaims.Subject)
	if err != nil {
		log.Printf("Error on apis.GetUserHistory(%s): %s", userClaims.Subject, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(resData)
	if err != nil {
		log.Printf("Error on apis.GetUserHistory(%s): %s", userClaims.Subject, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

/*
// GetWalletBalance Retrieves  the wallet balance
// @Description Retrieves  the user balance
// @Tags Client
// @ID GetWalletBalance
// @Success 200
// @Security UserAuth
// @Router /wallet/{code}/balance [get]
func (h *ApisHandler) GetWalletBalance(userClaims *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	code := vars["code"]

	if len(code) == 0 {
		log.Printf("Error on apis.getWalletBalance(%s): missing code param", userClaims.Subject)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	resData, err := h.app.Services.GetWalletBalance(userClaims.Subject, code)
	if err != nil {
		log.Printf("Error on apis.getWalletBalance(%s): %s", userClaims.Subject, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(resData)
	if err != nil {
		log.Printf("Error on apis.getWalletBalance(%s): %s", userClaims.Subject, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// GetWalletHistory Retrieves the user history
// @Description Retrieves the user history
// @Tags Client
// @ID GetWalletHistory
// @Success 200
// @Security UserAuth
// @Router /wallet/{code}/history [get]
func (h *ApisHandler) GetWalletHistory(userClaims *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	code := vars["code"]

	if len(code) == 0 {
		log.Printf("Error on apis.GetWalletHistory(%s): missing code param", userClaims.Subject)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	resData, err := h.app.Services.GetWalletHistoryEntries(userClaims.Subject)
	if err != nil {
		log.Printf("Error on apis.GetWalletHistory(%s): %s", userClaims.Subject, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(resData)
	if err != nil {
		log.Printf("Error on apis.GetWalletHistory(%s): %s", userClaims.Subject, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}*/
