package db

import (
	logger "modtest/gostudy/lesson1/log"
	"modtest/gostudy/lesson2/mercury/common"

	"github.com/jmoiron/sqlx"
)

//根据question_id获取相关的answer_id列表
func GetAnswerIdList(questionId int64, offset, limit int64) (answerIdList []int64, err error) {
	sqlstr := `select 
		answer_id
	from 
		question_answer_rel
	where 
		question_id = ?
	limit ?, ?`
	err = DB.Select(&answerIdList, sqlstr, questionId, offset, limit)
	if err != nil {
		logger.Error("get answerid list failed, err:%v", err)
		return
	}
	return
}

//根据answerid列表获取回答列表
func GetAnswerList(answerIds []int64) (answerList []*common.Answer, err error) {
	sqlstr := `select 
			answer_id, content, author_id, comment_count, voteup_count, status, can_comment, create_time, update_time
		from 
			answer 
		where 
			answer_id in (?)`
	var interfaceSlice []interface{}
	for _, v := range answerIds {
		interfaceSlice = append(interfaceSlice, v)
	}

	queryStr, paramsSlice, err := sqlx.In(sqlstr, interfaceSlice)
	if err != nil {
		logger.Error("get query sql failed, sql:%v, err : %v", sqlstr, err)
		return
	}
	err = DB.Select(&answerList, queryStr, paramsSlice...)
	if err != nil {
		logger.Error("获取回答列表失败1, answer_ids : %v, querysql: %v err :%v", answerIds, queryStr, err)
		return
	}
	return
}

//获取问题回答的总数
func GetAnswerCount(questionId int64) (answerCount int32, err error) {
	sqlstr := `select count(answer_id) from question_answer_rel where question_id = ?`
	err = DB.Get(&answerCount, sqlstr, questionId)
	if err != nil {
		logger.Error("获取回答总数失败, question_id:%v, 错误:", questionId, err)
		return
	}
	return
}
