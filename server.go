package main

import (
	"github.com/sfauvart/Agathadmin-api/routers"
	"github.com/sfauvart/Agathadmin-api/settings"
	"github.com/urfave/negroni"
	"github.com/rs/cors"
	"net/http"
	"log"
	"fmt"
)

// Variables to identiy the build
var (
	Version string
	Build   string
)

func main() {
	fmt.Println("Webserver Version:", Version, "- Git commit hash:", Build)

	settings.Init()
	router := routers.InitRoutes()

	n := negroni.Classic()

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowCredentials: true,
		AllowedHeaders: []string{"Authorization", "Content-Type"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		Debug: settings.Get().DebugMode,
	})

	n.Use(c)
	n.UseHandler(router)
	log.Fatal(http.ListenAndServe(":3000", n))
}
