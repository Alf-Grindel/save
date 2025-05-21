package db

import (
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

// CreateUser create user info
func CreateUser(account, hashPassword string) (int64, error) {
	snowId := utils.NewSnowflake(constant.MachineID).GenerateID()

	u := &User{
		Id:       snowId,
		Account:  account,
		Password: hashPassword,
	}
	result := DB.Select("id", "account", "password").Create(&u)
	if err := result.Error; err != nil {
		hlog.Error(err)
		return -1, errno.SystemErr
	}
	return u.Id, nil
}

// QueryUserByAccount query user by account
func QueryUserByAccount(account string) (*User, error) {
	u := &User{}
	result := DB.Where("user_account = ?", account).First(&u)
	if err := result.Error; err != nil {
		hlog.Error(err)
		return nil, errno.SystemErr
	}
	return u, nil
}

// QueryUserById query user by id
func QueryUserById(id int64) (*User, error) {
	u := &User{}
	result := DB.Where("id = ? and is_delete = 0", id).First(&u)
	if err := result.Error; err != nil {
		hlog.Error(err)
		return nil, errno.SystemErr
	}
	return u, nil
}

// QueryUserByTags query user by tags
func QueryUserByTags(tags []string) (*[]User, error) {
	var users *[]User
	result := DB
	for _, tag := range tags {
		jsonTag, err := json.Marshal([]string{tag})
		if err != nil {
			hlog.Error(err)
			return nil, errno.SystemErr
		}
		result = result.Where("JSON_CONTAINS(tags, ?)", string(jsonTag))
	}
	if err := result.Find(&users).Error; err != nil {
		hlog.Error(err)
		return nil, errno.SystemErr
	}
	return users, nil

}

// QueryUser query all user
func QueryUser() (*[]User, error) {
	var users *[]User
	if err := DB.Find(&users).Error; err != nil {
		hlog.Error(err)
		return nil, errno.SystemErr
	}
	return users, nil
}

// UpdateUser update user when user login
func UpdateUser(account string, req *user.UserUpdateReq) (*User, error) {
	u := &User{
		Password: strings.TrimSpace(req.Password),
		Name:     strings.TrimSpace(req.Name),
		Avatar:   strings.TrimSpace(req.Avatar),
		Profile:  strings.TrimSpace(req.Profile),
		Tags:     strings.TrimSpace(req.Tags),
	}
	current := &User{}
	DB.Where("user_account = ?", account).First(&current)
	result := DB.Model(&current).Updates(&u)
	if err := result.Error; err != nil {
		hlog.Error(err)
		return nil, errno.SystemErr
	}
	return current, nil
}

// DropUser delete login user
func DropUser(account string) (bool, error) {
	u := &User{}
	result := DB.Where("user_account = ?", account).First(&u).Update("is_delete", 1)
	if err := result.Error; err != nil {
		hlog.Fatal(err)
		return false, errno.SystemErr
	}
	return true, nil
}
