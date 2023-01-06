package common

//用户信息
/*
user: "",
nickname: "",
sex: 1,
email: "",
password: ""
*/

const (
	UserSexMan   = 1
	UserSexWomen = 2
)

type UserInfo struct {
	UserId   int64  `json:"user_id" db:"user_id"`
	Username string `json:"user" db:"username"`
	NickName string `json:"nickname" db:"nickname"`
	Email    string `json:"email" db:"email"`
	Sex      int    `json:"sex" db:"sex"`
	Password string `json:"password" db:"password"`
}
