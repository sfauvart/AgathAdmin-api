package dao

import (
  "github.com/sfauvart/Agathadmin-api/models"
	"gopkg.in/mgo.v2/bson"
	logger "github.com/Sirupsen/logrus"
)

// MockedAbstract is the abstract returned by this mocked interface
var MockedAbstract = models.Abstract{
  Id:     bson.NewObjectId(),
}

// AbstractDAOMock is the mocked implementation of the AbstractDAO
type BaseDAOMock interface {
	GetModel() interface{}
}
// AbstractDAOMongo is the mongo implementation of the UserDAO
type AbstractDAOMock struct {
	model			 interface{}
}

func (a *AbstractDAOMock) GetModel() interface{} {
	return a.model
}

// NewAbstractDAOMock creates a new AbstractDAO with a mocked implementation
func NewAbstractDAOMock(s *AbstractDAOMock) BaseDAO {
  if s.model == nil {
		logger.WithField("error", s.model).Warn("no model defined !")
	}
	return s
}

// GetUserByID returns a user by its ID
func (s *AbstractDAOMock) GetByID(ID string) (interface{}, error) {
	return &MockedAbstract, nil
}

// GetAllUsers returns all users with paging capability
func (s *AbstractDAOMock) GetAll(start, end int) (interface{}, error) {
	return []models.Abstract{MockedAbstract}, nil
}

// SaveUser saves the user
func (s *AbstractDAOMock) Save(model interface{}) error {
	return nil
}

// UpsertUser updates or creates a user
func (s *AbstractDAOMock) UpsertByID(ID string, model interface{}) (bool, error) {
	return true, nil
}

// DeleteUser deletes a users by its ID
func (s *AbstractDAOMock) DeleteByID(ID string) error {
	return nil
}
