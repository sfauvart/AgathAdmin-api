package structs

import (
	"github.com/urfave/negroni"
)

// Route is a structure of Route
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc negroni.HandlerFunc
	Auth        bool
	Roles       []string
}

// SpiritHandler is a handler of spirits
type Routes struct {
	Routes []Route
	Prefix string
}
