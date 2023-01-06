package account

import (
	"fmt"
	"log"
	"modtest/gostudy/lesson2/mercury/util"
	"time"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {
	ProcessRequest(c)
	isLogin := IsLogin(c)
	if !isLogin {
		fmt.Println("未登录")
		util.ResponseError(c, util.ErrCodeUserNotLogin)
		c.Abort() //中断当前请求
		return
	}
	c.Next()
}

func StatCost() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		c.Set("example", "12345")
		c.Next()
		latency := time.Since(t)
		log.Printf("共计耗时: %d 毫秒", latency)
	}
}
