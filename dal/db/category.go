package db

import (
	"database/sql"
	logger "modtest/gostudy/lesson1/log"
	"modtest/gostudy/lesson2/mercury/common"

	"github.com/jmoiron/sqlx"
)

//获取所有问题分类
func GetCategoryList() (categoryList []*common.Category, err error) {
	categoryList = []*common.Category{} //最好先初始化一下结构体. 避免返回出错
	sqlstr := "select category_id, category_name from category order by category_id asc"
	err = DB.Select(&categoryList, sqlstr)
	if err == sql.ErrNoRows {
		err = nil
		return
	}
	if err != nil {
		return
	}
	return
}

//根据分类id获取问题
func GetQuestionListByCategoryId(categoryId int64) (questionList []*common.Question, err error) {
	sqlstr := `select 
					question_id, caption, content, author_id, category_id, status, create_time
				from 
					question
				where 
					category_id = ?
				order by
					question_id asc`
	err = DB.Select(&questionList, sqlstr, categoryId)
	if err != nil {
		logger.Error("获取问题列表失败, 错误: %#v", err)
		return
	}
	return
}

//根据分类id获取分类信息(使用map查询多个)
func MGetCategory(categoryIds []int64) (categoryMap map[int64]*common.Category, err error) {
	sqlstr := `select
		category_id, category_name
	from 
		category
	where category_id in (?)`

	var interfaceSlice []interface{}
	for _, c := range categoryIds {
		interfaceSlice = append(interfaceSlice, c)
	}
	queryStr, params, err := sqlx.In(sqlstr, interfaceSlice)
	if err != nil {
		logger.Error("sql创建失败, 错误:%v", err)
		return
	}

	categoryMap = make(map[int64]*common.Category, len(categoryIds))
	var categoryList []*common.Category
	err = DB.Select(&categoryList, queryStr, params...) //这里记得展开切片
	if err != nil {
		logger.Error("查询失败, sqlstr: %v, 错误:%v", queryStr, err)
		return
	}

	for _, v := range categoryList {
		categoryMap[v.CategoryId] = v
	}

	return
}
