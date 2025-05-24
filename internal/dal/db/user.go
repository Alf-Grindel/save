package db

import (
	"context"
	"encoding/json"
	"github.com/Alf_Grindel/save/internal/model/basic/user"
	"github.com/Alf_Grindel/save/pkg/constant"
	"github.com/Alf_Grindel/save/pkg/utils"
	"github.com/Alf_Grindel/save/pkg/utils/errno"
	"github.com/Alf_Grindel/save/pkg/utils/hlog"
	"strings"
	"time"
)

type User struct {
	Id         int64     `json:"id"`
	Account    string    `json:"account"`
	Password   string    `json:"password"`
	Name       string    `json:"name"`
	Avatar     string    `json:"avatar"`
	Profile    string    `json:"profile"`
	Tags       string    `json:"tags"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
	IsDelete   int8      `json:"is_delete"`
}

func (User) TableName() string {
	return constant.UserTableName
}

var s = utils.NewSnowflake(constant.MachineID)

// CreateUser create user info
func CreateUser(ctx context.Context, account, hashPassword string) (int64, error) {
	snowId := s.GenerateID()

	u := &User{
		Id:       snowId,
		Account:  account,
		Password: hashPassword,
	}
	result := DB.WithContext(ctx).Select("id", "account", "password").Create(&u)
	if err := result.Error; err != nil {
		hlog.Error(err)
		return -1, errno.SystemErr
	}
	return u.Id, nil
}

// QueryUserByAccount query user by account
func QueryUserByAccount(ctx context.Context, account string) (*User, error) {
	u := &User{}
	result := DB.WithContext(ctx).Where("user_account = ?", account).First(&u)
	if err := result.Error; err != nil {
		hlog.Error(err)
		return nil, errno.SystemErr
	}
	return u, nil
}

// QueryUserById query user by id
func QueryUserById(ctx context.Context, id int64) (*User, error) {
	u := &User{}
	result := DB.WithContext(ctx).Where("id = ? and is_delete = 0", id).First(&u)
	if err := result.Error; err != nil {
		hlog.Error(err)
		return nil, errno.SystemErr
	}
	return u, nil
}

// QueryUserByTags query user by tags
func QueryUserByTags(ctx context.Context, tags []string) ([]User, error) {
	var users []User
	result := DB
	for _, tag := range tags {
		jsonTag, err := json.Marshal([]string{tag})
		if err != nil {
			hlog.Error(err)
			return nil, errno.SystemErr
		}
		result = result.Where("JSON_CONTAINS(tags, ?)", string(jsonTag))
	}
	if err := result.WithContext(ctx).Find(&users).Error; err != nil {
		hlog.Error(err)
		return nil, errno.SystemErr
	}
	return users, nil

}

// QueryUser query all user
func QueryUser(ctx context.Context) ([]User, error) {
	var users []User
	if err := DB.WithContext(ctx).Where("is_delete = 0").Find(&users).Error; err != nil {
		hlog.Error(err)
		return nil, errno.SystemErr
	}
	return users, nil
}

// QueryUserByList query all user by return list
func QueryUserByList(ctx context.Context, currentPage, pageSize int64) ([]User, error) {
	var users []User
	var total int64

	if err := DB.WithContext(ctx).Model(&users).Where("is_delete = 0").Count(&total).Error; err != nil {
		hlog.Error(err)
		return nil, errno.SystemErr
	}
	if total <= 0 {
		hlog.Error("sql having no data")
		return nil, errno.SystemErr
	}
	result := DB.WithContext(ctx).Limit(int(pageSize)).Offset(int(pageSize * (currentPage - 1))).Where("is_delete = 0").Find(&users)
	if err := result.Error; err != nil {
		hlog.Error(err)
		return nil, errno.SystemErr
	}
	return users, nil
}

// UpdateUser update user when user login
func UpdateUser(ctx context.Context, account string, req *user.UserUpdateReq) (*User, error) {
	u := &User{
		Password: strings.TrimSpace(req.Password),
		Name:     strings.TrimSpace(req.Name),
		Avatar:   strings.TrimSpace(req.Avatar),
		Profile:  strings.TrimSpace(req.Profile),
		Tags:     strings.TrimSpace(req.Tags),
	}
	current := &User{}
	DB.WithContext(ctx).Where("user_account = ?", account).First(&current)
	result := DB.WithContext(ctx).Model(&current).Updates(&u)
	if err := result.Error; err != nil {
		hlog.Error(err)
		return nil, errno.SystemErr
	}
	return current, nil
}

// DropUser delete login user
func DropUser(ctx context.Context, account string) (bool, error) {
	u := &User{}
	result := DB.WithContext(ctx).Where("user_account = ?", account).First(&u).Update("is_delete", 1)
	if err := result.Error; err != nil {
		hlog.Fatal(err)
		return false, errno.SystemErr
	}
	return true, nil
}
