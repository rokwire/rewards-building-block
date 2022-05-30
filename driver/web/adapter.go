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

package web

import (
	"fmt"
	"log"
	"net/http"
	"rewards/core"
	"rewards/core/model"
	"rewards/driver/web/rest"
	"rewards/utils"
	"strings"

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

	app *core.Application
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
	adminSubRouter.HandleFunc("/types", we.coreAuthWrapFunc(we.adminApisHandler.GetRewardTypes, nil)).Methods("GET")
	adminSubRouter.HandleFunc("/types", we.coreAuthWrapFunc(we.adminApisHandler.CreateRewardType, nil)).Methods("POST")
	adminSubRouter.HandleFunc("/types/{id}", we.coreAuthWrapFunc(we.adminApisHandler.GetRewardType, nil)).Methods("GET")
	adminSubRouter.HandleFunc("/types/{id}", we.coreAuthWrapFunc(we.adminApisHandler.UpdateRewardType, nil)).Methods("PUT")
	adminSubRouter.HandleFunc("/types/{id}", we.coreAuthWrapFunc(we.adminApisHandler.DeleteRewardType, nil)).Methods("DELETE")

	adminSubRouter.HandleFunc("/operations", we.coreAuthWrapFunc(we.adminApisHandler.GetRewardOperations, nil)).Methods("GET")
	adminSubRouter.HandleFunc("/operations", we.coreAuthWrapFunc(we.adminApisHandler.CreateRewardOperation, nil)).Methods("POST")
	adminSubRouter.HandleFunc("/operations/{id}", we.coreAuthWrapFunc(we.adminApisHandler.GetRewardOperation, nil)).Methods("GET")
	adminSubRouter.HandleFunc("/operations/{id}", we.coreAuthWrapFunc(we.adminApisHandler.UpdateRewardOperation, nil)).Methods("PUT")
	adminSubRouter.HandleFunc("/operations/{id}", we.coreAuthWrapFunc(we.adminApisHandler.DeleteRewardOperation, nil)).Methods("DELETE")

	adminSubRouter.HandleFunc("/inventories", we.coreAuthWrapFunc(we.adminApisHandler.GetRewardInventories, nil)).Methods("GET")
	adminSubRouter.HandleFunc("/inventories", we.coreAuthWrapFunc(we.adminApisHandler.CreateRewardInventory, nil)).Methods("POST")
	adminSubRouter.HandleFunc("/inventories/{id}", we.coreAuthWrapFunc(we.adminApisHandler.GetRewardInventory, nil)).Methods("GET")
	adminSubRouter.HandleFunc("/inventories/{id}", we.coreAuthWrapFunc(we.adminApisHandler.UpdateRewardInventory, nil)).Methods("PUT")

	adminSubRouter.HandleFunc("/claims", we.coreAuthWrapFunc(we.adminApisHandler.GetRewardClaims, nil)).Methods("GET")
	adminSubRouter.HandleFunc("/claims", we.coreAuthWrapFunc(we.adminApisHandler.CreateRewardClaim, nil)).Methods("POST")
	adminSubRouter.HandleFunc("/claims/{id}", we.coreAuthWrapFunc(we.adminApisHandler.GetRewardClaim, nil)).Methods("GET")
	adminSubRouter.HandleFunc("/claims/{id}", we.coreAuthWrapFunc(we.adminApisHandler.UpdateRewardClaim, nil)).Methods("PUT")

	log.Fatal(http.ListenAndServe(":"+we.port, router))
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

type coreAuthFunc = func(*tokenauth.Claims, http.ResponseWriter, *http.Request)

func (we Adapter) coreAuthWrapFunc(handler coreAuthFunc, authorization Authorization) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		utils.LogRequest(req)
		responseStatus, claims, err := authorization.check(req)
		if err != nil {
			log.Printf("error authorization check - %s", err)
			http.Error(w, http.StatusText(responseStatus), responseStatus)
			return
		}

		handler(claims, w, req)
	}
}

type apiKeysAuthFunc = func(http.ResponseWriter, *http.Request)

func (we Adapter) apiKeyOrTokenWrapFunc(handler apiKeysAuthFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		utils.LogRequest(req)

		// apply core token check
		coreAuth, _ := we.auth.coreAuth.Check(req)
		if coreAuth {
			handler(w, req)
			return
		}

		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}
}

type userAuthFunc = func(*tokenauth.Claims, http.ResponseWriter, *http.Request)

func (we Adapter) userAuthWrapFunc(handler userAuthFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		utils.LogRequest(req)

		coreAuth, claims := we.auth.coreAuth.Check(req)
		if coreAuth && claims != nil && !claims.Anonymous {
			handler(claims, w, req)
			return
		}
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}
}

type adminAuthFunc = func(*tokenauth.Claims, http.ResponseWriter, *http.Request)

func (we Adapter) adminAuthWrapFunc(handler adminAuthFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		utils.LogRequest(req)

		obj := req.URL.Path // the resource that is going to be accessed.
		act := req.Method   // the operation that the user performs on the resource.

		coreAuth, claims := we.auth.coreAuth.Check(req)
		if coreAuth {
			permissions := strings.Split(claims.Permissions, ",")

			HasAccess := false
			for _, s := range permissions {
				HasAccess = we.authorization.Enforce(s, obj, act)
				if HasAccess {
					break
				}
			}
			if HasAccess {
				handler(claims, w, req)
				return
			}
			log.Printf("Access control error - Core Subject: %s is trying to apply %s operation for %s\n", claims.Subject, act, obj)
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

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
func NewWebAdapter(host string, port string, app *core.Application, config model.Config) Adapter {
	auth := NewAuth(app, config)
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
	}
}

// AppListener implements core.ApplicationListener interface
type AppListener struct {
	adapter *Adapter
}
