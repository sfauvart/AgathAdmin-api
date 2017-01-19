package dao

import (
	"github.com/sfauvart/Agathadmin-api/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// AbstractDAOMongo is the mongo implementation of the UserDAO
type UserDAOMongo struct {
	*AbstractDAOMongo
}

// NewDAOMongo creates a new UserDAO mongo implementation
func NewUserDAOMongo(session *mgo.Session) UserDAO {
	userDaoMongo := UserDAOMongo{
		&AbstractDAOMongo{
			session:    session,
			collection: "users",
			index:      []string{"id", "email"},
			model:      models.User{},
		},
	}
	return &userDaoMongo
}

func (s *UserDAOMongo) GetByEmail(Email string) (*models.User, error) {
	session := s.session.Copy()
	defer session.Close()

	user := models.User{}
	c := session.DB("").C(s.collection)

	err := c.Find(bson.M{"email": Email}).One(&user)
	return &user, err
}
