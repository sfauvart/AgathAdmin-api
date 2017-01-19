package models

import "gopkg.in/mgo.v2/bson"

type Base interface {
	GetId() bson.ObjectId
}

type Abstract struct {
	Id bson.ObjectId `json:"id" bson:"_id,omitempty"`
}

func (m Abstract) GetId() bson.ObjectId {
	return m.Id
}
