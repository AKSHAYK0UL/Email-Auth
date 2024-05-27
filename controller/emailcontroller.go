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
