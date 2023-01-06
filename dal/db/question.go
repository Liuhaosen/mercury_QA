package db

import (
	logger "modtest/gostudy/lesson1/log"
	"modtest/gostudy/lesson2/mercury/common"
)

//创建问题
func CreateQuestion(question *common.Question) (err error) {
	sqlstr := `insert into question
		(question_id, author_id, caption, category_id, content) 
	values 
		(?, ?, ?, ?, ?)`
	_, err = DB.Exec(sqlstr, question.QuestionId, question.AuthorId, question.Caption, question.CategoryId, question.Content)
	if err != nil {
		logger.Error("创建问题失败, 错误:%v", err)
		return
	}
	return
}

//获取问题信息
func GetQuestionInfo(questionId int64) (question *common.Question, err error) {
	question = &common.Question{}
	sqlstr := `select 
		question_id, caption, content, author_id, category_id, create_time
	from
		question
	where 
		question_id = ?`
	err = DB.Get(question, sqlstr, questionId)
	if err != nil {
		logger.Error("get question failed, err : %v", err)
		return
	}

	return
}
