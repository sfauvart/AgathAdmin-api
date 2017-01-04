package api

import (
  "github.com/sfauvart/Agathadmin-api/routers/structs"
  "github.com/sfauvart/Agathadmin-api/controllers/users"
	"net/http"
)

const (
	prefix = "/api/users"
)

func NewUsersRoutes() *structs.Routes {
	authRoutes := structs.Routes{
		Prefix:    prefix,
	}

	// build routes
	routes := []structs.Route{}
	// GetAll
	routes = append(routes, structs.Route{
		Name:        "Get All",
		Method:      http.MethodGet,
		Pattern:     "",
		HandlerFunc: users.GetAll,
    Auth:      true,
    Roles:    []string{"ADMIN"},
	})
  // Get
	routes = append(routes, structs.Route{
		Name:        "Get one user",
		Method:      http.MethodGet,
		Pattern:     "/{id}",
		HandlerFunc: users.Get,
    Auth:      true,
    Roles:    []string{"ADMIN"},
	})
	// Create
	routes = append(routes, structs.Route{
		Name:        "Create a user",
		Method:      http.MethodPost,
		Pattern:     "",
		HandlerFunc: users.Create,
    Auth:      true,
    Roles:    []string{"ADMIN"},
	})
	// Update
	routes = append(routes, structs.Route{
		Name:        "Update a user",
		Method:      http.MethodPut,
		Pattern:     "/{id}",
		HandlerFunc: users.Update,
    Auth:      true,
    Roles:    []string{"ADMIN"},
	})
	// Delete
	routes = append(routes, structs.Route{
		Name:        "Delete a user",
		Method:      http.MethodDelete,
		Pattern:     "/{id}",
		HandlerFunc: users.Delete,
    Auth:      true,
    Roles:    []string{"ADMIN"},
	})

	authRoutes.Routes = routes

	return &authRoutes
}
