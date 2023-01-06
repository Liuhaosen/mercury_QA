package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//错误/信息返回处理

type ResponseData struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

/*
	"code":0, //0表示成功, 其他表示失败
	"message" :"success", //用来描述失败原因
	"data":{}
*/
//错误返回
func ResponseError(c *gin.Context, code int) {
	responseData := &ResponseData{
		Code:    code,
		Message: GetErrMessage(code),
		Data:    make(map[string]interface{}),
	}
	c.JSON(http.StatusOK, responseData)
}

//成功返回
func ResponseSuccess(c *gin.Context, data interface{}) {
	responseData := &ResponseData{
		Code:    ErrCodeSuccess,
		Message: GetErrMessage(ErrCodeSuccess),
		Data:    data,
	}
	c.JSON(http.StatusOK, responseData)
}
