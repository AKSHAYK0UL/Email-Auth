package routefunc

import (
	"github.com/AKSHAYK0UL/Email_Auth/controller"
	"github.com/AKSHAYK0UL/Email_Auth/middleware"
	"github.com/gin-gonic/gin"
)

func RouteTable() *gin.Engine {
	route := gin.Default()
	route.POST("/signup", controller.SendmailContoller)
	route.POST("/verify", controller.VerificationController)
	route.POST("/saveguser", controller.SaveGoogleUserController)
	route.POST("/reset", middleware.CheckEmailInDbMiddleware(), controller.ResetPasswordController)
	route.PATCH("/rverify", controller.ResetverifyController)
	route.POST("/login", middleware.CheckEmailInDbMiddleware(), controller.LoginController)
	route.POST("/uexist", controller.UserExistController)
	route.POST("/guexist", controller.GoogleUserExistController)
	route.POST("/sessignup", controller.SessionSignupController)

	return route

}
