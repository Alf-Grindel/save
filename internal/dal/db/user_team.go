package db

import (
	"context"
	"github.com/Alf_Grindel/save/internal/model"
	"github.com/Alf_Grindel/save/pkg/constant"
	"github.com/Alf_Grindel/save/pkg/utils"
	"github.com/Alf_Grindel/save/pkg/utils/errno"
	"github.com/Alf_Grindel/save/pkg/utils/hlog"
	"gorm.io/gorm"
)

var userTeamSnow = utils.NewSnowflake(constant.UserTeamTableMachineID)

// AddUserTeam
func AddUserTeam(ctx context.Context, tx *gorm.DB, userTeam *model.UserTeam) error {
	userTeam.Id = userTeamSnow.GenerateID()
	if err := tx.WithContext(ctx).Create(&userTeam).Error; err != nil {
		hlog.Error(err)
		return errno.SystemErr
	}
	return nil
}

func CountUserTeamByTeamId(ctx context.Context, teamId int64) (int64, error) {
	var total int64
	if err := DB.WithContext(ctx).Model(&model.UserTeam{}).Where("is_delete = 0 and team_id = ?", teamId).Count(&total).Error; err != nil {
		hlog.Error(err)
		return -1, errno.SystemErr
	}
	return total, nil
}

// CountUserTeamByUserId count by userId
func CountUserTeamByUserId(ctx context.Context, userId int64) (int64, error) {
	var total int64
	if err := DB.WithContext(ctx).Model(&model.UserTeam{}).Where("is_delete = 0 and user_id = ?", userId).Count(&total).Error; err != nil {
		hlog.Error(err)
		return -1, errno.SystemErr
	}
	return total, nil
}

func QueryTeamIdByUserId(ctx context.Context, userId int64) ([]int64, error) {
	var teamIds []int64
	var teamUsers []model.UserTeam
	if err := DB.WithContext(ctx).Model(&model.UserTeam{}).Where("is_delete = 0 and user_id = ?", userId).Find(&teamUsers).Error; err != nil {
		hlog.Error(err)
		return nil, errno.SystemErr
	}
	for _, teamUser := range teamUsers {
		teamIds = append(teamIds, teamUser.TeamId)
	}

	return teamIds, nil
}

func CheckJoin(ctx context.Context, teamId, userId int64) (bool, error) {
	userTeam := &model.UserTeam{}

	if err := DB.WithContext(ctx).Where("is_delete = 0 and user_id = ? and team_id = ?", userId, teamId).First(&userTeam).Error; err != nil {
		hlog.Error(err)
		return false, errno.SystemErr
	}
	return true, nil
}

func DeleteUserTeam(ctx context.Context, tx *gorm.DB, teamId, userId int64) (bool, error) {
	userTeam := &model.UserTeam{}
	res := tx.WithContext(ctx).Where("is_delete = 0 and user_id = ? and team_id = ?", userId, teamId).First(&userTeam)
	if err := res.Update("is_delete", 1).Error; err != nil {
		hlog.Error(err)
		return false, errno.SystemErr
	}
	return true, nil
}

func ChangeLeader(ctx context.Context, tx *gorm.DB, teamId int64) (int64, error) {
	var userTeams []model.UserTeam
	res := tx.WithContext(ctx).Where("is_delete = 0 and team_id = ?", teamId).Order("join_time asc").Limit(2).Find(&userTeams)
	if err := res.Error; err != nil {
		hlog.Error(err)
		return -1, errno.SystemErr
	}
	if len(userTeams) < 2 {
		return -1, errno.SystemErr
	}
	return userTeams[1].UserId, nil
}
