package model

import (
	"github.com/Alf_Grindel/save/pkg/constant"
	"time"
)

type User struct {
	Id         int64     `json:"id"`
	Account    string    `json:"account"`
	Password   string    `json:"password"`
	UserName   string    `json:"user_name"`
	Avatar     string    `json:"avatar"`
	Profile    string    `json:"profile"`
	Tags       string    `json:"tags"`
	CreateTime time.Time `json:"create_time" gorm:"autoUpdateTime"`
	UpdateTime time.Time `json:"update_time" gorm:"autoUpdateTime"`
	IsDelete   int8      `json:"is_delete"`
}

func (User) TableName() string {
	return constant.UserTableName
}

type Team struct {
	Id          int64     `json:"id"`
	TeamName    string    `json:"team_name"`
	Description string    `json:"description"`
	MaxNum      int       `json:"max_num"`
	ExpireTime  time.Time `json:"expire_time"`
	UserId      int64     `json:"user_id"`
	Status      string    `json:"status"`
	Password    string    `json:"password"`
	CreateTime  time.Time `json:"create_time" gorm:"autoUpdateTime"`
	UpdateTime  time.Time `json:"update_time" gorm:"autoUpdateTime"`
	IsDelete    int       `json:"is_delete"`
}

func (t Team) TableName() string {
	return constant.TeamTableName
}

type UserTeam struct {
	Id         int64     `json:"id"`
	UserId     int64     `json:"user_id"`
	TeamId     int64     `json:"team_id"`
	JoinTime   time.Time `json:"join_time"`
	CreateTime time.Time `json:"create_time" gorm:"autoUpdateTime"`
	UpdateTime time.Time `json:"update_time" gorm:"autoUpdateTime"`
	IsDelete   int       `json:"is_delete"`
}

func (ut UserTeam) TableName() string {
	return constant.UserTeamTableName
}
