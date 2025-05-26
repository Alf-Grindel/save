package service

import (
	"context"
	"github.com/Alf_Grindel/save/internal/dal/db"
	"github.com/Alf_Grindel/save/internal/middleware/redis"
	"github.com/Alf_Grindel/save/internal/model"
	"github.com/Alf_Grindel/save/internal/model/basic/team"
	"github.com/Alf_Grindel/save/internal/model/basic/user"
	"github.com/Alf_Grindel/save/pkg/constant"
	"github.com/Alf_Grindel/save/pkg/utils"
	"github.com/Alf_Grindel/save/pkg/utils/errno"
	"github.com/Alf_Grindel/save/pkg/utils/hlog"
	"github.com/go-redsync/redsync/v4"
	"gorm.io/gorm"
	"strings"
	"time"
)

type TeamService struct {
}

func NewTeamService() *TeamService { return &TeamService{} }

var validStatus = map[string]struct{}{
	"public":    {},
	"private":   {},
	"encrypted": {},
}

func checkLogin(ctx context.Context) (*user.UserVo, error) {
	userInfo, ok := ctx.Value(constant.CtxUserInfoKey).(*user.UserVo)
	if !ok || userInfo == nil {
		return nil, errno.NotLoginErr
	}
	return userInfo, nil
}

func (s *TeamService) AddTeam(ctx context.Context, req *team.AddTeamReq) (int64, error) {
	if req == nil {
		return -1, errno.ParamErr
	}
	userInfo, err := checkLogin(ctx)
	if err != nil {
		return -1, err
	}
	maxNum := req.MaxNum
	if maxNum < 1 || maxNum > 20 {
		return -1, errno.ParamErr.WithMessage("队伍人数不满足要求")
	}
	name := strings.TrimSpace(req.TeamName)
	if len(name) == 0 || len(name) > 20 {
		return -1, errno.ParamErr.WithMessage("队伍标题不满足要求")
	}
	description := strings.TrimSpace(req.Description)
	if len(description) == 0 || len(description) > 512 {
		return -1, errno.ParamErr.WithMessage("队伍描述过长")
	}
	status := strings.TrimSpace(req.Status)
	_, ok := validStatus[status]
	if len(status) > 0 && !ok {
		return -1, errno.ParamErr.WithMessage("队伍状态不满足要求")
	}
	if len(status) == 0 {
		status = "public"
	}
	password := strings.TrimSpace(req.Password)
	if status == "encrypted" && (len(password) == 0 || len(password) > 32) {
		return -1, errno.ParamErr.WithMessage("密码设置不正确")
	}
	loc, _ := time.LoadLocation("Asia/Shanghai")
	expireTime, err := time.ParseInLocation(time.DateTime, strings.TrimSpace(req.ExpireTime), loc)
	if err != nil {
		return -1, errno.ParamErr.WithMessage("超时时间格式应该为 2006-01-02 15:04:05")
	}
	if time.Now().After(expireTime) {
		return -1, errno.ParamErr.WithMessage("超时时间 > 当前时间")
	}
	total, err := db.CountTeamByUserId(ctx, userInfo.Id)
	if err != nil {
		return -1, err
	}
	if total > 5 {
		return -1, errno.ParamErr.WithMessage("用户最多创建5个队伍")
	}

	var id int64
	err = db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		newTeam := &model.Team{
			TeamName:    name,
			Description: description,
			MaxNum:      maxNum,
			UserId:      userInfo.Id,
			Status:      status,
		}
		id, err = db.AddTeam(ctx, tx, newTeam, password, expireTime)
		if err != nil {
			return err
		}
		userTeam := &model.UserTeam{
			UserId:   userInfo.Id,
			TeamId:   id,
			JoinTime: time.Now(),
		}
		if err = db.AddUserTeam(ctx, tx, userTeam); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (s *TeamService) UpdateTeam(ctx context.Context, req *team.UpdateTeamReq) (*model.Team, error) {
	if req == nil {
		return nil, errno.ParamErr
	}
	id := req.Id
	if id <= 0 {
		return nil, errno.ParamErr
	}
	oldTeam, err := db.QueryTeamById(ctx, id)
	if err != nil {
		return nil, errno.NotFoundErr.WithMessage("队伍不存在")
	}
	userInfo, err := checkLogin(ctx)
	if err != nil {
		return nil, err
	}
	if userInfo.Id != oldTeam.UserId {
		return nil, errno.NoAuthErr
	}
	status := req.Status
	_, ok := validStatus[status]
	if len(status) > 0 && !ok {
		return nil, errno.ParamErr.WithMessage("队伍状态不满足要求")
	}
	password := req.Password
	if status == "encrypted" && len(password) == 0 {
		return nil, errno.ParamErr.WithMessage("加密房间必须设置密码")
	}
	expireTime := oldTeam.ExpireTime
	if len(req.ExpireTime) > 0 {
		loc, _ := time.LoadLocation("Asia/Shanghai")
		expireTime, err = time.ParseInLocation(time.DateTime, strings.TrimSpace(req.ExpireTime), loc)
		if err != nil {
			return nil, errno.ParamErr.WithMessage("超时时间格式应该为 2006-01-02 15:04:05")
		}
		if time.Now().After(expireTime) {
			return nil, errno.ParamErr.WithMessage("超时时间 > 当前时间")
		}
	}
	updateTeam := &model.Team{
		Id:          req.Id,
		TeamName:    strings.TrimSpace(req.TeamName),
		Description: strings.TrimSpace(req.Description),
		ExpireTime:  expireTime,
		Status:      req.Status,
		Password:    strings.TrimSpace(req.Password),
	}
	newTeam, err := db.UpdateTeam(ctx, updateTeam)
	if err != nil {
		return nil, err
	}
	newTeam, err = db.QueryTeamById(ctx, newTeam.Id)
	if err != nil {
		return nil, err
	}
	return newTeam, nil
}

func (s *TeamService) ListTeams(ctx context.Context, req *team.ListTeamsReq) ([]team.TeamUserVo, error) {
	var teams []model.Team
	var err error
	if req != nil {
		if _, ok := validStatus[req.Status]; !ok {
			req.Status = ""
		}
		teams, err = db.QueryTeam(ctx, req)
		if err != nil {
			return nil, err
		}
	}
	var teamUserVoList []team.TeamUserVo
	for _, t := range teams {
		userId := t.UserId
		userInfo, err := db.QueryUserById(ctx, userId)
		if err != nil {
			return nil, err
		}
		createUser := GetSafeUser(userInfo)
		hasJoinNum, err := db.CountUserTeamByTeamId(ctx, t.Id)
		if err != nil {
			return nil, err
		}
		isJoin, err := db.CheckJoin(ctx, t.Id, userId)
		teamUserVo := &team.TeamUserVo{
			Id:          t.Id,
			TeamName:    t.TeamName,
			Description: t.Description,
			MaxNum:      t.MaxNum,
			ExpireTime:  t.ExpireTime.Format(time.DateTime),
			UserId:      t.UserId,
			Status:      t.Status,
			CreateTime:  t.CreateTime.Format(time.DateTime),
			UpdateTime:  t.UpdateTime.Format(time.DateTime),
			CreateUser:  *createUser,
			HasJoinNum:  int(hasJoinNum),
			HasJoin:     isJoin,
		}
		teamUserVoList = append(teamUserVoList, *teamUserVo)
	}
	return teamUserVoList, nil
}

func (s *TeamService) JoinTeam(ctx context.Context, req *team.JoinTeamReq) (bool, error) {
	if req == nil {
		return false, errno.ParamErr
	}
	teamId := req.Id
	if teamId <= 0 {
		return false, errno.ParamErr
	}
	userInfo, err := checkLogin(ctx)
	if err != nil {
		return false, err
	}
	teamInfo, err := db.QueryTeamById(ctx, teamId)
	if err != nil {
		return false, errno.NotFoundErr.WithMessage("队伍不存在")
	}
	if teamInfo.ExpireTime.Before(time.Now()) {
		return false, errno.ParamErr.WithMessage("队伍已过期")
	}
	if teamInfo.Status == "private" {
		return false, errno.ParamErr.WithMessage("禁止加入私有队伍")
	}
	if teamInfo.Status == "encrypted" && (!utils.ComparePassword(teamInfo.Password, req.Password) || len(req.Password) == 0) {
		return false, errno.ParamErr.WithMessage("密码错误")
	}
	lockTime := 30 * time.Second

	mx := redis.RedSync.NewMutex(constant.JoinTeamLockRedisKey,
		redsync.WithExpiry(lockTime),
		redsync.WithTries(1),
	)
	if err = mx.LockContext(ctx); err != nil {
		hlog.Error("Lock failed", err)
		return false, errno.SystemErr
	}
	defer mx.UnlockContext(ctx)

	watchdogCtx, watchdogCancel := context.WithCancel(ctx)
	defer watchdogCancel()

	go func() {
		ticker := time.NewTicker(lockTime / 2)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				ok, err := mx.Extend()
				if !ok || err != nil {
					hlog.Error("Failed to extend lock", err)
				}
			case <-watchdogCtx.Done():
				return
			}
		}
	}()

	if hasJoinNum, _ := NewUserTeamService().CountByUserId(ctx, userInfo.Id); hasJoinNum > 5 {
		return false, errno.ParamErr.WithMessage("用户最多创建和加入5个队伍")
	}
	if isJoin, _ := db.CheckJoin(ctx, teamId, userInfo.Id); isJoin {
		return false, errno.ParamErr.WithMessage("用户不能重复加入队伍")
	}
	if teamJoinNum, _ := NewUserTeamService().CountByTeamId(ctx, teamId); teamJoinNum >= int64(teamInfo.MaxNum) {
		return false, errno.ParamErr.WithMessage("队伍已满")
	}
	userTeam := &model.UserTeam{
		UserId:   userInfo.Id,
		TeamId:   teamId,
		JoinTime: time.Now(),
	}
	if err = db.AddUserTeam(ctx, db.DB, userTeam); err != nil {
		return false, err
	}
	return true, nil
}

func (s *TeamService) QuitTeam(ctx context.Context, req *team.QuitTeamReq) (bool, error) {
	if req == nil || req.Id <= 0 {
		return false, errno.ParamErr
	}
	id := req.Id
	teamInfo, err := db.QueryTeamById(ctx, id)
	if err != nil {
		return false, err
	}
	userInfo, err := checkLogin(ctx)
	if err != nil {
		return false, err
	}
	isJoin, err := db.CheckJoin(ctx, id, userInfo.Id)
	if err != nil {
		return false, err
	}
	if !isJoin {
		return false, errno.ParamErr.WithMessage("未加入该队伍")
	}
	teamJoinNum, err := NewUserTeamService().CountByTeamId(ctx, id)
	if err != nil {
		return false, err
	}
	err = db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if teamJoinNum == 1 {
			if _, err = db.DeleteTeam(ctx, tx, id); err != nil {
				return err
			}
		} else if userInfo.Id == teamInfo.UserId {
			nextTeamLeaderId, err := db.ChangeLeader(ctx, tx, id)
			if err != nil {
				return err
			}
			if err := tx.WithContext(ctx).Model(&model.Team{}).
				Where("is_delete = 0 and id = ?", id).Update("user_id", nextTeamLeaderId).Error; err != nil {
				hlog.Error(err)
				return errno.SystemErr
			}

		}
		_, err = db.DeleteUserTeam(ctx, tx, id, userInfo.Id)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *TeamService) DeleteTeam(ctx context.Context, req *team.DeleteTeamReq) (bool, error) {
	userInfo, err := checkLogin(ctx)
	if err != nil {
		return false, err
	}
	if req == nil || req.Id <= 0 {
		return false, errno.ParamErr
	}
	teamId := req.Id
	teamInfo, err := db.QueryTeamById(ctx, teamId)
	if err != nil {
		return false, err
	}
	if teamInfo.UserId != userInfo.Id {
		return false, errno.NoAuthErr
	}
	err = db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err = tx.WithContext(ctx).Model(&model.Team{}).
			Where("is_delete = 0 and id = ?", teamId).Update("is_delete", 1).Error; err != nil {
			hlog.Error(err)
			return errno.SystemErr
		}
		if err = tx.WithContext(ctx).Model(&model.UserTeam{}).
			Where("is_delete = 0 and team_id = ?", teamId).Update("is_delete", 1).Error; err != nil {
			hlog.Error(err)
			return errno.SystemErr
		}
		return nil
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *TeamService) ListMyCreateTeams(ctx context.Context) ([]team.TeamUserVo, error) {
	userInfo, err := checkLogin(ctx)
	if err != nil {
		return nil, err
	}
	userId := userInfo.Id
	total, err := db.CountTeamByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}
	if total < 1 {
		return nil, errno.ParamErr.WithMessage("未创建队伍")
	}
	var teamUserVoList []team.TeamUserVo
	teams, err := db.QueryTeamByUserId(ctx, userId)
	for _, t := range teams {
		userInfo, err := db.QueryUserById(ctx, userId)
		if err != nil {
			return nil, err
		}
		createUser := GetSafeUser(userInfo)
		hasJoinNum, err := NewUserTeamService().CountByTeamId(ctx, t.Id)
		if err != nil {
			return nil, err
		}
		teamUserVo := &team.TeamUserVo{
			Id:          t.Id,
			TeamName:    t.TeamName,
			Description: t.Description,
			MaxNum:      t.MaxNum,
			ExpireTime:  t.ExpireTime.Format(time.DateTime),
			UserId:      t.UserId,
			Status:      t.Status,
			CreateTime:  t.CreateTime.Format(time.DateTime),
			UpdateTime:  t.UpdateTime.Format(time.DateTime),
			CreateUser:  *createUser,
			HasJoinNum:  int(hasJoinNum),
			HasJoin:     true,
		}
		teamUserVoList = append(teamUserVoList, *teamUserVo)
	}
	return teamUserVoList, nil
}

func (s *TeamService) ListMyJoinTeams(ctx context.Context) ([]team.TeamUserVo, error) {
	userInfo, err := checkLogin(ctx)
	if err != nil {
		return nil, err
	}
	userId := userInfo.Id
	total, err := NewUserTeamService().CountByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}
	if total < 1 {
		return nil, errno.ParamErr.WithMessage("未加入队伍")
	}
	var teamUserVoList []team.TeamUserVo
	var teams []model.Team
	teamIds, err := db.QueryTeamIdByUserId(ctx, userId)
	for _, teamId := range teamIds {
		teamInfo, err := db.QueryTeamById(ctx, teamId)
		if err != nil {
			return nil, err
		}
		teams = append(teams, *teamInfo)
	}
	for _, t := range teams {
		userInfo, err := db.QueryUserById(ctx, userId)
		if err != nil {
			return nil, err
		}
		createUser := GetSafeUser(userInfo)
		hasJoinNum, err := NewUserTeamService().CountByTeamId(ctx, t.Id)
		if err != nil {
			return nil, err
		}
		teamUserVo := &team.TeamUserVo{
			Id:          t.Id,
			TeamName:    t.TeamName,
			Description: t.Description,
			MaxNum:      t.MaxNum,
			ExpireTime:  t.ExpireTime.Format(time.DateTime),
			UserId:      t.UserId,
			Status:      t.Status,
			CreateTime:  t.CreateTime.Format(time.DateTime),
			UpdateTime:  t.UpdateTime.Format(time.DateTime),
			CreateUser:  *createUser,
			HasJoinNum:  int(hasJoinNum),
			HasJoin:     true,
		}
		teamUserVoList = append(teamUserVoList, *teamUserVo)
	}
	return teamUserVoList, nil
}
