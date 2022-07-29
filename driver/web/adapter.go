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

package web

import (
	"fmt"
	"log"
	"net/http"
	"rewards/core"
	"rewards/core/model"
	"rewards/driver/web/rest"
	"rewards/utils"

	"github.com/rokwire/logging-library-go/logs"

	"github.com/rokwire/core-auth-library-go/tokenauth"

	"github.com/casbin/casbin"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

//Adapter entity
type Adapter struct {
	host          string
	port          string
	auth          *Auth
	authorization *casbin.Enforcer

	apisHandler         rest.ApisHandler
	adminApisHandler    rest.AdminApisHandler
	internalApisHandler rest.InternalApisHandler

	app    *core.Application
	logger *logs.Logger
}

// @title Rewards Building Block API
// @description RoRewards Building Block API Documentation.
// @version 1.0.6
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost
// @BasePath /content
// @schemes https

// @securityDefinitions.apikey InternalApiAuth
// @in header (add INTERNAL-API-KEY with correct value as a header)
// @name Authorization

// @securityDefinitions.apikey AdminUserAuth
// @in header (add Bearer prefix to the Authorization value)
// @name Authorization

// @securityDefinitions.apikey AdminGroupAuth
// @in header
// @name GROUP

//Start starts the module
func (we Adapter) Start() {

	router := mux.NewRouter().StrictSlash(true)

	subrouter := router.PathPrefix("/rewards").Subrouter()
	subrouter.PathPrefix("/doc/ui").Handler(we.serveDocUI())
	subrouter.HandleFunc("/doc", we.serveDoc)
	subrouter.HandleFunc("/version", we.wrapFunc(we.apisHandler.Version)).Methods("GET")

	// handle apis
	apiRouter := subrouter.PathPrefix("/api").Subrouter()

	// Internal APIs called from other BBs
	apiRouter.HandleFunc("/int/reward", we.internalAPIKeyAuthWrapFunc(we.internalApisHandler.CreateReward)).Methods("POST")
	apiRouter.HandleFunc("/int/stats", we.internalAPIKeyAuthWrapFunc(we.internalApisHandler.GetRewardStats)).Methods("GET")

	// Client APIs
	apiRouter.HandleFunc("/user/balance", we.userAuthWrapFunc(we.apisHandler.GetUserBalance)).Methods("GET")
	apiRouter.HandleFunc("/user/history", we.userAuthWrapFunc(we.apisHandler.GetUserRewardsHistory)).Methods("GET")
	apiRouter.HandleFunc("/user/claims", we.userAuthWrapFunc(we.apisHandler.GetUserRewardClaim)).Methods("GET")
	apiRouter.HandleFunc("/user/claims", we.userAuthWrapFunc(we.apisHandler.CreateUserRewardClaim)).Methods("POST")

	// handle student guide admin apis
	adminSubRouter := apiRouter.PathPrefix("/admin").Subrouter()
	adminSubRouter.HandleFunc("/types", we.adminAuthWrapFunc(we.adminApisHandler.GetRewardTypes)).Methods("GET")
	adminSubRouter.HandleFunc("/types", we.adminAuthWrapFunc(we.adminApisHandler.CreateRewardType)).Methods("POST")
	adminSubRouter.HandleFunc("/types/{id}", we.adminAuthWrapFunc(we.adminApisHandler.GetRewardType)).Methods("GET")
	adminSubRouter.HandleFunc("/types/{id}", we.adminAuthWrapFunc(we.adminApisHandler.UpdateRewardType)).Methods("PUT")
	adminSubRouter.HandleFunc("/types/{id}", we.adminAuthWrapFunc(we.adminApisHandler.DeleteRewardType)).Methods("DELETE")

	adminSubRouter.HandleFunc("/operations", we.adminAuthWrapFunc(we.adminApisHandler.GetRewardOperations)).Methods("GET")
	adminSubRouter.HandleFunc("/operations", we.adminAuthWrapFunc(we.adminApisHandler.CreateRewardOperation)).Methods("POST")
	adminSubRouter.HandleFunc("/operations/{id}", we.adminAuthWrapFunc(we.adminApisHandler.GetRewardOperation)).Methods("GET")
	adminSubRouter.HandleFunc("/operations/{id}", we.adminAuthWrapFunc(we.adminApisHandler.UpdateRewardOperation)).Methods("PUT")
	adminSubRouter.HandleFunc("/operations/{id}", we.adminAuthWrapFunc(we.adminApisHandler.DeleteRewardOperation)).Methods("DELETE")

	adminSubRouter.HandleFunc("/inventories", we.adminAuthWrapFunc(we.adminApisHandler.GetRewardInventories)).Methods("GET")
	adminSubRouter.HandleFunc("/inventories", we.adminAuthWrapFunc(we.adminApisHandler.CreateRewardInventory)).Methods("POST")
	adminSubRouter.HandleFunc("/inventories/{id}", we.adminAuthWrapFunc(we.adminApisHandler.GetRewardInventory)).Methods("GET")
	adminSubRouter.HandleFunc("/inventories/{id}", we.adminAuthWrapFunc(we.adminApisHandler.UpdateRewardInventory)).Methods("PUT")

	adminSubRouter.HandleFunc("/claims", we.adminAuthWrapFunc(we.adminApisHandler.GetRewardClaims)).Methods("GET")
	adminSubRouter.HandleFunc("/claims", we.adminAuthWrapFunc(we.adminApisHandler.CreateRewardClaim)).Methods("POST")
	adminSubRouter.HandleFunc("/claims/{id}", we.adminAuthWrapFunc(we.adminApisHandler.GetRewardClaim)).Methods("GET")
	adminSubRouter.HandleFunc("/claims/{id}", we.adminAuthWrapFunc(we.adminApisHandler.UpdateRewardClaim)).Methods("PUT")

	//log.Fatal(http.ListenAndServe(":"+we.port, router))
	log.Fatal(http.ListenAndServe(":81", router))
}

func (we Adapter) serveDoc(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("access-control-allow-origin", "*")
	http.ServeFile(w, r, "./docs/swagger.yaml")
}

func (we Adapter) serveDocUI() http.Handler {
	url := fmt.Sprintf("%s/rewards/doc", we.host)
	return httpSwagger.Handler(httpSwagger.URL(url))
}

func (we Adapter) wrapFunc(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		utils.LogRequest(req)

		handler(w, req)
	}
}

type apiKeysAuthFunc = func(http.ResponseWriter, *http.Request)

func (we Adapter) apiKeyOrTokenWrapFunc(handler apiKeysAuthFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		utils.LogRequest(req)

		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}
}

type userAuthFunc = func(*tokenauth.Claims, http.ResponseWriter, *http.Request)

func (we Adapter) userAuthWrapFunc(handler userAuthFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		utils.LogRequest(req)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}
}

type adminAuthFunc = func(*tokenauth.Claims, http.ResponseWriter, *http.Request)

func (we Adapter) adminAuthWrapFunc(handler adminAuthFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		utils.LogRequest(req)

		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}
}

type internalAPIKeyAuthFunc = func(http.ResponseWriter, *http.Request)

func (we Adapter) internalAPIKeyAuthWrapFunc(handler internalAPIKeyAuthFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		utils.LogRequest(req)

		apiKeyAuthenticated := we.auth.internalAuth.check(w, req)

		if apiKeyAuthenticated {
			handler(w, req)
		} else {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
	}
}

// NewWebAdapter creates new WebAdapter instance
func NewWebAdapter(host string, port string, app *core.Application, config model.Config, logger *logs.Logger) Adapter {
	auth := NewAuth(app, config, logger)
	authorization := casbin.NewEnforcer("driver/web/authorization_model.conf", "driver/web/authorization_policy.csv")

	apisHandler := rest.NewApisHandler(app)
	adminApisHandler := rest.NewAdminApisHandler(app)
	internalApisHandler := rest.NewInternalApisHandler(app)
	return Adapter{
		host:                host,
		port:                port,
		auth:                auth,
		authorization:       authorization,
		apisHandler:         apisHandler,
		adminApisHandler:    adminApisHandler,
		internalApisHandler: internalApisHandler,
		app:                 app,
		logger:              logger,
	}
}

// AppListener implements core.ApplicationListener interface
type AppListener struct {
	adapter *Adapter
}
