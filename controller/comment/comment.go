package comment

import (
	"html"
	logger "modtest/gostudy/lesson1/log"
	"modtest/gostudy/lesson2/mercury/common"
	"modtest/gostudy/lesson2/mercury/dal/db"
	"modtest/gostudy/lesson2/mercury/id_gen"
	"modtest/gostudy/lesson2/mercury/util"

	"github.com/gin-gonic/gin"
)

const (
	MinCommentContentSize = 10
)

//发表评论的接口
func PostCommentHandler(c *gin.Context) {
	var comment common.Comment
	err := c.ShouldBind(&comment)
	if err != nil {
		logger.Error("bind params failed, comment:%#v\n, err: %v", comment, err)
		util.ResponseError(c, util.ErrCodeParameter)
		return
	}
	// logger.Debug("测试参数绑定成功, comment %#v", comment)
	if len(comment.Content) <= MinCommentContentSize || comment.AnswerId == 0 {
		logger.Error("参数有误,len(content) : %v, question_id: %v", len(comment.Content), comment.AnswerId)
		util.ResponseError(c, util.ErrCodeParameter)
		return
	}
	//用户id
	userId := int64(123123)
	// userId, err := account.GetUserId(c)
	if err != nil || userId <= 0 {
		logger.Error("用户未登录")
		util.ResponseError(c, util.ErrCodeUserNotLogin)
		return
	}
	comment.AuthorId = userId

	//1. 针对content做一个转义, 防止xss漏洞
	comment.Content = html.EscapeString(comment.Content)
	//2. 生成评论id
	uintCommentId, err := id_gen.GetId()
	if err != nil {
		logger.Error("生成comment_id失败, 错误:%v", err)
		util.ResponseError(c, util.ErrCodeServerBusy)
		return
	}
	comment.CommentId = int64(uintCommentId)
	err = db.CreatePostComment(&comment)
	if err != nil {
		logger.Error("生成评论失败, 错误:%v", err)
		util.ResponseError(c, util.ErrCodeServerBusy)
		return
	}
	util.ResponseSuccess(c, nil)
}

//发表回复的接口
func PostReplyHandler(c *gin.Context) {
	var comment common.Comment
	err := c.BindJSON(&comment)
	if err != nil {
		logger.Error("绑定参数失败, 错误:%v", err)
		util.ResponseError(c, util.ErrCodeParameter)
		return
	}

	//校验字段
	if len(comment.Content) < MinCommentContentSize || comment.AnswerId == 0 || comment.ReplyCommentId == 0 || comment.ParentCommentId == 0 {
		logger.Error("参数有误,len(content) : %v, answer_id: %v, reply_comment_id: %v, parent_comment_id:%v", len(comment.Content), comment.AnswerId, comment.ReplyCommentId, comment.ParentCommentId)
		util.ResponseError(c, util.ErrCodeParameter)
		return
	}

	comment.Content = html.EscapeString(comment.Content)

	//userid
	userId := int64(123123)
	// userId, err := account.GetUserId(c)
	if err != nil || userId == 0 {
		logger.Error("用户未登录")
		util.ResponseError(c, util.ErrCodeUserNotLogin)
		return
	}
	comment.AuthorId = userId

	uintCommentId, err := id_gen.GetId()
	if err != nil {
		logger.Error("生成comment_id失败, 错误:%v", err)
		util.ResponseError(c, util.ErrCodeServerBusy)
		return
	}

	//要查询ReplyCommentId的author_id, 作为reply_author_id
	comment.CommentId = int64(uintCommentId)
	err = db.CreateReplyComment(&comment)
	if err != nil {
		logger.Error("生成回复失败, 错误:%v", err)
		util.ResponseError(c, util.ErrCodeServerBusy)
		return
	}
	util.ResponseSuccess(c, nil)
}

//获取评论列表接口
func CommentListHandler(c *gin.Context) {
	//1. 解析参数
	answerId, err := util.GetQueryInt64(c, "answer_id")
	if err != nil {
		util.ResponseError(c, util.ErrCodeParameter)
		return
	}

	offset, err := util.GetQueryInt64(c, "offset")
	if err != nil {
		util.ResponseError(c, util.ErrCodeParameter)
		return
	}

	limit, err := util.GetQueryInt64(c, "limit")
	if err != nil {
		util.ResponseError(c, util.ErrCodeParameter)
		return
	}

	//2. 获取评论列表
	commentList, count, err := db.GetCommentList(answerId, offset, limit)
	if err != nil {
		logger.Error("获取评论列表失败, answerId :%v, commentList: %#v, count:%v, err:%v", answerId, commentList, count, err)
		util.ResponseError(c, util.ErrCodeServerBusy)
		return
	}

	var userIdList []int64
	for _, v := range commentList {
		userIdList = append(userIdList, v.AuthorId, v.ReplyAuthorId)
	}

	userList, err := db.GetUserInfoList(userIdList)
	if err != nil {
		logger.Error("获取用户信息列表失败, 错误:%v", err)
		util.ResponseError(c, util.ErrCodeServerBusy)
		return
	}

	userInfoMap := make(map[int64]*common.UserInfo, len(userIdList))
	for _, user := range userList {
		userInfoMap[user.UserId] = user
	}

	for _, comment := range commentList {
		user, ok := userInfoMap[comment.AuthorId]
		if ok {
			comment.AuthorName = user.Username
		}
		replyUser, ok := userInfoMap[comment.ReplyAuthorId]
		if ok {
			comment.ReplyAuthorName = replyUser.Username
		}
	}

	var apiCommentList = &common.ApiCommentList{}
	apiCommentList.CommentList = commentList
	apiCommentList.TotalCount = count
	util.ResponseSuccess(c, apiCommentList)
}

//获取回复列表接口
func ReplyListHandler(c *gin.Context) {
	//1. 解析参数
	//回复某条评论, 所以要获取评论id
	commentId, err := util.GetQueryInt64(c, "comment_id")
	if err != nil {
		util.ResponseError(c, util.ErrCodeParameter)
		return
	}
	offset, err := util.GetQueryInt64(c, "offset")
	if err != nil {
		util.ResponseError(c, util.ErrCodeParameter)
		return
	}

	limit, err := util.GetQueryInt64(c, "limit")
	if err != nil {
		util.ResponseError(c, util.ErrCodeParameter)
		return
	}

	//查询parent_comment_id=comment_id的所有回复
	commentList, count, err := db.GetReplyList(commentId, offset, limit)
	if err != nil {
		logger.Error("get commentlist failed, commentList:%#v, count:%v, comment_id:%v, offset: %v, limit:%v, err:%v", commentList, count, commentId, offset, limit, err)
		util.ResponseError(c, util.ErrCodeServerBusy)
		return
	}

	util.ResponseSuccess(c, commentList)
}

//点赞接口
func LikeHandler(c *gin.Context) {
	var like common.Like
	err := c.BindJSON(&like)
	if err != nil {
		logger.Error("like handler param bind failed, err:%v", err)
		util.ResponseError(c, util.ErrCodeParameter)
		return
	}

	if like.Id == 0 || (like.LikeType != common.LikeTypeAnswer && like.LikeType != common.LikeTypeComment) {
		logger.Error("非法参数, id:%v, like.LikeType : %v, err:%v", like.Id, like.LikeType, err)
		util.ResponseError(c, util.ErrCodeParameter)
		return
	}

	if like.LikeType == common.LikeTypeAnswer {
		//回答点赞
		err = db.UpdateAnswerVoteupCount(like.Id)
		if err != nil {
			util.ResponseError(c, util.ErrCodeServerBusy)
			return
		}
	}

	if like.LikeType == common.LikeTypeComment {
		//评论或回复点赞
		err = db.UpdateCommentLikeCount(like.Id)
		if err != nil {
			util.ResponseError(c, util.ErrCodeServerBusy)
			return
		}
	}

	util.ResponseSuccess(c, nil)
}
