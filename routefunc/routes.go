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
	route.POST("/reset", middleware.CheckEmailInDbMiddleware(), controller.ResetPasswordController)
	route.POST("/rverify", controller.ResetverifyController)
	route.POST("/login", middleware.CheckEmailInDbMiddleware(), controller.LoginController)
	return route

}
