package account

import (
	"modtest/gostudy/lesson2/session"
)

//初始化session
func InitSession(provider, addr string, options ...string) (err error) {
	return session.Init(provider, addr, options...)
}
