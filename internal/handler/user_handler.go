package handler

import (
	"context"
	"encoding/json"
	"github.com/Alf_Grindel/save/internal/dal/db"
	"github.com/Alf_Grindel/save/internal/model/basic/user"
	"github.com/Alf_Grindel/save/internal/service"
	"github.com/Alf_Grindel/save/pkg/constant"
	"github.com/Alf_Grindel/save/pkg/utils"
	"github.com/Alf_Grindel/save/pkg/utils/errno"
	"github.com/gorilla/sessions"
	"net/http"
)

type UserHandler struct {
	ctx   context.Context
	store sessions.Store
}

func NewUserHandler(ctx context.Context, store sessions.Store) *UserHandler {
	return &UserHandler{
		ctx:   ctx,
		store: store,
	}
}

func (uh *UserHandler) UserRegister(w http.ResponseWriter, r *http.Request) {
	var req user.UserRegisterReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespWithErr(w, errno.ConvertErr(err))
		return
	}
	defer r.Body.Close()
	id, err := service.NewUserService(uh.ctx).UserRegister(&req)
	if err != nil {
		utils.RespWithErr(w, err)
		return
	}

	resp := &user.UserRegisterResp{
		Base: utils.BaseResp{
			StatusCode: errno.SuccessCode,
			StatusMsg:  "OK",
		},
		Id: id,
	}
	_ = json.NewEncoder(w).Encode(resp)
}

func (uh *UserHandler) UserLogin(w http.ResponseWriter, r *http.Request) {
	var req user.UserLoginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespWithErr(w, errno.ConvertErr(err))
		return
	}
	defer r.Body.Close()
	u, err := service.NewUserService(uh.ctx).UserLogin(&req)
	if err != nil {
		utils.RespWithErr(w, err)
		return
	}
	session, err := uh.store.Get(r, constant.UserLoginState)
	if err != nil {
		utils.RespWithErr(w, errno.ConvertErr(err))
		return
	}
	session.Values["login"] = true
	session.Values["user_account"] = u.Account
	_ = session.Save(r, w)

	resp := &user.UserLoginResp{
		Base: utils.BaseResp{StatusCode: errno.SuccessCode, StatusMsg: "OK"},
		User: *u,
	}
	_ = json.NewEncoder(w).Encode(resp)
}

func (uh *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	u, ok := r.Context().Value(constant.CtxUserInfoKey).(*user.UserVo)
	if !ok || u == nil {
		utils.RespWithErr(w, errno.NotLoginErr)
		return
	}
	resp := &user.GetCurrentResp{
		Base: utils.BaseResp{StatusCode: errno.SuccessCode, StatusMsg: "OK"},
		User: *u,
	}
	_ = json.NewEncoder(w).Encode(resp)
}

func (uh *UserHandler) UserLogout(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(constant.SessionKey).(*sessions.Session)
	if !ok {
		utils.RespWithErr(w, errno.SystemErr)
	}
	session.Values["login"] = false
	delete(session.Values, "user_account")
	session.Options.MaxAge = -1
	_ = session.Save(r, w)
	utils.RespWithErr(w, errno.Success)
}

func (uh *UserHandler) UserUpdate(w http.ResponseWriter, r *http.Request) {
	var req user.UserUpdateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespWithErr(w, errno.ConvertErr(err))
		return
	}
	defer r.Body.Close()
	u, err := service.NewUserService(r.Context()).UserUpdate(&req)
	if err != nil {
		utils.RespWithErr(w, err)
		return
	}
	resp := &user.UserUpdateResp{
		Base: utils.BaseResp{StatusCode: errno.SuccessCode, StatusMsg: "OK"},
		User: *u,
	}
	_ = json.NewEncoder(w).Encode(resp)
}

func (uh *UserHandler) UserDrop(w http.ResponseWriter, r *http.Request) {
	var req user.UserDropReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespWithErr(w, errno.ConvertErr(err))
		return
	}
	userInfo, ok := r.Context().Value(constant.CtxUserInfoKey).(*user.UserVo)
	if !ok || userInfo == nil {
		utils.RespWithErr(w, errno.SystemErr)
		return
	}
	if userInfo.Account != req.Account {
		utils.RespWithErr(w, errno.ParamErr.WithMessage("account is not match"))
		return
	}
	if isDrop, err := db.DropUser(req.Account); err != nil && !isDrop {
		utils.RespWithErr(w, err)
		return
	}
	resp := &user.UserDropResp{
		Base: utils.BaseResp{StatusCode: errno.SuccessCode, StatusMsg: "OK"},
	}
	_ = json.NewEncoder(w).Encode(resp)
}

func (uh *UserHandler) SearchUserByTags(w http.ResponseWriter, r *http.Request) {
	var req user.SearchUserByTagsReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespWithErr(w, errno.ConvertErr(err))
		return
	}
	defer r.Body.Close()
	users, err := service.NewUserService(uh.ctx).SearchUserByTagsBySQL(&req)
	if err != nil {
		utils.RespWithErr(w, err)
		return
	}
	var userList []user.UserVo
	for _, u := range users {
		userList = append(userList, *u)
	}
	resp := &user.SearchUserByTagsResp{
		Base:  utils.BaseResp{StatusCode: errno.SuccessCode, StatusMsg: "OK"},
		Users: userList,
	}
	_ = json.NewEncoder(w).Encode(resp)
}
