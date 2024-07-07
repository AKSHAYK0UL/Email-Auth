package model

type UserIdtype struct {
	UserId string `json:"userid,omitempty" bson:"_id,omitempty"`
}

type GEmailUserType struct {
	UserEmail string `json:"useremail,omitempty" bson:"useremail,omitempty"`
}
