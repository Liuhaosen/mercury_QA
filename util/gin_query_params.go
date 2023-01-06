package util

import (
	"fmt"
	logger "modtest/gostudy/lesson1/log"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

//处理请求参数的方法
func GetQueryInt64(c *gin.Context, key string) (value int64, err error) {
	idStr, ok := c.GetQuery(key)
	idStr = strings.TrimSpace(idStr)
	if !ok {
		if key == "offset" {
			idStr = "0"
		} else if key == "limit" {
			idStr = "10"
		} else {
			logger.Error("非法参数 %s, 未找到参数%s", key)
			err = fmt.Errorf("非法参数, 没有找到参数:%s", key)
			return
		}

	}

	value, err = strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Error("参数转换失败, 参数: %v, 错误:%v", idStr, err)
		err = fmt.Errorf("参数转换失败, 参数: %v, 错误:%v", idStr, err)
		return
	}
	return
}
