package dao

import (
	"github.com/sfauvart/Agathadmin-api/settings"
)

// GetUserDAO returns a UserDAO according to type and params
func GetUserDAO(daoType int) (UserDAO, error) {
	switch daoType {
	case DAOMongo:
		return NewUserDAOMongo(settings.GetMongoSession()), nil
	case DAOMock:
		return NewUserDAOMock(), nil
	default:
		return nil, ErrorDAONotFound
	}
}
