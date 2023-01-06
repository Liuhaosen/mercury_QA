package util

import "fmt"

//自定义错误码
const (
	ErrCodeSuccess             = 0    //返回成功
	ErrCodeParameter           = 1001 //参数错误
	ErrCodeUserExist           = 1002 //用户已存在
	ErrCodeServerBusy          = 1003 //服务器繁忙
	ErrCodeUserNotExist        = 1004 //用户不存在
	ErrCodeUserPasswordWrong   = 1005 //用户名或密码错误
	ErrCodeCaptionHitSensitive = 1006 //问题标题有敏感词
	ErrCodeContentHitSensitive = 1007 //问题内容有敏感词
	ErrCodeUserNotLogin        = 1008 //用户未登录
	ErrCodeRecordExist         = 1009 //记录已存在
)

//错误描述
func GetErrMessage(code int) (message string) {
	switch code {
	case ErrCodeSuccess:
		message = "success(成功)"
	case ErrCodeParameter:
		message = "参数错误"
	case ErrCodeUserExist:
		message = "用户已存在"
	case ErrCodeServerBusy:
		message = "服务器繁忙"
	case ErrCodeUserNotExist:
		message = "用户名不存在"
	case ErrCodeUserPasswordWrong:
		message = "用户名或密码错误"
	case ErrCodeCaptionHitSensitive:
		message = "标题中含有非法内容, 请修改后发表"
	case ErrCodeContentHitSensitive:
		message = "内容含有非法内容, 请修改后发表"
	case ErrCodeUserNotLogin:
		message = "用户未登录"
	case ErrCodeRecordExist:
		message = "记录已存在"
	default:
		message = "未知错误"
		fmt.Println("未知错误, 错误码:", code)
	}
	return
}
