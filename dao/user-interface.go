package dao

import (
	"github.com/sfauvart/Agathadmin-api/models"
)

type UserDAO interface {
	BaseDAO
	GetByEmail(Email string) (*models.User, error)
}
