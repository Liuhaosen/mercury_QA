package category

import (
	logger "modtest/gostudy/lesson1/log"
	"modtest/gostudy/lesson2/mercury/dal/db"
	"modtest/gostudy/lesson2/mercury/util"

	"github.com/gin-gonic/gin"
)

//获取所有问题分类
func GetCategoryListHandler(c *gin.Context) {
	categoryList, err := db.GetCategoryList()
	if err != nil {
		logger.Error("获取问题分类失败, 错误:%v", err)
		util.ResponseError(c, util.ErrCodeServerBusy)
		return
	}

	util.ResponseSuccess(c, categoryList)
}
