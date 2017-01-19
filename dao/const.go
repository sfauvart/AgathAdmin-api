package dao

import (
	"errors"
)

const (
	// DAOMongo is used for Mongo implementation of UserDAO
	DAOMongo int = iota
	// DAOMock is used for mocked implementation of UserDAO
	DAOMock

	NoPaging               = -1
	MaxElementsPerPage     = 100
	DefaultElementsPerPage = 10
)

var (
	// ErrorDAONotFound is used for unknown DAO type
	ErrorDAONotFound = errors.New("Unknown DAO type")
)
