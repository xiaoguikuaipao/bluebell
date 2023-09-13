package controller

type ResCode int64

const (
	CodeSuccess ResCode = 1000 + iota
	CodeInvalidParam
	CodeUserExisted
	CodeUserNotExisted
	CodeInvalidPassword
	CodeServerBusy

	CodeInvalidToken
	CodeNeedLogin
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess:         "success",
	CodeInvalidParam:    "Request Params Error",
	CodeUserExisted:     "username is existed",
	CodeUserNotExisted:  "user don't exist",
	CodeInvalidPassword: "username or password error",
	CodeServerBusy:      "service busy",

	CodeInvalidToken: "token is invalid",
	CodeNeedLogin:    "need login",
}

func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeServerBusy]
	}
	return msg
}
