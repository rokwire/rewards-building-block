package web

import (
	"github.com/rokmetro/auth-library/authservice"
	"github.com/rokmetro/auth-library/tokenauth"
	"log"
	"net/http"
	"rewards/core"
	"rewards/core/model"
)

// CoreAuth implementation
type CoreAuth struct {
	app       *core.Application
	tokenAuth *tokenauth.TokenAuth
}

// NewCoreAuth creates new CoreAuth
func NewCoreAuth(app *core.Application, config model.Config) *CoreAuth {

	serviceLoader := authservice.NewRemoteServiceRegLoader(config.CoreServiceRegLoaderURL, []string{"core"})
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
