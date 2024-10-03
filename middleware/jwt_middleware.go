package middleware

import (
	"net/http"

	"github.com/AKSHAYK0UL/Email_Auth/jwtauth"
	"github.com/AKSHAYK0UL/Email_Auth/model"
	"github.com/gin-gonic/gin"
)

func JwtMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var userexist model.UserIdtype
		//Get the JWT from the "Authorization" header
		token := ctx.Request.Header.Get("Authorization")
		if err := ctx.ShouldBindBodyWithJSON(&userexist); err != nil {
			ctx.String(http.StatusNotFound, "User EROROR found")
			return

		}

		err := jwtauth.VerifyAuthToken(token, userexist.UserId)
		if err != nil {
			ctx.String(http.StatusUnauthorized, err.Error())
			ctx.Abort()
		}
		ctx.Next()
	}

}
