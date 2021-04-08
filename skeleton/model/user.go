package model

import (
	"github.com/daodao97/egin/db"
)

// 后台用户
type UserEntity struct {
	Id            int    `json:"id"`
	Username      string `json:"username" comment:"用户名"`
	Realname      string `json:"realname"`
	Password      string `json:"password"`
	Mobile        string `json:"mobile"`
	Email         string `json:"email"`
	Status        int    `json:"status"`
	LoginTime     string `json:"login_time"`
	LoginIp       string `json:"login_ip"`
	IsAdmin       int    `json:"is_admin" comment:"is admin"`
	IsDefaultPass int    `json:"is_default_pass" comment:"是否初始密码1:是,0:否"`
	Qq            string `json:"qq" comment:"用户qq"`
	Roles         string `json:"roles"`
	Sign          string `json:"sign" comment:"签名"`
	Avatar        string `json:"avatar"`
	AvatarSmall   string `json:"avatar_small"`
	CreateAt      string `json:"create_at"`
	UpdateAt      string `json:"update_at"`
}

type UserModel struct {
	db.BaseModel
}

var User UserModel

func init() {
	User = *NewUserModel()
}

func NewUserModel() *UserModel {
	return &UserModel{
		BaseModel: db.BaseModel{
			Connection: "default",
			Table:      "user",
		},
	}
}
