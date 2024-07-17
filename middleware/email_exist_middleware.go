package middleware

import (
	"net/http"

	"github.com/AKSHAYK0UL/Email_Auth/model"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CheckEmailInDbMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var mdwEmail model.MiddlewareEmail
		if err := ctx.ShouldBindBodyWithJSON(&mdwEmail); err != nil {
			ctx.String(http.StatusBadRequest, "invalid request")
			ctx.Abort()
			return
		}
		filter := bson.D{{Key: "useremail", Value: mdwEmail.UserEmail}}
		result := model.MongoInstance.Mdatabase.Collection("Account").FindOne(ctx, filter)
		if err := result.Err(); err != nil {
			if err == mongo.ErrNoDocuments {
				ctx.String(http.StatusNotFound, "account not found (mdw)")
				ctx.Abort()
				return
			} else {
				ctx.String(http.StatusInternalServerError, "error retrieving account")
				ctx.Abort()
				return
			}
		}
		ctx.Next()
	}
}
