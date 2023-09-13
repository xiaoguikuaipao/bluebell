package mysql

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"

	"web_app/models"
)

const secret = "xuanyizhao"

// InsertUser Insert a new user data
func InsertUser(user *models.User) (err error) {
	//encrypt the password with MD5
	user.Password = encryptPassword(user.Password)

	//insert the data into db
	sqlStr := `insert into user(user_id, username, password) values (?, ?, ?)`
	_, err = db.Exec(sqlStr, user.UserID, user.Username, user.Password)
	return
}

// CheckUserExist check whether the user exist or not
func CheckUserExist(username string) (err error) {
	sqlStr := `select count(user_id) from user where username = ?`
	var count int
	if err = db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	if count > 0 {
		return ErrorUserExist
	}
	return err
	//...
}

func Login(user *models.User) (err error) {
	oPassword := user.Password
	sqlStr := `select user_id, username, password from user where username = ?`
	err = db.Get(user, sqlStr, user.Username)
	if err == sql.ErrNoRows {
		return ErrorUserNotExist
	}
	if err != nil {
		return err
	}
	//estimate the password correctness
	password := encryptPassword(oPassword)
	if password != user.Password {
		return ErrorPassword
	}
	return
}

func encryptPassword(opassword string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(opassword)))
}

func GetUserByID(uid int64) (*models.User, error) {
	user := new(models.User)
	sqlStr := `select user_id, username from user where user_id = ?`
	err := db.Get(user, sqlStr, uid)
	return user, err
}
