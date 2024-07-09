package model

type SignUpResponse struct {
	UserId    string `json:"userid,omitempty" bson:"_id,omitempty"`
	UserName  string `json:"username,omitempty" bson:"username,omitempty"`
	UserEmail string `json:"useremail,omitempty" bson:"useremail,omitempty"`
	Phone     string `json:"phone,omitempty" bson:"phone,omitempty"`
	Vcode     string `json:"vcode,omitempty" bson:"vcode,omitempty"`
	Password  string `json:"password,omitempty" bson:"password,omitempty"`
	Status    string `json:"status,omitempty" bson:"status,omitempty"`
}
