package service

import (
	"context"
	"github.com/Alf_Grindel/save/internal/dal/db"
)

type UserTeamService struct {
}

func NewUserTeamService() *UserTeamService { return &UserTeamService{} }

// 检查用户加入了多少队伍
func (s *UserTeamService) CountByUserId(ctx context.Context, userId int64) (int64, error) {
	hasJoinNum, err := db.CountUserTeamByUserId(ctx, userId)
	if err != nil {
		return -1, err
	}
	return hasJoinNum, nil
}

// 检查队伍加入了多少用户
func (s *UserTeamService) CountByTeamId(ctx context.Context, teamId int64) (int64, error) {
	total, err := db.CountUserTeamByTeamId(ctx, teamId)
	if err != nil {
		return -1, err
	}
	return total, nil
}
