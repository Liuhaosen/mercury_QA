package question

import (
	logger "modtest/gostudy/lesson1/log"
	"modtest/gostudy/lesson2/mercury/common"
	"modtest/gostudy/lesson2/mercury/dal/db"
	"modtest/gostudy/lesson2/mercury/filter"
	"modtest/gostudy/lesson2/mercury/id_gen"
	"modtest/gostudy/lesson2/mercury/middleware/account"
	"modtest/gostudy/lesson2/mercury/util"
	"strconv"

	"github.com/gin-gonic/gin"
)

//提问 提交(提问完成后发送到kafka队列)
func QuestionSubmitHandler(c *gin.Context) {
	var question *common.Question
	err := c.BindJSON(&question)
	if err != nil {
		util.ResponseError(c, util.ErrCodeParameter)
		return
	}
	//1. 创建问题
	result, hit := filter.Replace(question.Caption, "***")
	if hit {
		logger.Error("标题含有敏感词, result:%#v", result)
		util.ResponseError(c, util.ErrCodeCaptionHitSensitive)
		return
	}

	result, hit = filter.Replace(question.Content, "***")
	if hit {
		logger.Error("内容含有敏感词, result:%#v", result)
		util.ResponseError(c, util.ErrCodeContentHitSensitive)
		return
	}

	qid, err := id_gen.GetId()
	if err != nil {
		logger.Error("问题id创建失败, err:%v", err)
		util.ResponseError(c, util.ErrCodeServerBusy)
		return
	}

	question.QuestionId = int64(qid)
	userId, err := account.GetUserId(c)

	if err != nil || userId <= 0 {
		logger.Error("用户未登录, err:%v", err)
		util.ResponseError(c, util.ErrCodeUserNotLogin)
		return
	}
	question.AuthorId = userId
	err = db.CreateQuestion(question)
	if err != nil {
		util.ResponseError(c, util.ErrCodeServerBusy)
		return
	}
	logger.Debug("create question success, question:%#v", question)
	util.ResponseSuccess(c, nil)

	//2. 把问题发送到kafka队列, 等待被消费, 消费成功后再放到es里
	err = util.SendtoKafka("mercury_topic", question)
	if err != nil {
		logger.Error("question send to kafka failed, err :%v", err)
		util.ResponseError(c, util.ErrCodeServerBusy)
		return
	}
}

//根据分类id获取所属问题
func GetQuestionListByCategoryidHandler(c *gin.Context) {
	categoryIdStr, ok := c.GetQuery("category_id")
	if !ok {
		logger.Error("category_id错误")
		util.ResponseError(c, util.ErrCodeParameter)
		return
	}
	categoryId, err := strconv.ParseInt(categoryIdStr, 10, 64)
	if err != nil {
		logger.Error("字符串转换失败, err:%v, 字符串:%v", err, categoryIdStr)
		util.ResponseError(c, util.ErrCodeParameter)
		return
	}
	questionList, err := db.GetQuestionListByCategoryId(categoryId)
	// util.ResponseSuccess(c, questionList)
	if err != nil {
		logger.Error("获取问题列表失败, err:%v, category_id = %d", err, categoryId)
		util.ResponseError(c, util.ErrCodeServerBusy)
		return
	}

	if len(questionList) == 0 {
		logger.Warn("问题列表数量为0")
		util.ResponseSuccess(c, questionList)
		return
	}

	var userIdList []int64
	userIdMap := make(map[int64]bool, 16)
	for _, question := range questionList {

		_, ok := userIdMap[question.AuthorId]
		if ok {
			continue
		}

		userIdMap[question.AuthorId] = true
		userIdList = append(userIdList, question.AuthorId)
	}

	userInfoList, err := db.GetUserInfoList(userIdList)
	if err != nil {
		logger.Error("get user info list failed, userids : %#v, err :%v", userIdList, err)
		return
	}
	var apiQuestionList []*common.ApiQuestion
	for _, question := range questionList {
		var apiQuestion = &common.ApiQuestion{}
		apiQuestion.Question = *question
		apiQuestion.CreateTimeStr = apiQuestion.CreateTime.Format("2006/01/02 15:04:05")
		for _, userInfo := range userInfoList {
			if question.AuthorId == userInfo.UserId {
				apiQuestion.AuthorName = userInfo.Username
				break
			}
		}
		apiQuestionList = append(apiQuestionList, apiQuestion)
	}

	util.ResponseSuccess(c, apiQuestionList)
}

//问题详情
func QuestionDetailHandler(c *gin.Context) {
	questionIdStr, ok := c.GetQuery("question_id")
	if !ok {
		logger.Error("非法参数,未找到question_id")
		util.ResponseError(c, util.ErrCodeParameter)
		return
	}

	questionId, err := strconv.ParseInt(questionIdStr, 10, 64)
	if err != nil {
		logger.Error("字符串转换失败, 错误:%v", err)
		util.ResponseError(c, util.ErrCodeParameter)
		return
	}

	//获取问题信息
	questionInfo, err := db.GetQuestionInfo(questionId)
	if err != nil {
		logger.Error("获取问题信息失败, 错误:%v", err)
		util.ResponseError(c, util.ErrCodeServerBusy)
		return
	}

	//获取问题的分类名
	categoryMap, err := db.MGetCategory([]int64{questionInfo.CategoryId})
	if err != nil {
		util.ResponseError(c, util.ErrCodeServerBusy)
	}
	category, ok := categoryMap[questionInfo.CategoryId]
	if !ok {
		logger.Error("获取分类失败, 错误:%v", err)
		util.ResponseError(c, util.ErrCodeServerBusy)
		return
	}

	//获取用户名
	userInfoList, err := db.GetUserInfoList([]int64{questionInfo.AuthorId})
	if err != nil || len(userInfoList) != 1 {
		logger.Error("获取用户信息失败, 错误:%v", err)
		util.ResponseError(c, util.ErrCodeServerBusy)
		return
	}
	apiQuestionDetail := &common.ApiQuestionDetail{}
	apiQuestionDetail.Question = *questionInfo
	apiQuestionDetail.AuthorName = userInfoList[0].Username
	apiQuestionDetail.CategoryName = category.CategoryName
	util.ResponseSuccess(c, apiQuestionDetail)
}
