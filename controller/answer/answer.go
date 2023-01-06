package answer

import (
	logger "modtest/gostudy/lesson1/log"
	"modtest/gostudy/lesson2/mercury/common"
	"modtest/gostudy/lesson2/mercury/dal/db"
	"modtest/gostudy/lesson2/mercury/util"

	"github.com/gin-gonic/gin"
)

//获取问题回答的列表接口
func GetAnswerListHandler(c *gin.Context) {
	//准备参数
	questionId, err := util.GetQueryInt64(c, "question_id")
	if err != nil {
		logger.Error("get question_id failed, err:", err)
		util.ResponseError(c, util.ErrCodeParameter)
		return
	}

	offset, err := util.GetQueryInt64(c, "offset")
	if err != nil {
		logger.Error("get offset failed, err:", err)
		util.ResponseError(c, util.ErrCodeParameter)
		return
	}

	limit, err := util.GetQueryInt64(c, "limit")
	if err != nil {
		logger.Error("get limit failed, err:", err)
		util.ResponseError(c, util.ErrCodeParameter)
		return
	}

	logger.Debug("获取answer列表参数完毕, qit:%v, offset:%v, limit:%v", questionId, offset, limit)

	//获取回答id列表
	answerIdList, err := db.GetAnswerIdList(questionId, offset, limit)
	if err != nil {
		logger.Error("db.GetAnswerIdList failed, question_id : %v, offset: %v, limit: %v", questionId, offset, limit)
		util.ResponseError(c, util.ErrCodeServerBusy)
		return
	}

	//获取回答列表
	answerList, err := db.GetAnswerList(answerIdList)
	if err != nil {
		logger.Error("获取回答列表失败2, answerids:%v, 错误:%v", answerIdList, err)
		util.ResponseError(c, util.ErrCodeServerBusy)
		return
	}

	//获取用户名
	var userIdList []int64
	for _, v := range answerList {
		userIdList = append(userIdList, v.AuthorId)
	}
	userInfoList, err := db.GetUserInfoList(userIdList)
	if err != nil {
		logger.Error("获取用户列表失败, userIdList:%v, 错误:%v", userIdList, err)
		util.ResponseError(c, util.ErrCodeServerBusy)
		return
	}
	apiAnswerList := &common.ApiAnsweList{}
	for _, v := range answerList {
		apiAnswer := &common.ApiAnswer{}
		apiAnswer.Answer = *v
		for _, userInfo := range userInfoList {
			if v.AuthorId == userInfo.UserId {
				apiAnswer.AuthorName = userInfo.Username
				break
			}
		}
		apiAnswerList.AnswerList = append(apiAnswerList.AnswerList, apiAnswer)
	}
	totalCount, err := db.GetAnswerCount(questionId)
	if err != nil {
		logger.Error("获取回答总数失败, question_id:%v, 错误:%v", questionId, err)
		util.ResponseError(c, util.ErrCodeServerBusy)
		return
	}
	apiAnswerList.TotalCount = totalCount
	util.ResponseSuccess(c, apiAnswerList)
}
