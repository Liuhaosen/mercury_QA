package common

const (
	LikeTypeAnswer  = 1 //回答点赞
	LikeTypeComment = 2 //评论或回复点赞
)

type Like struct {
	Id       int64 `json:"id"` //要点赞的answer_id或者comment_id
	LikeType int   `json:"type"`
}
