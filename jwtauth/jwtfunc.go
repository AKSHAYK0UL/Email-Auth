package jwtauth

import (
	"errors"
	"os"

	"github.com/AKSHAYK0UL/Email_Auth/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

// JWT functions
func GenerateAuthToken(data model.UserAccount) (string, error) {
	godotenv.Load()
	key := os.Getenv("SecretKey")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":   data.UserId,
		"name":     data.UserName,
		"email":    data.UserEmail,
		"phone":    data.Phone,
		"authType": data.AuthType,
	})
	return token.SignedString([]byte(key))
}
func VerifyAuthToken(authToken string, userId string) error {
	if authToken == "" {
		return errors.New("unauthorized")
	}
	token, err := jwt.Parse(authToken, func(t *jwt.Token) (interface{}, error) {
		godotenv.Load()
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("invalid signin method!")
		}
		key := os.Getenv("SecretKey")

		return []byte(key), nil
	})
	if err != nil {
		return errors.New("could not parse token")
	}
	tokenIsValid := token.Valid //check if the token is sign with the valid sercet key or not
	if !tokenIsValid {
		return errors.New("token is not valid")

	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid token claims")

	}
	claimsuserId := claims["userId"].(string)
	if claimsuserId != userId {
		return errors.New("invalid user credentials")
	}
	return nil

}
