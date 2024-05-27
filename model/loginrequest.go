package model

type Login struct {
	UserEmail string `json:"useremail,omitempty" bson:"useremail,omitempty"`
	Password  string `json:"password,omitempty" bson:"password,omitempty"`
}
