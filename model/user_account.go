package model

import "time"

type UserAccount struct {
	UserId    string `json:"userid,omitempty" bson:"_id,omitempty"`
	UserName  string `json:"username,omitempty" bson:"username,omitempty"`
	UserEmail string `json:"useremail,omitempty" bson:"useremail,omitempty"`
	AuthType  string `json:"authtype,omitempty" bson:"authtype,omitempty"`
	Phone     string `json:"phone,omitempty" bson:"phone,omitempty"`
	AuthToken string `json:"authtoken,omitempty" bson:"authtoken,omitempty"`
}
type UserAccountStoreDb struct {
	UserId    string    `json:"userid,omitempty" bson:"_id,omitempty"`
	UserName  string    `json:"username,omitempty" bson:"username,omitempty"`
	UserEmail string    `json:"useremail,omitempty" bson:"useremail,omitempty"`
	Phone     string    `json:"phone,omitempty" bson:"phone,omitempty"`
	AuthType  string    `json:"authtype,omitempty" bson:"authtype,omitempty"`
	Password  string    `json:"password,omitempty" bson:"password,omitempty"`
	LoginAt   time.Time `json:"login_at,omitempty" bson:"login_at,omitempty"`
	SignoutAt time.Time `json:"signout_at,omitempty" bson:"signout_at,omitempty"`
	CreateAt  time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdateAt  time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
