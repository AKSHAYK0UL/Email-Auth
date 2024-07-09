package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RequestModel struct {
	UserId    primitive.ObjectID `json:"userid,omitempty" bson:"_id,omitempty"`
	UserName  string             `json:"username,omitempty" bson:"username,omitempty"`
	UserEmail string             `json:"useremail,omitempty" bson:"useremail,omitempty"`
	Phone     string             `json:"phone,omitempty" bson:"phone,omitempty"`
	Password  string             `json:"password,omitempty" bson:"password,omitempty"`
	Vcode     string             `json:"vcode,omitempty" bson:"vcode,omitempty"`
	SendAt    string             `json:"send_at,omitempty" bson:"send_at,omitempty"`
}
