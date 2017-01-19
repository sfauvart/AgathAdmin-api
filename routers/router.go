package routers

import (
	logger "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	middleJwt "github.com/sfauvart/Agathadmin-api/middlewares/jwt"
	"github.com/sfauvart/Agathadmin-api/routers/api"
	"github.com/sfauvart/Agathadmin-api/routers/structs"
	"github.com/urfave/negroni"
)

func InitRoutes() *mux.Router {
	muxRouter := mux.NewRouter().StrictSlash(false)

	allRoutes := []*structs.Routes{NewAuthRoutes(), api.NewUsersRoutes()}

	for _, routes := range allRoutes {
		for _, route := range routes.Routes {
			logger.WithField("route", route).Debug("adding route to mux")

			_n := negroni.New()
			if route.Auth == true {
				_n.Use(negroni.HandlerFunc(middleJwt.RequireTokenAuthentication))
				_n.Use(negroni.HandlerFunc(middleJwt.RequireRoles(route.Roles)))
			}
			_n.Use(negroni.HandlerFunc(route.HandlerFunc))

			muxRouter.Handle(routes.Prefix+route.Pattern, _n).Methods(route.Method).Name(route.Name)
		}
	}

	return muxRouter
}
