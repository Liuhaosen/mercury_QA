package db

import "errors"

//错误返回
var (
	ErrUserExist         = errors.New("用户名已存在")
	ErrUserNotExist      = errors.New("用户不存在")
	ErrUserPasswordWrong = errors.New("用户名或密码不正确")
	ErrRecordExist       = errors.New("记录已存在")
)
