package account

import (
	"fmt"
	"modtest/gostudy/lesson2/mercury/common"
	"modtest/gostudy/lesson2/mercury/dal/db"
	"modtest/gostudy/lesson2/mercury/id_gen"
	"modtest/gostudy/lesson2/mercury/middleware/account"
	"modtest/gostudy/lesson2/mercury/util"

	"github.com/gin-gonic/gin"
)

//账号模块, 包含注册和登录验证功能. 使用中间件保存session

//登录功能
func LoginHandler(c *gin.Context) {
	var userInfo common.UserInfo
	account.ProcessRequest(c)
	var err error
	defer func() {
		if err != nil {
			return
		}
		//用户登录成功, 那么我们设置user_id到用户的session中
		account.SetUserId(int64(userInfo.UserId), c)
		//要在数据返回前执行response
		account.ProcessReponse(c)

		util.ResponseSuccess(c, nil)
	}()

	//校验账号密码

	err = c.BindJSON(&userInfo)
	if err != nil {
		util.ResponseError(c, util.ErrCodeParameter)
		return
	}

	if len(userInfo.Username) == 0 || len(userInfo.Password) == 0 {
		util.ResponseError(c, util.ErrCodeParameter)
		return
	}

	err = db.Login(&userInfo)

	if err == db.ErrUserNotExist {
		util.ResponseError(c, util.ErrCodeUserExist)
		return
	}

	if err == db.ErrUserPasswordWrong {
		util.ResponseError(c, util.ErrCodeUserPasswordWrong)
		return
	}

	if err != nil {
		util.ResponseError(c, util.ErrCodeServerBusy)
		return
	}

}

//注册功能
func RegisterHandler(c *gin.Context) {
	var userInfo common.UserInfo
	err := c.BindJSON(&userInfo)
	if err != nil {
		fmt.Println(err)
		util.ResponseError(c, util.ErrCodeParameter)
		return
	}

	if len(userInfo.Email) == 0 || len(userInfo.Username) == 0 || len(userInfo.Password) == 0 {
		fmt.Println(err)
		util.ResponseError(c, util.ErrCodeParameter)
		return
	}

	//sex = 1男, sex = 2女
	if userInfo.Sex != common.UserSexMan && userInfo.Sex != common.UserSexWomen {
		fmt.Println(err)
		util.ResponseError(c, util.ErrCodeParameter)
		return
	}

	userId, err := id_gen.GetId()
	userInfo.UserId = int64(userId)
	if err != nil {
		fmt.Println(err)
		util.ResponseError(c, util.ErrCodeServerBusy)
		return
	}

	err = db.Register(&userInfo)
	if err == db.ErrUserExist {
		fmt.Println(err)
		util.ResponseError(c, util.ErrCodeUserExist)
		return
	}

	if err != nil {
		fmt.Println(err)
		util.ResponseError(c, util.ErrCodeServerBusy)
		return
	}

	util.ResponseSuccess(c, nil)
}
