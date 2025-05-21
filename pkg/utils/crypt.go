package utils

import (
	"github.com/Alf_Grindel/save/pkg/constant"
	"github.com/Alf_Grindel/save/pkg/utils/hlog"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

func addSalt(password string) string {
	var builder strings.Builder
	builder.WriteString(password)
	builder.WriteString(constant.Salt)
	return builder.String()
}

func HashPassword(password string) string {
	plus := addSalt(password)
	b, err := bcrypt.GenerateFromPassword([]byte(plus), bcrypt.DefaultCost)
	if err != nil {
		hlog.Error("can not hash password")
		return ""
	}
	return string(b)
}

func ComparePassword(hashPassword, password string) bool {
	plus := addSalt(password)
	if err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(plus)); err != nil {
		hlog.Error("password can not match")
		return false
	}
	return true
}
