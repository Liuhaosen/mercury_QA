package favorite

import (
	logger "modtest/gostudy/lesson1/log"
	"modtest/gostudy/lesson2/mercury/common"
	"modtest/gostudy/lesson2/mercury/dal/db"
	"modtest/gostudy/lesson2/mercury/id_gen"
	"modtest/gostudy/lesson2/mercury/util"
	"strings"

	"github.com/gin-gonic/gin"
)

//添加收藏夹接口
func AddFavoriteDirHandler(c *gin.Context) {
	var favoriteDir common.FavoriteDir
	err := c.BindJSON(&favoriteDir)
	if err != nil {
		logger.Error("参数绑定失败, favorite_dir: %#v, err:%v", favoriteDir, err)
		util.ResponseError(c, util.ErrCodeParameter)
		return
	}

	favoriteDir.DirName = strings.TrimSpace(favoriteDir.DirName)

	if len(favoriteDir.DirName) == 0 {
		logger.Error("参数非法, dir_name:%v, err:%v", favoriteDir.DirName, err)
		util.ResponseError(c, util.ErrCodeParameter)
		return
	}

	uintDirId, err := id_gen.GetId()
	if err != nil {
		logger.Error("dir_id生成失败, 错误:%v", err)
		util.ResponseError(c, util.ErrCodeServerBusy)
		return
	}

	userId := int64(123123)
	// userId, err := account.GetUserId(c)
	if err != nil || userId == 0 {
		logger.Error("获取用户id失败, 错误: %v", err)
		util.ResponseError(c, util.ErrCodeUserNotLogin)
		return
	}
	favoriteDir.UserId = userId

	favoriteDir.DirId = int64(uintDirId)
	err = db.CreateFavoriteDir(&favoriteDir)
	if err != nil {
		if err == db.ErrRecordExist {
			util.ResponseError(c, util.ErrCodeRecordExist)
		} else {
			logger.Error("create favorite failed, favoriteDir: %#v, err:%v", favoriteDir, err)
			util.ResponseError(c, util.ErrCodeServerBusy)
		}
		return
	}
	util.ResponseSuccess(c, nil)
}

//添加收藏接口
func AddFavoriteHandler(c *gin.Context) {
	var favorite common.Favorite
	err := c.BindJSON(&favorite)
	if err != nil {
		logger.Error("参数绑定失败, favorite:%#v, err:%v", favorite, err)
		util.ResponseError(c, util.ErrCodeParameter)
		return
	}

	if favorite.AnswerId == 0 || favorite.DirId == 0 {
		logger.Error("参数有误, answer_id: %v, dir_id: %v, err:%v", favorite.AnswerId, favorite.DirId, err)
		util.ResponseError(c, util.ErrCodeParameter)
		return
	}

	userId := int64(123123)
	// userId, err := account.GetUserId()
	if err != nil || userId == 0 {
		util.ResponseError(c, util.ErrCodeUserNotLogin)
		return
	}

	favorite.UserId = userId
	err = db.CreateFavorite(&favorite)
	if err != nil {
		if err == db.ErrRecordExist {
			util.ResponseError(c, util.ErrCodeRecordExist)
		} else {
			logger.Error("create favorite failed, favorite: %#v, err:%v", favorite, err)
			util.ResponseError(c, util.ErrCodeServerBusy)
		}
		return
	}
	util.ResponseSuccess(c, favorite)
}

//收藏夹列表接口
func GetFavoriteDirListHandler(c *gin.Context) {
	//根据当前用户id获取列表
	var err error
	userId := int64(123123)
	// userId, err := account.GetUserId()
	if err != nil || userId == 0 {
		util.ResponseError(c, util.ErrCodeUserNotLogin)
		return
	}

	favoriteDirList, err := db.GetFavoriteDirList(userId)
	if err != nil {
		logger.Error("获取收藏夹列表失败, user_id:%v, err:%v", userId, err)
		util.ResponseError(c, util.ErrCodeServerBusy)
		return
	}
	util.ResponseSuccess(c, favoriteDirList)
}

//收藏列表接口
func GetFavoriteListHandler(c *gin.Context) {
	//根据收藏夹id获取收藏列表
	dirId, err := util.GetQueryInt64(c, "dir_id")
	if err != nil {
		util.ResponseError(c, util.ErrCodeParameter)
		return
	}

	userId := int64(123123)
	// userId, err := account.GetUserId()
	if err != nil || userId == 0 {
		util.ResponseError(c, util.ErrCodeUserNotLogin)
		return
	}

	offset, err := util.GetQueryInt64(c, "offset")
	if err != nil {
		util.ResponseError(c, util.ErrCodeParameter)
		offset = 0
	}

	limit, err := util.GetQueryInt64(c, "limit")
	if err != nil {
		util.ResponseError(c, util.ErrCodeParameter)
		limit = 10
	}

	favoriteList, err := db.GetFavoriteList(userId, dirId, offset, limit)
	if err != nil {
		logger.Error("获取收藏列表失败, userid:%v, dirid:%v, offset:%v, limit:%v, err:%v", userId, dirId, offset, limit, err)
		util.ResponseError(c, util.ErrCodeServerBusy)
		return
	}

	//返回给用户看的应该是answer的信息.
	//获取回答id列表
	var answerIdList []int64
	for _, favorite := range favoriteList {
		answerIdList = append(answerIdList, favorite.AnswerId)
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

	util.ResponseSuccess(c, apiAnswerList)
}
