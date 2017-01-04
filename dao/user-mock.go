package dao

import (
  "github.com/sfauvart/Agathadmin-api/models"
	"time"
	"gopkg.in/mgo.v2/bson"
)

// MockedUser is the user returned by this mocked interface
var MockedUser = models.User {
  Abstract: models.Abstract{
    Id:     bson.NewObjectId(),
  },
	Username:         "fakeuser",
	Email:            "fake@user.com",
	Enabled:          true,
	Password:         "secret",
	LastLogin:        time.Date(2015, 01, 01, 0, 0, 0, 0, time.UTC),
	Locked:           false,
	FirstName:        "Fake",
	LastName:         "USER",
	Roles:            []string{"ADMIN","FAKEROLE","GUEST"},
}

// UserDAOMock is the mocked implementation of the UserDAO
type UserDAOMock struct {
  *AbstractDAOMock
}

// NewUserDAOMock creates a new UserDAO with a mocked implementation
func NewUserDAOMock() UserDAO {
  userDaoMock := UserDAOMock{
		&AbstractDAOMock{
			model: models.User{},
		},
	}
	return &userDaoMock
}

// GetUserByEmail returns a user by its Email
func (s *UserDAOMock) GetByEmail(Email string) (*models.User, error) {
	return &MockedUser, nil
}
