package db

import (
	"database/sql"
	"fmt"
	logger "modtest/gostudy/lesson1/log"
	"modtest/gostudy/lesson2/mercury/common"
	"modtest/gostudy/lesson2/mercury/util"

	"github.com/jmoiron/sqlx"
)

const (
	PasswordSalt = "E8N2O0xiEwk4OloTYJrMstroN2xJCasb"
)

//用户注册
func Register(user *common.UserInfo) (err error) {
	//查询数据库是否存在.
	var count int64
	sqlstr := "select count(user_id) from user where username = ?"
	err = DB.Get(&count, sqlstr, user.Username)
	if err != sql.ErrNoRows && err != nil {
		fmt.Println(err)
		return
	}
	if count > 0 {
		err = ErrUserExist
		fmt.Println(err)
		return
	}

	//插入数据库
	password := user.Password + PasswordSalt
	dbPassword := util.Md5([]byte(password))

	sqlstr = `insert into 
		user 
			(nickname, username, password, email, user_id, sex) 
		values 
			(?, ?, ?, ?, ?, ?)`
	_, err = DB.Exec(sqlstr, user.NickName, user.Username, dbPassword, user.Email, user.UserId, user.Sex)
	return
}

//用户登录
func Login(user *common.UserInfo) (err error) {
	originPassword := user.Password
	sqlstr := `select username, password, user_id from user 
		where
		username = ?`
	err = DB.Get(user, sqlstr, user.Username)
	if err != nil && err != sql.ErrNoRows {
		return
	}
	if err == sql.ErrNoRows {
		err = ErrUserNotExist
		return
	}

	password := originPassword + PasswordSalt
	originPassword = util.Md5([]byte(password))
	if originPassword != user.Password {
		err = ErrUserPasswordWrong
		return
	}

	return
}

//批量获取用户信息
func GetUserInfoList(userIdList []int64) (userInfoList []*common.UserInfo, err error) {
	if len(userIdList) == 0 {
		return
	}
	userInfoList = []*common.UserInfo{}
	sqlstr := `select user_id, username, sex, nickname, email from user where user_id in (?)`
	querySql, args, err := sqlx.In(sqlstr, userIdList)
	if err != nil {
		logger.Error("sqlx in failed, sqlstr : %v, user_ids: %#v, err: %v", sqlstr, userIdList, err)
		return
	}
	err = DB.Select(&userInfoList, querySql, args...)
	return
}
