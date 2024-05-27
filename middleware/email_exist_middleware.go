package middleware

import (
	"net/http"

	"github.com/AKSHAYK0UL/Email_Auth/model"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func CheckEmailInDbMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var mdwEmail model.MiddlewareEmail
		if err := ctx.ShouldBindBodyWithJSON(&mdwEmail); err != nil {
			ctx.String(http.StatusBadRequest, "invalid request")
			ctx.Abort()
			return
		}
		filter := bson.D{{}}
		curr, err := model.MongoInstance.Mdatabase.Collection("Account").Find(ctx, filter)
		if err != nil {
			ctx.String(http.StatusBadRequest, "invalid request")
			ctx.Abort()
			return
		}
		for curr.Next(ctx) {
			var email model.MiddlewareEmail
			err := curr.Decode(&email)
			if err != nil {
				ctx.String(http.StatusBadRequest, "invalid request")

			} else if email.UserEmail == mdwEmail.UserEmail {

				defer curr.Close(ctx)

				ctx.Next()
				return

			}

		}

		ctx.String(http.StatusNotFound, "no user found(middleware)")
		ctx.Abort()

		defer curr.Close(ctx)

	}
}
