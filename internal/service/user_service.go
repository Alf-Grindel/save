package service

import (
	"context"
	"encoding/json"
	"github.com/Alf_Grindel/save/internal/dal/db"
	"github.com/Alf_Grindel/save/internal/model/basic/user"
	"github.com/Alf_Grindel/save/pkg/constant"
	"github.com/Alf_Grindel/save/pkg/utils"
	"github.com/Alf_Grindel/save/pkg/utils/errno"
	"regexp"
)

type UserService struct {
	ctx context.Context
}

func NewUserService(ctx context.Context) *UserService {
	return &UserService{
		ctx: ctx,
	}
}

// UserRegister register user and return user id
func (s *UserService) UserRegister(req *user.UserRegisterReq) (int64, error) {
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
	if _, err := db.QueryUserByAccount(account); err == nil {
		return -1, errno.ParamErr.WithMessage("account already exist or password is wrong")
	}
	hashPassword := utils.HashPassword(password)
	id, err := db.CreateUser(account, hashPassword)
	if err != nil {
		return -1, err
	}
	return id, nil
}

// UserLogin login user and return safety user info
func (s *UserService) UserLogin(req *user.UserLoginReq) (*user.UserVo, error) {
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
	current, err := db.QueryUserByAccount(account)
	if err != nil {
		return nil, errno.ParamErr.WithMessage("account not exist or password is wrong ")
	}
	if !utils.ComparePassword(current.Password, password) {
		return nil, errno.ParamErr.WithMessage("account not exist or password is wrong ")
	}
	return GetSafeUser(current), nil
}

// UserUpdate update user and return safety user info
func (s *UserService) UserUpdate(req *user.UserUpdateReq) (*user.UserVo, error) {
	if len(req.Password) != 0 {
		if len(req.Password) < 8 {
			return nil, errno.ParamErr.WithMessage("password length must greater than 8")
		}
		req.Password = utils.HashPassword(req.Password)
	}

	userinfo, ok := s.ctx.Value(constant.CtxUserInfoKey).(*user.UserVo)
	if !ok || userinfo == nil {
		return nil, errno.SystemErr
	}
	current, err := db.UpdateUser(userinfo.Account, req)
	if err != nil {
		return nil, err
	}
	user := GetSafeUser(current)
	return user, nil
}

// SearchUserByTags search user by tags in memory return safety user info
func (s *UserService) SearchUserByTags(req *user.SearchUserByTagsReq) ([]*user.UserVo, error) {
	tags := req.Tags
	if len(tags) == 0 {
		return nil, errno.ParamErr.WithMessage("parameter is empty")
	}
	currents, err := db.QueryUser()
	if err != nil {
		return nil, err
	}
	var users []*user.UserVo
	for _, current := range *currents {
		var userTags []string
		if err := json.Unmarshal([]byte(current.Tags), &userTags); err != nil {
			continue
		}
		if containsAllTags(userTags, tags) {
			users = append(users, GetSafeUser(&current))
		}
	}
	return users, nil
}

// SearchUserByTagsBySQL search user by tags in sql return safety user info
func (s *UserService) SearchUserByTagsBySQL(req *user.SearchUserByTagsReq) ([]*user.UserVo, error) {
	tags := req.Tags
	if len(tags) == 0 {
		return nil, errno.ParamErr.WithMessage("parameter is empty")
	}
	currents, err := db.QueryUserByTags(tags)
	if err != nil {
		return nil, err
	}
	var users []*user.UserVo
	for _, current := range *currents {
		users = append(users, GetSafeUser(&current))
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
func GetSafeUser(current *db.User) *user.UserVo {
	if current == nil {
		return nil
	}
	user := &user.UserVo{
		Id:         current.Id,
		Account:    current.Account,
		Name:       current.Name,
		Avatar:     current.Avatar,
		Profile:    current.Profile,
		Tags:       current.Tags,
		CreateTime: current.CreateTime.Format("2006-01-02 15:04:05"),
	}
	return user
}
