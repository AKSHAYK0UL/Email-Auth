package model

import "time"

type UserAccount struct {
	UserId    string `json:"userid,omitempty" bson:"_id,omitempty"`
	UserName  string `json:"username,omitempty" bson:"username,omitempty"`
	UserEmail string `json:"useremail,omitempty" bson:"useremail,omitempty"`
}
type UserAccountStoreDb struct {
	UserId    string    `json:"userid,omitempty" bson:"_id,omitempty"`
	UserName  string    `json:"username,omitempty" bson:"username,omitempty"`
	UserEmail string    `json:"useremail,omitempty" bson:"useremail,omitempty"`
	Password  string    `json:"password,omitempty" bson:"password,omitempty"`
	CreateAt  time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdateAt  time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
