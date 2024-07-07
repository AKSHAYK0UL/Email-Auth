package controller

import (
	"fmt"
	"net/http"

	"github.com/AKSHAYK0UL/Email_Auth/helper"
	"github.com/AKSHAYK0UL/Email_Auth/model"
	"github.com/gin-gonic/gin"
)

// Send email controller
func SendmailContoller(ctx *gin.Context) {
	var requestData model.RequestModel
	ctx.ShouldBindBodyWithJSON(&requestData)
	domainval, err := helper.IfEmailIsAllowed(requestData.UserEmail)
	if err != nil {
		// ctx.JSON(http.StatusNotAcceptable, "try another email")
		ctx.String(http.StatusNotAcceptable, "try another email")

	} else {

		response, err := helper.SendEmail(domainval, requestData)
		if err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
		} else {
			ctx.JSON(http.StatusAccepted, response)

		}

	}
}
func VerificationController(ctx *gin.Context) {
	var requestData model.SignUpResponse
	ctx.ShouldBindBodyWithJSON(&requestData)
	response, err := helper.VerifyCode(requestData.UserId, requestData.Vcode)
	if err != nil {
		fmt.Println(err)
		ctx.String(http.StatusBadRequest, "invalid")
		// ctx.JSON(http.StatusBadRequest, "invalid")
	} else {
		ctx.JSON(http.StatusAccepted, response)

	}

}
func ResetPasswordController(ctx *gin.Context) {
	var userReset model.RequestModel
	ctx.ShouldBindBodyWithJSON(&userReset)
	domainval, err := helper.IfEmailIsAllowed(userReset.UserEmail)
	if err != nil {
		// ctx.JSON(http.StatusNotAcceptable, "try another email")
		ctx.String(http.StatusNotAcceptable, "try another email")

	} else {
		response, err := helper.ResetpasswordSendEmail(domainval, userReset)
		if err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
		} else {
			ctx.JSON(http.StatusAccepted, response)

		}
	}
}
func ResetverifyController(ctx *gin.Context) {
	var requestData model.SignUpResponse
	ctx.ShouldBindBodyWithJSON(&requestData)
	response, err := helper.Resetverify(requestData.UserId, requestData.Vcode)
	if err != nil {
		fmt.Println(err)
		ctx.String(http.StatusBadRequest, "invalid")
		// ctx.JSON(http.StatusBadRequest, "invalid")
	} else {
		ctx.JSON(http.StatusAccepted, response)

	}
}
func LoginController(ctx *gin.Context) {
	var Loginreqdata model.Login
	ctx.ShouldBindBodyWithJSON(&Loginreqdata)
	_, err := helper.IfEmailIsAllowed(Loginreqdata.UserEmail)
	if err != nil {
		// ctx.JSON(http.StatusNotAcceptable, "try another email")
		ctx.String(http.StatusNotAcceptable, "try another email")

	} else {
		response, err := helper.LoginToAccount(Loginreqdata.UserEmail, Loginreqdata.Password)
		if err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
		} else {
			ctx.JSON(http.StatusAccepted, response)

		}
	}

}
func UserExistController(ctx *gin.Context) {
	var userexist model.UserIdtype
	if err := ctx.ShouldBindBodyWithJSON(&userexist); err != nil {
		ctx.String(http.StatusNotFound, "User EROROR found")
		return

	}
	// ctx.ShouldBindBodyWithJSON(&userexist.UserId)
	response, err := helper.UserExist(userexist.UserId)
	if err != nil {
		ctx.String(http.StatusNotFound, err.Error())
		return

	}

	ctx.JSON(http.StatusFound, response)

}

func SaveGoogleUserController(ctx *gin.Context) {
	var G_user model.UserAccount
	if err := ctx.ShouldBindBodyWithJSON(&G_user); err != nil {
		ctx.String(http.StatusNotFound, "error in binding")
		return
	}
	response, err := helper.SaveGUser(G_user)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return

	} else {
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
}

func GoogleUserExistController(ctx *gin.Context) {
	var userexist model.GEmailUserType
	if err := ctx.ShouldBindBodyWithJSON(&userexist); err != nil {
		ctx.String(http.StatusNotFound, "User EROROR found")
		return

	}
	// ctx.ShouldBindBodyWithJSON(&userexist.UserId)
	response, err := helper.GoogleUserExist(userexist.UserEmail)
	if err != nil {
		ctx.String(http.StatusNotFound, err.Error())
		return

	}

	ctx.JSON(http.StatusFound, response)

}

// session SignUp to register the use info[name,email,uid]
func SessionSignupController(ctx *gin.Context) {
	var requestData model.RequestModel
	ctx.ShouldBindBodyWithJSON(&requestData)
	domainval, err := helper.IfEmailIsAllowed(requestData.UserEmail)
	if err != nil {
		// ctx.JSON(http.StatusNotAcceptable, "try another email")
		ctx.String(http.StatusNotAcceptable, "try another email")

	} else {

		response, err := helper.SendEmail(domainval, requestData)
		if err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
		} else {
			ctx.JSON(http.StatusAccepted, response)

		}

	}
}
