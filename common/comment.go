package common

import "time"

type Comment struct {
	CommentId       int64     `json:"comment_id" db:"comment_id"`
	Content         string    `json:"content" db:"content"`
	AuthorId        int64     `json:"author_id" db:"author_id"`
	LikeCount       int32     `json:"like_count" db:"like_count"`
	CommentCount    int32     `json:"comment_count" db:"comment_count"`
	ParentCommentId int64     `json:"parent_comment_id" db:"parent_comment_id"`
	AnswerId        int64     `json:"answer_id" db:"answer_id"`
	ReplyAuthorId   int64     `json:"reply_author_id" db:"reply_author_id"`
	ReplyCommentId  int64     `json:"reply_comment_id" db:"reply_comment_id"`
	CreateTime      time.Time `json:"create_time" db:"create_time"`
	AuthorName      string    `json:"author_name"`
	ReplyAuthorName string    `json:"reply_author_name"`
}

type ApiCommentList struct {
	CommentList []*Comment `json:"comment_list"`
	TotalCount  int64      `json:"total_count"`
}
