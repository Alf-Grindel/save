package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Alf_Grindel/save/internal/dal/db"
	"github.com/Alf_Grindel/save/internal/middleware/redis"
	"github.com/Alf_Grindel/save/internal/model"
	"github.com/Alf_Grindel/save/internal/model/basic/user"
	"github.com/Alf_Grindel/save/pkg/constant"
	"github.com/Alf_Grindel/save/pkg/utils"
	"github.com/Alf_Grindel/save/pkg/utils/errno"
	"github.com/Alf_Grindel/save/pkg/utils/hlog"
	"regexp"
	"sort"
)

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}

// UserRegister register user and return user id
func (s *UserService) UserRegister(ctx context.Context, req *user.UserRegisterReq) (int64, error) {
	account := req.Account
	password := req.Password
	checkPassword := req.CheckPassword
	if len(account) == 0 || len(password) == 0 || len(checkPassword) == 0 {
		return -1, errno.ParamErr.WithMessage("parameter is empty")
	}
	if len(account) < 4 {
		return -1, errno.ParamErr.WithMessage("account length must greater than 4")
	}
	if len(password) < 8 || len(checkPassword) < 8 {
		return -1, errno.ParamErr.WithMessage("password length must greater than 8")
	}
	if password != checkPassword {
		return -1, errno.ParamErr.WithMessage("both input password is not equal")
	}
	reg := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	if !reg.MatchString(account) {
		return -1, errno.ParamErr.WithMessage("account can only contain numbers or letters")
	}
	if _, err := db.QueryUserByAccount(ctx, account); err == nil {
		return -1, errno.ParamErr.WithMessage("account already exist or password is wrong")
	}
	hashPassword := utils.HashPassword(password)
	id, err := db.CreateUser(ctx, account, hashPassword)
	if err != nil {
		return -1, err
	}
	return id, nil
}

// UserLogin login user and return safety user info
func (s *UserService) UserLogin(ctx context.Context, req *user.UserLoginReq) (*user.UserVo, error) {
	account := req.Account
	password := req.Password
	if len(account) == 0 || len(password) == 0 {
		return nil, errno.ParamErr.WithMessage("parameter is empty")
	}
	if len(account) < 4 {
		return nil, errno.ParamErr.WithMessage("account length must greater than 4")
	}
	if len(password) < 8 {
		return nil, errno.ParamErr.WithMessage("password length must greater than 8")
	}
	reg := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	if !reg.MatchString(account) {
		return nil, errno.ParamErr.WithMessage("account can only contain numbers or letters")
	}
	current, err := db.QueryUserByAccount(ctx, account)
	if err != nil {
		return nil, errno.ParamErr.WithMessage("account not exist or password is wrong ")
	}
	if !utils.ComparePassword(current.Password, password) {
		return nil, errno.ParamErr.WithMessage("account not exist or password is wrong ")
	}
	return GetSafeUser(current), nil
}

// UserUpdate update user and return safety user info
func (s *UserService) UserUpdate(ctx context.Context, req *user.UserUpdateReq) (*user.UserVo, error) {
	if len(req.Password) != 0 {
		if len(req.Password) < 8 {
			return nil, errno.ParamErr.WithMessage("password length must greater than 8")
		}
		req.Password = utils.HashPassword(req.Password)
	}

	userinfo, ok := ctx.Value(constant.CtxUserInfoKey).(*user.UserVo)
	if !ok || userinfo == nil {
		return nil, errno.SystemErr
	}
	current, err := db.UpdateUser(ctx, userinfo.Account, req)
	if err != nil {
		return nil, err
	}
	u := GetSafeUser(current)
	return u, nil
}

// SearchUserByTags search user by tags in memory return safety user info
func (s *UserService) SearchUserByTags(ctx context.Context, req *user.SearchUserByTagsReq) ([]user.UserVo, error) {
	tags := req.Tags
	if len(tags) == 0 {
		return nil, errno.ParamErr.WithMessage("parameter is empty")
	}
	currents, err := db.QueryUser(ctx)
	if err != nil {
		return nil, err
	}
	var users []user.UserVo
	for _, current := range currents {
		var userTags []string
		if err := json.Unmarshal([]byte(current.Tags), &userTags); err != nil {
			continue
		}
		if containsAllTags(userTags, tags) {
			users = append(users, *GetSafeUser(&current))
		}
	}
	return users, nil
}

// RecommendUser recommend user
func (s *UserService) RecommendUser(ctx context.Context, req *user.RecommendUserReq) ([]user.UserVo, error) {
	pageSize := req.PageSize
	currentPage := req.CurrentPage
	if pageSize <= 0 {
		pageSize = constant.PageSize
	}
	if currentPage <= 0 {
		currentPage = constant.CurrentPage
	}

	userInfo, ok := ctx.Value(constant.CtxUserInfoKey).(*user.UserVo)
	if !ok || userInfo == nil {
		return nil, errno.NotLoginErr
	}
	var rdbRecommend redis.Recommend
	var users []user.UserVo

	redisKey := fmt.Sprintf(constant.UserRecommendRedisKey, userInfo.Id, currentPage, pageSize)
	// check cache is or not having data
	if rdbRecommend.ExistRecommend(ctx, redisKey) {
		// if exist return in cache
		dataJson, err := rdbRecommend.GetRecommend(ctx, redisKey)
		if err != nil {
			hlog.Error("redis error, ", err)
			return nil, errno.SystemErr
		}
		err = json.Unmarshal([]byte(dataJson), &users)
		if err != nil {
			hlog.Error("json unmarshal error, ", err)
			return nil, errno.SystemErr
		}
		return users, nil
	}
	currents, err := db.QueryUserByList(ctx, currentPage, pageSize)
	if err != nil {
		return nil, err
	}
	for _, current := range currents {
		users = append(users, *GetSafeUser(&current))
	}
	// add data in cache
	data, err := json.Marshal(users)
	if err != nil {
		hlog.Error("json marshal error, ", err)
		return nil, errno.SystemErr
	}
	err = rdbRecommend.AddRecommend(ctx, redisKey, data)
	if err != nil {
		hlog.Error("add key to redis failed, ", err)
	}
	return users, nil
}

// containsAllTags if contains all tags return true
func containsAllTags(userTags, searchTags []string) bool {
	tagMap := make(map[string]struct{})
	for _, tag := range userTags {
		tagMap[tag] = struct{}{}
	}
	for _, t := range searchTags {
		if _, ok := tagMap[t]; !ok {
			return false
		}
	}
	return true
}

// GetSafeUser return safety user info
func GetSafeUser(current *model.User) *user.UserVo {
	if current == nil {
		return nil
	}
	u := &user.UserVo{
		Id:         current.Id,
		Account:    current.Account,
		UserName:   current.UserName,
		Avatar:     current.Avatar,
		Profile:    current.Profile,
		Tags:       current.Tags,
		CreateTime: current.CreateTime.Format("2006-01-02 15:04:05"),
	}
	return u
}

// SearchUserByTagsBySQL search user by tags in sql return safety user info
func (s *UserService) SearchUserByTagsBySQL(ctx context.Context, req *user.SearchUserByTagsReq) ([]user.UserVo, error) {
	tags := req.Tags
	if len(tags) == 0 {
		return nil, errno.ParamErr.WithMessage("parameter is empty")
	}
	currents, err := db.QueryUserByTags(ctx, tags)
	if err != nil {
		return nil, err
	}
	var users []user.UserVo
	for _, current := range currents {
		users = append(users, *GetSafeUser(&current))
	}
	return users, nil
}

func (s *UserService) MatchUsers(ctx context.Context) ([]user.UserVo, error) {
	userInfo, ok := ctx.Value(constant.CtxUserInfoKey).(*user.UserVo)
	if !ok || userInfo == nil {
		return nil, errno.NotLoginErr
	}
	if userInfo.Tags == "" {
		return nil, errno.ParamErr.WithMessage("当前用户标签为空")
	}
	var allUsers []model.User
	if err := db.DB.WithContext(ctx).Select("id", "tags").Where("tags is not null").
		Find(&allUsers).Error; err != nil {
		hlog.Error(err)
		return nil, errno.SystemErr
	}
	var loginTags []string
	if err := json.Unmarshal([]byte(userInfo.Tags), &loginTags); err != nil {
		hlog.Error("用户标签解析失败")
		return nil, errno.SystemErr
	}
	type userScore struct {
		Id       int64
		Distance int
	}
	var userScores []userScore
	for _, u := range allUsers {
		if u.Id == userInfo.Id || u.Tags == "" {
			continue
		}
		var userTags []string
		if err := json.Unmarshal([]byte(u.Tags), &userTags); err != nil {
			hlog.Error("用户标签解析失败")
			continue
		}
		distance := utils.MinDistance(loginTags, userTags)
		userScores = append(userScores, userScore{Id: u.Id, Distance: distance})
	}
	sort.Slice(userScores, func(i, j int) bool {
		return userScores[i].Distance < userScores[j].Distance
	})
	topUserScores := userScores[:constant.MatchNum]
	var matchUsers []user.UserVo
	for _, userIdScore := range topUserScores {
		current, err := db.QueryUserById(ctx, userIdScore.Id)
		if err != nil {
			return nil, err
		}
		matchUsers = append(matchUsers, *GetSafeUser(current))
	}
	return matchUsers, nil
}
