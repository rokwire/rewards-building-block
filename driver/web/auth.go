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

	"github.com/rokwire/core-auth-library-go/authservice"
	"github.com/rokwire/core-auth-library-go/tokenauth"

	"github.com/rokwire/logging-library-go/logs"
)

// CoreAuth implementation
type CoreAuth struct {
	app       *core.Application
	tokenAuth *tokenauth.TokenAuth
}

// NewCoreAuth creates new CoreAuth
func NewCoreAuth(app *core.Application, config model.Config) *CoreAuth {

	remoteConfig := authservice.RemoteAuthDataLoaderConfig{
		AuthServicesHost: config.CoreBBHost,
	}

	serviceLoader, err := authservice.NewRemoteAuthDataLoader(remoteConfig, []string{"core"}, logs.NewLogger("groupsbb", &logs.LoggerOpts{}))
	authService, err := authservice.NewAuthService("rewards", config.ContentServiceURL, serviceLoader)
	if err != nil {
		log.Fatalf("Error initializing auth service: %v", err)
	}
	tokenAuth, err := tokenauth.NewTokenAuth(true, authService, nil, nil)
	if err != nil {
		log.Fatalf("Error intitializing token auth: %v", err)
	}

	auth := CoreAuth{app: app, tokenAuth: tokenAuth}
	return &auth
}

// Check checks the request contains a valid Core access token
func (ca CoreAuth) Check(r *http.Request) (bool, *tokenauth.Claims) {
	claims, err := ca.tokenAuth.CheckRequestTokens(r)
	if err != nil {
		log.Printf("error validate token: %s", err)
		return false, nil
	}

	if claims != nil {
		if claims.Valid() == nil {
			return true, claims
		}
	}

	return false, nil
}

// Auth handler
type Auth struct {
	internalAuth *InternalAuth
	coreAuth     *CoreAuth
	logger       *logs.Logger
}

func (auth *Auth) clientIDCheck(w http.ResponseWriter, r *http.Request) bool {
	clientID := r.Header.Get("APP")
	if len(clientID) == 0 {
		clientID = "edu.illinois.rokwire"
	}

	log.Println(fmt.Sprintf("400 - Bad Request"))
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("Bad Request"))
	return false
}

// NewAuth creates new auth handler
func NewAuth(app *core.Application, config model.Config, logger *logs.Logger) *Auth {

	internalAuth := newInternalAuth(config)
	auth := Auth{internalAuth: internalAuth, logger: logger}
	return &auth
}

// InternalAuth handling the internal calls fromother BBs
type InternalAuth struct {
	internalAPIKey string
}

func newInternalAuth(config model.Config) *InternalAuth {
	auth := InternalAuth{internalAPIKey: config.InternalAPIKey}
	return &auth
}

func (auth *InternalAuth) check(w http.ResponseWriter, r *http.Request) bool {
	apiKey := r.Header.Get("INTERNAL-API-KEY")
	//check if there is api key in the header
	if len(apiKey) == 0 {
		//no key, so return 400
		log.Println(fmt.Sprintf("400 - Bad Request"))

		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
		return false
	}

	exist := auth.internalAPIKey == apiKey

	if !exist {
		//not exist, so return 401
		log.Println(fmt.Sprintf("401 - Unauthorized for key %s", apiKey))

		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
		return false
	}
	return true
}
