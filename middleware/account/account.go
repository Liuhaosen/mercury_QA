package account

import (
	"errors"
	"fmt"
	logger "modtest/gostudy/lesson1/log"
	"modtest/gostudy/lesson2/session"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
	账号中间件
	1. 具体请求处理之前, 通过cookie获取sessionid, 加载session
	2. 具体请求处理之后, 如果session有修改, 则设置cookie
	3. 暴露获取user_id和判断登录状态的两个接口
*/

//处理request.
//获取session里用户的登录信息.
//先获取cookie是否有sessionId, 有的话根据sessionid获取session对象. 用session获取userid
//已登录 : userid = session里的userid, loginstatus = 0
//未登录 : userid, loginstatus = 0
func ProcessRequest(c *gin.Context) {
	var userSession session.Session
	var err error
	defer func() {
		//如果没登录, 创建一个session
		if userSession == nil {
			userSession, err = session.CreateSession()
		}
		c.Set(MercurySessionName, userSession)
	}()

	cookie, err := c.Request.Cookie(CookieSessionId)
	if err != nil {
		//说明用户未登录, 把这个状态设置到gin框架里
		c.Set(MercuryUserId, int64(0))
		c.Set(MercuryUserLoginStatus, int64(0))
		return
	}

	//获取session_id
	sessionId := cookie.Value
	if len(sessionId) == 0 {
		//还是未登录
		c.Set(MercuryUserId, int64(0))
		c.Set(MercuryUserLoginStatus, int64(0))
		return
	}

	//根据session_id 获取对应session
	userSession, err = session.Get(sessionId)
	if err != nil {
		c.Set(MercuryUserId, int64(0))
		c.Set(MercuryUserLoginStatus, int64(0))

		return
	}

	tempUserId, err := userSession.Get(MercuryUserId)
	if err != nil {
		c.Set(MercuryUserId, int64(0))
		c.Set(MercuryUserLoginStatus, int64(0))
		return
	}

	userId, ok := tempUserId.(int64)
	if !ok || userId == 0 {
		c.Set(MercuryUserId, int64(0))
		c.Set(MercuryUserLoginStatus, int64(0))
		return
	}

	//登录
	c.Set(MercuryUserId, int64(userId))
	c.Set(MercuryUserLoginStatus, int64(1))
}

//设置用户id
func SetUserId(userId int64, c *gin.Context) {
	var userSession session.Session
	tempSession, exist := c.Get(MercurySessionName)
	if !exist {
		return
	}
	userSession, ok := tempSession.(session.Session)
	if !ok {
		return
	}
	userSession.Set(MercuryUserId, userId)
}

//获取用户userid
func GetUserId(c *gin.Context) (userId int64, err error) {
	tempUserid, ok := c.Get(MercuryUserId)
	if !ok {
		logger.Error("userid 不存在")
		err = errors.New("userid 不存在")
		return
	}

	userId, ok = tempUserid.(int64)
	if !ok {
		logger.Error("userid 转为64位失败")
		err = errors.New("userid 转为64位失败")
		return
	}
	return
}

//获取用户loginstatus
func IsLogin(c *gin.Context) (login bool) {
	login = false
	tempLoginStatus, ok := c.Get(MercuryUserLoginStatus)
	if !ok {
		return
	}

	loginStatus, ok := tempLoginStatus.(int64)
	if !ok {
		return
	}

	if loginStatus == 0 {
		return
	}
	login = true
	return
}

//处理response
//如果session有修改, 那么保存到cookie里.
func ProcessReponse(c *gin.Context) {
	var userSession session.Session
	tempSession, exist := c.Get(MercurySessionName)
	if !exist {
		//session不存在
		return
	}

	userSession, ok := tempSession.(session.Session)
	if !ok {
		return
	}

	if userSession == nil {
		return
	}
	//查看是否修改
	if !userSession.IsModify() {
		return
	}
	err := userSession.Save()
	if err != nil {
		return
	}
	fmt.Println("save session success")
	sessionId := userSession.GetId()
	cookie := &http.Cookie{
		Name:     CookieSessionId,
		Value:    sessionId,
		MaxAge:   CookieMaxAge,
		HttpOnly: true,
		Path:     "/",
	}

	http.SetCookie(c.Writer, cookie)
}
