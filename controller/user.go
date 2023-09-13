package controller

import (
	"errors"
	"fmt"

	"web_app/dao/mysql"
	"web_app/logic"
	"web_app/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func SignUpHandler(c *gin.Context) {
	//1. deal with the parameters
	p := new(models.ParamSignUp)

	//the request parameters Error (Formatting Error)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponsError(c, CodeInvalidParam)
			return
		}
		ResponsErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}

	//Deprecates: Estimate the param with business logic (Manually)
	//if len(p.Password) == 0 || len(p.RePassword) == 0 || len(p.Username) == 0 || p.Password != p.RePassword {
	//	zap.L().Error("SignUp with invalid params")
	//	c.JSON(http.StatusOK, gin.H{
	//		"msg": "logic of request parameters Error",
	//	})
	//	return
	//}

	//2. deal with the business logic
	if err := logic.SignUp(p); err != nil {
		zap.L().Error("logic.SignUp failed", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserExist) {
			ResponsError(c, CodeUserExisted)
			return
		}
		ResponsError(c, CodeServerBusy)
		return
	}

	//3. return the response
	ResponsSuccess(c, nil)
}

func LoginHandler(c *gin.Context) {
	//1. get the request params and validate
	var p models.ParamLogin
	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Error("Login with invalid params", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponsError(c, CodeInvalidParam)
			return
		}
		ResponsErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}

	//2. deal with the business logic
	user, err := logic.Login(&p)
	if err != nil {
		zap.L().Error("logic.Login failed", zap.String("username", p.Username), zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponsError(c, CodeUserNotExisted)
			return
		}
		ResponsError(c, CodeInvalidPassword)
		return
	}

	//3. return the response
	ResponsSuccess(c, gin.H{
		"user_id":   fmt.Sprintf("%d", user.UserID), //if the id > 1<<53 - 1(the number limit in json), the num will distort
		"user_name": user.Username,
		"token":     user.Token,
	})
}
