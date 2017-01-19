package routers

import (
	"github.com/sfauvart/Agathadmin-api/controllers/auth"
	"github.com/sfauvart/Agathadmin-api/routers/structs"
	"net/http"
)

const (
	prefix = "/auth"
)

// NewSpiritHandler creates a new spirit handler to manage spirits
func NewAuthRoutes() *structs.Routes {
	authRoutes := structs.Routes{
		Prefix: prefix,
	}

	// build routes
	routes := []structs.Route{}
	// SignIn
	routes = append(routes, structs.Route{
		Name:        "Sign in",
		Method:      http.MethodPost,
		Pattern:     "/signin",
		HandlerFunc: auth.SignInController,
		Auth:        false,
	})
	// Refresh Token
	routes = append(routes, structs.Route{
		Name:        "Check & Refresh token",
		Method:      http.MethodPost,
		Pattern:     "/check",
		HandlerFunc: auth.RefreshTokenController,
		Auth:        false,
	})

	authRoutes.Routes = routes

	return &authRoutes
}
