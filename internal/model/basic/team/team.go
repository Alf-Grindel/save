package team

import (
	"github.com/Alf_Grindel/save/internal/model"
	"github.com/Alf_Grindel/save/internal/model/basic/user"
	"github.com/Alf_Grindel/save/pkg/utils"
)

type TeamUserVo struct {
	Id          int64       `json:"id"`
	TeamName    string      `json:"team_name"`
	Description string      `json:"description"`
	MaxNum      int         `json:"max_num"`
	ExpireTime  string      `json:"expire_time"`
	UserId      int64       `json:"user_id"`
	Status      string      `json:"status"`
	CreateTime  string      `json:"create_time"`
	UpdateTime  string      `json:"update_time"`
	CreateUser  user.UserVo `json:"create_user"`
	HasJoinNum  int         `json:"has_join_num"`
	HasJoin     bool        `json:"has_join"`
}

type AddTeamReq struct {
	TeamName    string `json:"team_name"`
	Description string `json:"description"`
	MaxNum      int    `json:"max_num"`
	ExpireTime  string `json:"expire_time"`
	Status      string `json:"status"`
	Password    string `json:"password"`
}

type AddTeamResp struct {
	Base utils.BaseResp `json:"base"`
	Id   int64          `json:"id"`
}

type UpdateTeamReq struct {
	Id          int64  `json:"id"`
	TeamName    string `json:"team_name"`
	Description string `json:"description"`
	ExpireTime  string `json:"expire_time"`
	Status      string `json:"status"`
	Password    string `json:"password"`
}

type UpdateTeamResp struct {
	Base utils.BaseResp `json:"base"`
	Team model.Team     `json:"team"`
}

type GetTeamReq struct {
	Id int64 `json:"team_id"`
}

type GetTeamResp struct {
	Base utils.BaseResp `json:"base"`
	Team model.Team     `json:"team"`
}

type ListTeamsReq struct {
	Id          int64   `json:"team_id"`
	IdList      []int64 `json:"team_id_list"`
	SearchText  string  `json:"search_text"`
	TeamName    string  `json:"team_name"`
	Description string  `json:"description"`
	MaxNum      int     `json:"max_num"`
	UserId      int64   `json:"user_id"`
	Status      string  `json:"status"`
}

type ListTeamsResp struct {
	Base utils.BaseResp `json:"base"`
	Data []TeamUserVo   `json:"data"`
}

type JoinTeamReq struct {
	Id       int64  `json:"team_id"`
	Password string `json:"password"`
}

type QuitTeamReq struct {
	Id int64 `json:"team_id"`
}

type DeleteTeamReq struct {
	Id int64 `json:"team_id"`
}

type ListMyCreateTeamsResp struct {
	Base utils.BaseResp `json:"base"`
	Data []TeamUserVo   `json:"data"`
}

type ListMyJoinTeams struct {
	Base utils.BaseResp `json:"base"`
	Data []TeamUserVo   `json:"data"`
}
