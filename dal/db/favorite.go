package db

import (
	logger "modtest/gostudy/lesson1/log"
	"modtest/gostudy/lesson2/mercury/common"
)

//新增收藏夹
func CreateFavoriteDir(favoriteDir *common.FavoriteDir) (err error) {
	//1. 开启事务
	tx, err := DB.Beginx()
	if err != nil {
		logger.Error("开启事务失败, 错误:%v", err)
		return
	}

	//2. 首先查看是否已存在相同名字的收藏夹.
	//通过收藏夹名dir_name和user_id来查询效率会低, 所以要添加user_id和dir_name的联合索引
	var dirCount int64
	sqlstr := `select
		count(dir_id) 
	from 
		favorite_dir 
	where 
		user_id = ? and dir_name = ?`
	err = tx.Get(&dirCount, sqlstr, favoriteDir.UserId, favoriteDir.DirName)

	if err != nil {
		logger.Error("查询收藏夹失败, 错误:%v, favoriteDir:%#v", err, favoriteDir)
		tx.Rollback()
		return
	}

	if dirCount > 0 {
		tx.Rollback()
		err = ErrRecordExist
		logger.Error("记录已存在, dir_count : %v, err: %v", dirCount, err)
		return
	}

	//3. 执行插入
	sqlstr = `insert into
	favorite_dir
		(dir_id, dir_name, user_id)
	values
		(?, ?, ?)`

	_, err = tx.Exec(sqlstr, favoriteDir.DirId, favoriteDir.DirName, favoriteDir.UserId)
	if err != nil {
		logger.Error("insert favorite_dir failed, favorite_dir:%#v, err:%v", favoriteDir, err)
		tx.Rollback()
		return
	}
	err = tx.Commit()
	if err != nil {
		logger.Error("事务提交失败, err:%v", err)
		tx.Rollback()
		return
	}
	return
}

//新增收藏夹
func CreateFavorite(favorite *common.Favorite) (err error) {
	//1. 开启事务
	tx, err := DB.Beginx()
	if err != nil {
		logger.Error("开启事务失败, 错误:%v", err)
		return
	}

	//2. 首先查看是否已存在相同的收藏.
	var favoriteCount int64
	sqlstr := `select
		count(answer_id) 
	from 
		favorite
	where 
		user_id = ? and dir_id = ?`
	err = tx.Get(&favoriteCount, sqlstr, favorite.UserId, favorite.DirId)

	if err != nil {
		logger.Error("查询收藏失败, 错误:%v, favorite:%#v", err, favorite)
		tx.Rollback()
		return
	}

	if favoriteCount > 0 {
		tx.Rollback()
		err = ErrRecordExist
		logger.Error("收藏记录已存在, favorite_count : %v, err: %v", favoriteCount, err)
		return
	}

	//3. 执行插入
	sqlstr = `insert into
	favorite
		(dir_id, answer_id, user_id)
	values
		(?, ?, ?)`

	_, err = tx.Exec(sqlstr, favorite.DirId, favorite.AnswerId, favorite.UserId)
	if err != nil {
		logger.Error("insert favorite failed, favorite:%#v, err:%v", favorite, err)
		tx.Rollback()
		return
	}
	err = tx.Commit()
	if err != nil {
		logger.Error("事务提交失败, err:%v", err)
		tx.Rollback()
		return
	}
	return
}

//获取收藏夹列表
func GetFavoriteDirList(userId int64) (favoriteDirList []*common.FavoriteDir, err error) {
	sqlstr := `select
		dir_id, dir_name, count 
	from 
		favorite_dir
	where 
		user_id = ?`
	err = DB.Select(&favoriteDirList, sqlstr, userId)
	if err != nil {
		logger.Error("get favoriteDirList failed, sqlstr : %v, user_id:%v, err:%v", sqlstr, userId, err)
		return
	}

	return
}

//获取收藏夹列表
func GetFavoriteList(userId, dirId, offset, limit int64) (favoriteList []*common.Favorite, err error) {
	sqlstr := `select
		answer_id, dir_id, user_id
	from 
		favorite
	where 
		user_id = ? and dir_id = ?
	limit ?, ?`
	err = DB.Select(&favoriteList, sqlstr, userId, dirId, offset, limit)
	if err != nil {
		logger.Error("get favoriteList failed, sqlstr : %v, user_id:%v, dirId: %v, err:%v", sqlstr, userId, dirId, err)
		return
	}

	return
}
