package logic

import (
	"web_app/dao/mysql"
	"web_app/models"
	"web_app/pkg/jwt"
	"web_app/pkg/snowflake"
)

func SignUp(p *models.ParamSignUp) (err error) {
	//0. estimate the existent of user
	if err = mysql.CheckUserExist(p.Username); err != nil {
		// search failed
		return err
	}
	//1. generate the UID
	userID := snowflake.GenID()

	//2. construct a User Instance
	u := &models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}

	//2. stores in the db
	return mysql.InsertUser(u)
}

func Login(p *models.ParamLogin) (user *models.User, err error) {
	user = &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	if err := mysql.Login(user); err != nil {
		return nil, err
	}

	//generate jwt token
	token, err := jwt.GenToken(user.UserID, user.Username)
	if err != nil {
		return
	}
	user.Token = token
	return
}
