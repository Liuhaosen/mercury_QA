package db

import (
	"fmt"
	logger "modtest/gostudy/lesson1/log"
	"modtest/gostudy/lesson2/mercury/common"

	"github.com/jmoiron/sqlx"
)

//新增评论
func CreatePostComment(comment *common.Comment) (err error) {
	//要插入两个表, 所以需要开启事务
	tx, err := DB.Beginx()
	if err != nil {
		logger.Error("事务开启失败,错误: %v", err)
		return
	}
	//1. 插入comment表
	sqlstr := `insert into
		comment
			(comment_id, content, author_id)
		values
			(?, ?, ?)`
	_, err = tx.Exec(sqlstr, comment.CommentId, comment.Content, comment.AuthorId)
	if err != nil {
		logger.Error("插入comment失败, comment :%#v, err:%v", err)
		tx.Rollback()
		return
	}
	//2. 插入comment_rel评论关系表
	sqlstr = `insert into
		comment_rel
			(comment_id, parent_comment_id, level, answer_id, reply_author_id)
		values
			(?, ?, ?, ?, ?)`
	_, err = tx.Exec(sqlstr, comment.CommentId, comment.ParentCommentId, 1, comment.AnswerId, comment.ReplyAuthorId)
	if err != nil {
		logger.Error("插入comment_rel失败, comment:%#v, err:%v", comment, err)
		tx.Rollback()
		return
	}

	//3. 回答评论数加1
	sqlstr = `update answer set comment_count = comment_count + 1 where answer_id = ?`
	_, err = tx.Exec(sqlstr, comment.AnswerId)
	if err != nil {
		logger.Error("回答的评论数增加失败, sqlstr: %v, answer_id: %v, 错误: %v", sqlstr, comment.AnswerId, err)
		tx.Rollback()
		return
	}

	//4. 事务提交
	err = tx.Commit()
	if err != nil {
		logger.Error("事务提交失败, err:%v", err)
		tx.Rollback()
		return
	}
	return
}

//新增回复
func CreateReplyComment(reply *common.Comment) (err error) {
	//要插入两个表, 所以需要开启事务
	tx, err := DB.Beginx()
	if err != nil {
		logger.Error("事务开启失败,错误: %v", err)
		return
	}
	//1. 插入comment表
	sqlstr := `insert into
		comment
			(comment_id, content, author_id)
		values
			(?, ?, ?)`
	_, err = tx.Exec(sqlstr, reply.CommentId, reply.Content, reply.AuthorId)
	if err != nil {
		logger.Error("插入回复失败, comment :%#v, err:%v", err)
		tx.Rollback()
		return
	}
	//2. 根据reply_comment_id查询reply_author_id
	/* commentInfo, err := GetCommentInfo(reply.CommentId)
	if err != nil {
		tx.Rollback()
		return
	}
	reply.ReplyAuthorId = commentInfo.AuthorId */
	var replyAuthorId int64
	sqlstr = `select author_id from comment where comment_id = ?`
	err = tx.Get(&replyAuthorId, sqlstr, reply.CommentId)
	if err != nil {
		logger.Error("获取author_id信息失败, comment_id = %v, err: ", reply.CommentId, err)
		return
	}
	if replyAuthorId == 0 {
		err = fmt.Errorf("author_id invalid")
		return
	}
	reply.ReplyAuthorId = replyAuthorId

	//3. 插入comment_rel评论关系表
	sqlstr = `insert into
		comment_rel
			(comment_id, parent_comment_id, level, answer_id, reply_author_id, reply_comment_id)
		values
			(?, ?, ?, ?, ?, ?)`
	_, err = tx.Exec(sqlstr, reply.CommentId, reply.ParentCommentId, 2, reply.AnswerId, reply.ReplyAuthorId, reply.ReplyCommentId)
	if err != nil {
		logger.Error("插入comment_rel失败, reply:%#v, err:%v", reply, err)
		tx.Rollback()
		return
	}

	//4. 回答评论数加1
	sqlstr = `update answer set comment_count = comment_count + 1 where answer_id = ?`
	_, err = tx.Exec(sqlstr, reply.AnswerId)
	if err != nil {
		logger.Error("回答的评论数增加失败, sqlstr: %v, answer_id: %v, 错误: %v", sqlstr, reply.AnswerId, err)
		tx.Rollback()
		return
	}

	//5. 当前回复的评论的评论数也加一. 即comment_id = 当前reply的parent_comment_id
	sqlstr = `update comment set comment_count = comment_count + 1 where comment_id = ?`
	_, err = tx.Exec(sqlstr, reply.ParentCommentId)
	if err != nil {
		logger.Error("回复的评论的评论数增加失败, sqlstr: %v, parent_comment_id: %v, 错误: %v", sqlstr, reply.ParentCommentId, err)
		tx.Rollback()
		return
	}

	//6. 提交事务
	err = tx.Commit()
	if err != nil {
		logger.Error("事务提交失败, err:%v", err)
		tx.Rollback()
		return
	}
	return
}

//获取comment信息
func GetCommentInfo(commentId int64) (comment *common.Comment, err error) {
	sqlstr := `select comment_id, content, author_id, like_count, comment_count, create_time, update_time from comment where comment_id = ?`
	err = DB.Get(comment, sqlstr, commentId)
	if err != nil {
		logger.Error("获取comment信息失败, comment_id = %v, err: ", commentId, err)
		return
	}
	return
}

//获取评论列表
func GetCommentList(answerId int64, offset, limit int64) (commentList []*common.Comment, count int64, err error) {
	//1. 先查评论关系表
	var commentIdList []int64
	sqlstr := `select comment_id from comment_rel where answer_id = ? and level = 1 limit ?, ?`
	err = DB.Select(&commentIdList, sqlstr, answerId, offset, limit)
	if err != nil {
		logger.Error("get commentIdList failed, answer_id : %v, err: %v", answerId, err)
		return
	}
	if len(commentIdList) == 0 {
		logger.Error("没有获取到comment列表")
		return
	}
	//2. 查询评论表(其实这里最好使用redis map来进行k-v查询)
	sqlstr = `select 
		comment_id, content, author_id, like_count, comment_count, create_time
	from 
		comment
	where
		comment_id in (?)`
	// var interfaceSlice []interface{}
	// for _, v := range commentIdList {
	// 	interfaceSlice = append(interfaceSlice, v)
	// }
	queryStr, params, err := sqlx.In(sqlstr, commentIdList)
	if err != nil {
		logger.Error("sqlx in failed, sqlstr : %v, err: %v", sqlstr, err)
		return
	}
	err = DB.Select(&commentList, queryStr, params...)
	if err != nil {
		logger.Error("select comment list failed, sqlstr : %v, err: %v", queryStr, err)
		return
	}
	//3. 查询总记录条数
	sqlstr = `select comment_id, reply_author_id from comment_rel where answer_id = ? and level = 1`
	var commentRelList []*common.Comment
	commentRelListMap := make(map[int64]*common.Comment)
	err = DB.Select(&commentRelList, sqlstr, answerId)

	if err != nil {
		logger.Error("get commentrelList failed, answer_id: %v, err: %v", answerId, err)
		return
	}

	count = int64(len(commentRelList))
	for _, commentRel := range commentRelList {
		commentRelListMap[commentRel.CommentId] = commentRel
	}
	for _, comment := range commentList {
		commentTemp, ok := commentRelListMap[comment.CommentId]
		if ok {
			comment.ReplyAuthorId = commentTemp.ReplyAuthorId
		}

	}
	// err = DB.Get(&count, sqlstr, answerId)
	// if err != nil {
	// 	logger.Error("get comment count failed, answer_id: %v, count:%v, err: %v", count, err)
	// 	return
	// }
	return
}

//获取回复列表, 根据comment_id查询parent_comment_id=comment_id的所有回复
func GetReplyList(commentId int64, offset, limit int64) (commentList []*common.Comment, count int64, err error) {
	//1. 先查评论关系表
	var commentIdList []int64
	sqlstr := `select comment_id from comment_rel where parent_comment_id = ? and level = 2 limit ?, ?`
	err = DB.Select(&commentIdList, sqlstr, commentId, offset, limit)
	if err != nil {
		logger.Error("get commentIdList failed, comment_id : %v, err: %v", commentId, err)
		return
	}
	if len(commentIdList) == 0 {
		logger.Error("没有获取到comment列表")
		return
	}
	//2. 查询评论表(其实这里最好使用redis map来进行k-v查询)
	sqlstr = `select 
		comment_id, content, author_id, like_count, comment_count, create_time
	from 
		comment
	where
		comment_id in (?)`
	// var interfaceSlice []interface{}
	// for _, v := range commentIdList {
	// 	interfaceSlice = append(interfaceSlice, v)
	// }
	queryStr, params, err := sqlx.In(sqlstr, commentIdList)
	if err != nil {
		logger.Error("sqlx in failed, sqlstr : %v, err: %v", sqlstr, err)
		return
	}
	err = DB.Select(&commentList, queryStr, params...)
	if err != nil {
		logger.Error("select comment list failed, sqlstr : %v, err: %v", queryStr, err)
		return
	}
	//3. 查询总记录条数
	sqlstr = `select comment_id, reply_author_id from comment_rel where parent_comment_id = ? and level = 2`
	var commentRelList []*common.Comment
	commentRelListMap := make(map[int64]*common.Comment)
	err = DB.Select(&commentRelList, sqlstr, commentId)

	if err != nil {
		logger.Error("get commentrelList failed, parent_comment_id=: %v, err: %v", commentId, err)
		return
	}

	count = int64(len(commentRelList))
	for _, commentRel := range commentRelList {
		commentRelListMap[commentRel.CommentId] = commentRel
	}
	for _, comment := range commentList {
		commentTemp, ok := commentRelListMap[comment.CommentId]
		if ok {
			comment.ReplyAuthorId = commentTemp.ReplyAuthorId
		}
	}
	// err = DB.Get(&count, sqlstr, answerId)
	// if err != nil {
	// 	logger.Error("get comment count failed, answer_id: %v, count:%v, err: %v", count, err)
	// 	return
	// }
	return
}

//回答点赞
func UpdateAnswerVoteupCount(answerId int64) (err error) {
	sqlstr := `update answer 
		set
			voteup_count = voteup_count + 1
		where answer_id = ?`
	_, err = DB.Exec(sqlstr, answerId)
	if err != nil {
		logger.Error("回答点赞失败, answer_id: %v, err: %v", answerId, err)
		return
	}
	return
}

//评论或回复点赞
func UpdateCommentLikeCount(commentId int64) (err error) {
	sqlstr := `update comment 
		set
			like_count = like_count + 1
		where comment_id = ?`
	_, err = DB.Exec(sqlstr, commentId)
	if err != nil {
		logger.Error("回复/评论点赞失败, comment_id: %v, err: %v", commentId, err)
		return
	}
	return
}
