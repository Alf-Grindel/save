package handler

import (
	"encoding/json"
	"github.com/Alf_Grindel/save/internal/dal/db"
	"github.com/Alf_Grindel/save/internal/model/basic/team"
	"github.com/Alf_Grindel/save/internal/service"
	"github.com/Alf_Grindel/save/pkg/utils"
	"github.com/Alf_Grindel/save/pkg/utils/errno"
	"net/http"
)

type TeamHandler struct {
}

func NewTeamHandler() *TeamHandler {
	return &TeamHandler{}
}

func (th *TeamHandler) AddTeam(w http.ResponseWriter, r *http.Request) {
	var req team.AddTeamReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespWithErr(w, errno.ConvertErr(err))
		return
	}
	id, err := service.NewTeamService().AddTeam(r.Context(), &req)
	if err != nil {
		utils.RespWithErr(w, err)
		return
	}
	resp := &team.AddTeamResp{
		Base: utils.BaseResp{StatusCode: errno.SuccessCode, StatusMsg: "OK"},
		Id:   id,
	}
	_ = json.NewEncoder(w).Encode(resp)
}

func (th *TeamHandler) UpdateTeam(w http.ResponseWriter, r *http.Request) {
	var req team.UpdateTeamReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespWithErr(w, errno.ConvertErr(err))
		return
	}
	t, err := service.NewTeamService().UpdateTeam(r.Context(), &req)
	if err != nil {
		utils.RespWithErr(w, err)
		return
	}
	resp := &team.UpdateTeamResp{
		Base: utils.BaseResp{StatusCode: errno.SuccessCode, StatusMsg: "OK"},
		Team: *t,
	}
	_ = json.NewEncoder(w).Encode(resp)
}

func (th *TeamHandler) GetTeamById(w http.ResponseWriter, r *http.Request) {
	var req team.GetTeamReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespWithErr(w, errno.ConvertErr(err))
		return
	}
	if req.Id <= 0 {
		utils.RespWithErr(w, errno.ParamErr)
	}
	t, err := db.QueryTeamById(r.Context(), req.Id)
	if err != nil {
		utils.RespWithErr(w, err)
		return
	}
	resp := &team.UpdateTeamResp{
		Base: utils.BaseResp{StatusCode: errno.SuccessCode, StatusMsg: "OK"},
		Team: *t,
	}
	_ = json.NewEncoder(w).Encode(resp)
}

func (th *TeamHandler) ListTeams(w http.ResponseWriter, r *http.Request) {
	var req team.ListTeamsReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespWithErr(w, errno.ConvertErr(err))
		return
	}

	userTeamList, err := service.NewTeamService().ListTeams(r.Context(), &req)
	if err != nil {
		utils.RespWithErr(w, err)
		return
	}
	resp := &team.ListTeamsResp{
		Base: utils.BaseResp{StatusCode: errno.SuccessCode, StatusMsg: "OK"},
		Data: userTeamList,
	}
	_ = json.NewEncoder(w).Encode(resp)
}

func (th *TeamHandler) JoinTeam(w http.ResponseWriter, r *http.Request) {
	var req team.JoinTeamReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespWithErr(w, errno.ConvertErr(err))
		return
	}
	_, err := service.NewTeamService().JoinTeam(r.Context(), &req)
	if err != nil {
		utils.RespWithErr(w, err)
		return
	}
	utils.RespWithErr(w, errno.Success)
}

func (th *TeamHandler) QuitTeam(w http.ResponseWriter, r *http.Request) {
	var req team.QuitTeamReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespWithErr(w, errno.ConvertErr(err))
		return
	}
	_, err := service.NewTeamService().QuitTeam(r.Context(), &req)
	if err != nil {
		utils.RespWithErr(w, err)
		return
	}
	utils.RespWithErr(w, errno.Success)
}

func (th *TeamHandler) DeleteTeam(w http.ResponseWriter, r *http.Request) {
	var req team.DeleteTeamReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespWithErr(w, errno.ConvertErr(err))
		return
	}
	_, err := service.NewTeamService().DeleteTeam(r.Context(), &req)
	if err != nil {
		utils.RespWithErr(w, err)
		return
	}
	utils.RespWithErr(w, errno.Success)
}

func (th *TeamHandler) ListMyCreateTeams(w http.ResponseWriter, r *http.Request) {
	userTeamList, err := service.NewTeamService().ListMyCreateTeams(r.Context())
	if err != nil {
		utils.RespWithErr(w, err)
		return
	}
	resp := &team.ListTeamsResp{
		Base: utils.BaseResp{StatusCode: errno.SuccessCode, StatusMsg: "OK"},
		Data: userTeamList,
	}
	_ = json.NewEncoder(w).Encode(resp)
}

func (th *TeamHandler) ListMyJoinTeams(w http.ResponseWriter, r *http.Request) {
	userTeamList, err := service.NewTeamService().ListMyJoinTeams(r.Context())
	if err != nil {
		utils.RespWithErr(w, err)
		return
	}
	resp := &team.ListTeamsResp{
		Base: utils.BaseResp{StatusCode: errno.SuccessCode, StatusMsg: "OK"},
		Data: userTeamList,
	}
	_ = json.NewEncoder(w).Encode(resp)
}
