package db

import (
	"context"
	"github.com/Alf_Grindel/save/internal/model"
	"github.com/Alf_Grindel/save/internal/model/basic/team"
	"github.com/Alf_Grindel/save/pkg/constant"
	"github.com/Alf_Grindel/save/pkg/utils"
	"github.com/Alf_Grindel/save/pkg/utils/errno"
	"github.com/Alf_Grindel/save/pkg/utils/hlog"
	"gorm.io/gorm"
	"time"
)

var teamSnow = utils.NewSnowflake(constant.TeamTableMachineID)

// AddTeam create a new team
func AddTeam(ctx context.Context, tx *gorm.DB, team *model.Team, password string, expireTime time.Time) (int64, error) {
	team.Id = teamSnow.GenerateID()
	if password != "" {
		team.Password = utils.HashPassword(password)
	}
	if !expireTime.IsZero() {
		team.ExpireTime = expireTime
	}
	if err := tx.WithContext(ctx).Create(&team).Error; err != nil {
		hlog.Error(err)
		return -1, errno.SystemErr
	}
	return team.Id, nil
}

// QueryTeamById  query team by id
func QueryTeamById(ctx context.Context, id int64) (*model.Team, error) {
	res := &model.Team{}
	if err := DB.WithContext(ctx).Where("is_delete = 0 and id = ?", id).First(&res).Error; err != nil {
		hlog.Error(err)
		return nil, errno.SystemErr
	}
	return res, nil
}

func QueryTeamByUserId(ctx context.Context, userId int64) ([]model.Team, error) {
	var res []model.Team
	if err := DB.WithContext(ctx).Where("is_delete = 0 and user_id = ?", userId).Find(&res).Error; err != nil {
		hlog.Error(err)
		return nil, errno.SystemErr
	}
	return res, nil
}

// QueryTeam return model.Team list
func QueryTeam(ctx context.Context, query *team.ListTeamsReq) ([]model.Team, error) {
	var teams []model.Team

	res := DB.WithContext(ctx).Model(&model.Team{}).Where("is_delete = 0")

	if query.Id > 0 {
		res = res.Where("id = ?", query.Id)
	}
	if len(query.IdList) > 0 {
		res = res.Where("id IN ?", query.IdList)
	}
	if len(query.SearchText) > 0 {
		like := "%" + query.SearchText + "%"
		res = res.Where("team_name like ? or description like ?", like, like)
	}
	if len(query.TeamName) > 0 {
		res = res.Where("team_name = ?", query.TeamName)
	}
	if len(query.Description) > 0 {
		res = res.Where("description = ?", query.Description)
	}
	if query.MaxNum > 0 {
		res = res.Where("max_num = ?", query.MaxNum)
	}
	if query.UserId > 0 {
		res = res.Where("user_id = ?", query.UserId)
	}
	if len(query.Status) > 0 {
		res = res.Where("status = ?", query.Status)
	}
	res = res.Where("expire_time IS NULL OR expire_time > ?", time.Now())

	if err := res.Find(&teams).Error; err != nil {
		hlog.Error(err)
		return nil, errno.SystemErr
	}
	return teams, nil
}

// CountTeamByUserId query team by userId
func CountTeamByUserId(ctx context.Context, userId int64) (int64, error) {
	var total int64
	if err := DB.WithContext(ctx).Model(&model.Team{}).Where("is_delete = 0 and user_id = ?", userId).Count(&total).Error; err != nil {
		hlog.Error(err)
		return -1, errno.SystemErr
	}
	return total, nil
}

// UpdateTeam update team info
func UpdateTeam(ctx context.Context, team *model.Team) (*model.Team, error) {
	current := &model.Team{}

	DB.WithContext(ctx).Where("is_delete = 0 and id = ?", team.Id).First(&current)

	if err := DB.WithContext(ctx).Model(&current).Updates(&team).Error; err != nil {
		hlog.Error(err)
		return nil, errno.SystemErr
	}
	return team, nil
}

// delete team
func DeleteTeam(ctx context.Context, tx *gorm.DB, id int64) (bool, error) {
	res := &model.Team{}
	result := tx.WithContext(ctx).Where("id = ? and is_delete = 0", id).First(&res)
	if err := result.Update("is_delete", 1).Error; err != nil {
		hlog.Error(err)
		return false, errno.SystemErr
	}
	return true, nil
}
