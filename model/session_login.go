package model

type SessionLogin struct {
	UserEmail string `json:"useremail,omitempty" bson:"useremail,omitempty"`
	UserName  string `json:"username,omitempty" bson:"username,omitempty"`
}
