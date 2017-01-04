package dao

import (
	"errors"
	"reflect"
	logger "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type BaseDAOMongo interface {
	GetSession() *mgo.Session
	GetCollection() string
	GetIndex() []string
	GetModel() interface{}
}
// AbstractDAOMongo is the mongo implementation of the UserDAO
type AbstractDAOMongo struct {
	session *mgo.Session
	collection string
	index      []string
	model			 interface{}
}

func (a *AbstractDAOMongo) GetSession() *mgo.Session {
	return a.session
}
func (a *AbstractDAOMongo) GetCollection() string {
	return a.collection
}
func (a *AbstractDAOMongo) GetIndex() []string {
	return a.index
}
func (a *AbstractDAOMongo) GetModel() interface{} {
	return a.model
}

// NewDAOMongo creates a new UserDAO mongo implementation
func NewAbstractDAOMongo(s *AbstractDAOMongo) BaseDAO {
	// create index
	err := s.GetSession().DB("").C(s.GetCollection()).EnsureIndex(mgo.Index{
		Key:        s.GetIndex(),
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	})

	if err != nil {
		logger.WithField("error", err).Warn("mongo db connection")
	}

	return s
}

// GetUserByID returns a user by its ID
func (s *AbstractDAOMongo) GetByID(ID string) (interface{}, error) {
	// check ID
	if !bson.IsObjectIdHex(ID) {
		return nil, errors.New("Invalid input to ObjectIdHex")
	}

	session := s.session.Copy()
	defer session.Close()

	m := reflect.New(reflect.TypeOf(s.GetModel()))
	c := session.DB("").C(s.collection)
	err := c.Find(bson.M{"_id": bson.ObjectIdHex(ID)}).One(&m)
	return &m, err
}

// getAllUsersByQuery returns users by query and paging capability
func (s *AbstractDAOMongo) getAllByQuery(query interface{}, start, end int) (interface{}, error) {
	session := s.session.Copy()
	defer session.Close()
	c := session.DB("").C(s.collection)

	// check param
	hasPaging := start > NoPaging && end > NoPaging && end > start

	// perform request
	var err error
	slice := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(s.GetModel())), 0, 0) //[]models.Abstract{}
	// Create a pointer to a slice value and set it to the slice
	x := reflect.New(slice.Type())
	x.Elem().Set(slice)

	if hasPaging {
		err = c.Find(query).Skip(start).Limit(end - start).All(x.Interface())
	} else {
		err = c.Find(query).All(x.Interface())
	}

	return x.Interface(), err
}

// GetAllUsers returns all users with paging capability
func (s *AbstractDAOMongo) GetAll(start, end int) (interface{}, error) {
	return s.getAllByQuery(nil, start, end)
}

// SaveUser saves the user
func (s *AbstractDAOMongo) Save(model interface{}) error {
	session := s.session.Copy()
	defer session.Close()
	c := session.DB("").C(s.collection)
	return c.Insert(model)
}

// UpsertUser updates or creates a user
func (s *AbstractDAOMongo) UpsertByID(ID string, model interface{}) (bool, error) {

	// check ID
	if !bson.IsObjectIdHex(ID) {
		return false, errors.New("Invalid input to ObjectIdHex")
	}

	session := s.session.Copy()
	defer session.Close()
	c := session.DB("").C(s.collection)
	chg, err := c.Upsert(bson.M{"_id": bson.ObjectIdHex(ID)}, model)
	if err != nil {
		return false, err
	}
	return chg.Updated > 0, err
}

// DeleteUser deletes a users by its ID
func (s *AbstractDAOMongo) DeleteByID(ID string) error {

	// check ID
	if !bson.IsObjectIdHex(ID) {
		return errors.New("Invalid input to ObjectIdHex")
	}

	session := s.session.Copy()
	defer session.Close()
	c := session.DB("").C(s.collection)
	err := c.Remove(bson.M{"_id": bson.ObjectIdHex(ID)})
	return err
}
