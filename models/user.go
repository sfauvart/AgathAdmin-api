package models

import "time"

type User struct {
	Abstract  `bson:",inline"`
	Username  string    `json:"username" bson:"username"`
	Email     string    `json:"email" bson:"email"`
	Enabled   bool      `json:"enabled" bson:"enabled"`
	Password  string    `json:"password" bson:"password"`
	LastLogin *time.Time `json:"last_login,omitempty" bson:"last_login,omitempty"`
	Locked    bool      `json:"locked" bson:"locked"`
	FirstName string    `json:"first_name" bson:"first_name"`
	LastName  string    `json:"last_name" bson:"last_name"`
	Roles     []string  `json:"roles" bson:"roles"`
}
