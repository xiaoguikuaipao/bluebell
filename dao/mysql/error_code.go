package mysql

import "errors"

var (
	ErrorUserExist    = errors.New("user is exited")
	ErrorUserNotExist = errors.New("the username don't exist")
	ErrorPassword     = errors.New("password errors")
	ErrorInvalidID    = errors.New("invalid id")
)
