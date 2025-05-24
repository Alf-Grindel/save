package user

import (
	"github.com/Alf_Grindel/save/pkg/utils"
)

type UserVo struct {
	Id         int64  `json:"id"`
	Account    string `json:"account"`
	Name       string `json:"name"`
	Avatar     string `json:"avatar"`
	Profile    string `json:"profile"`
	Tags       string `json:"tags"`
	CreateTime string `json:"create_time"`
}

type UserRegisterReq struct {
	Account       string `json:"account"`
	Password      string `json:"password"`
	CheckPassword string `json:"check_password"`
}

type UserRegisterResp struct {
	Base utils.BaseResp `json:"base"`
	Id   int64          `json:"id"`
}

type UserLoginReq struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type UserLoginResp struct {
	Base utils.BaseResp `json:"base"`
	User UserVo         `json:"user"`
}

type GetCurrentResp struct {
	Base utils.BaseResp `json:"base"`
	User UserVo         `json:"user"`
}

type UserUpdateReq struct {
	Password string `json:"password"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
	Profile  string `json:"profile"`
	Tags     string `json:"tags"`
}

type UserUpdateResp struct {
	Base utils.BaseResp `json:"base"`
	User UserVo         `json:"user"`
}

type UserDropReq struct {
	Account string `json:"account"`
}

type UserDropResp struct {
	Base utils.BaseResp `json:"base"`
}

type SearchUserByTagsReq struct {
	Tags []string `json:"tags"`
}

type SearchUserByTagsResp struct {
	Base  utils.BaseResp `json:"base"`
	Users []UserVo       `json:"user"`
}

type RecommendUserReq struct {
	PageSize    int64 `json:"page_size"`
	CurrentPage int64 `json:"current_page"`
}

type RecommendUserResp struct {
	Base  utils.BaseResp `json:"base"`
	Users []UserVo       `json:"user"`
}
