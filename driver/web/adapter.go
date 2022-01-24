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

	apisHandler      rest.ApisHandler
	adminApisHandler rest.AdminApisHandler

	app *core.Application
}

// @title Rewards Building Block API
// @description RoRewards Building Block API Documentation.
// @version 0.0.1
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost
// @BasePath /content
// @schemes https

// @securityDefinitions.apikey UserAuth
// @in header (add Bearer prefix to the Authorization value)
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

	// handle apis
	contentRouter := router.PathPrefix("/rewards/api").Subrouter()
	contentRouter.PathPrefix("/doc/ui").Handler(we.serveDocUI())
	contentRouter.HandleFunc("/doc", we.serveDoc)
	contentRouter.HandleFunc("/version", we.wrapFunc(we.apisHandler.Version)).Methods("GET")

	// handle student guide admin apis
	adminSubRouter := contentRouter.PathPrefix("/admin").Subrouter()
	adminSubRouter.HandleFunc("/reward_types", we.adminAuthWrapFunc(we.adminApisHandler.GetRewardTypes)).Methods("GET")
	adminSubRouter.HandleFunc("/reward_types", we.adminAuthWrapFunc(we.adminApisHandler.CreateRewardType)).Methods("POST")
	adminSubRouter.HandleFunc("/reward_types/{id}", we.adminAuthWrapFunc(we.adminApisHandler.GetRewardType)).Methods("GET")
	adminSubRouter.HandleFunc("/reward_types/{id}", we.adminAuthWrapFunc(we.adminApisHandler.UpdateRewardType)).Methods("PUT")
	adminSubRouter.HandleFunc("/reward_types/{id}", we.adminAuthWrapFunc(we.adminApisHandler.DeleteRewardType)).Methods("DELETE")

	adminSubRouter.HandleFunc("/reward_pools", we.adminAuthWrapFunc(we.adminApisHandler.GetRewardPools)).Methods("GET")
	adminSubRouter.HandleFunc("/reward_pools", we.adminAuthWrapFunc(we.adminApisHandler.CreateRewardPool)).Methods("POST")
	adminSubRouter.HandleFunc("/reward_pools/{id}", we.adminAuthWrapFunc(we.adminApisHandler.GetRewardPool)).Methods("GET")
	adminSubRouter.HandleFunc("/reward_pools/{id}", we.adminAuthWrapFunc(we.adminApisHandler.UpdateRewardPool)).Methods("PUT")
	adminSubRouter.HandleFunc("/reward_pools/{id}", we.adminAuthWrapFunc(we.adminApisHandler.DeleteRewardPool)).Methods("DELETE")

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

type userAuthFunc = func(http.ResponseWriter, *http.Request)

func (we Adapter) userAuthWrapFunc(handler userAuthFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		utils.LogRequest(req)

		coreAuth, claims := we.auth.coreAuth.Check(req)
		if coreAuth && claims != nil && !claims.Anonymous {
			handler(w, req)
			return
		}
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}
}

type adminAuthFunc = func(http.ResponseWriter, *http.Request)

func (we Adapter) adminAuthWrapFunc(handler adminAuthFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		utils.LogRequest(req)

		handler(w, req)
		return

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
				handler(w, req)
				return
			}
			log.Printf("Access control error - Core Subject: %s is trying to apply %s operation for %s\n", claims.Subject, act, obj)
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}
}

// NewWebAdapter creates new WebAdapter instance
func NewWebAdapter(host string, port string, app *core.Application, config model.Config) Adapter {
	auth := NewAuth(app, config)
	authorization := casbin.NewEnforcer("driver/web/authorization_model.conf", "driver/web/authorization_policy.csv")

	apisHandler := rest.NewApisHandler(app)
	adminApisHandler := rest.NewAdminApisHandler(app)
	return Adapter{host: host, port: port, auth: auth, authorization: authorization, apisHandler: apisHandler, adminApisHandler: adminApisHandler, app: app}
}

// AppListener implements core.ApplicationListener interface
type AppListener struct {
	adapter *Adapter
}
