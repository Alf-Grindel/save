package utils

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Alf_Grindel/save/pkg/utils/errno"
)

type BaseResp struct {
	StatusCode int32
	StatusMsg  string
}

func buildBaseResp(err error) *BaseResp {
	if err == nil {
		return baseResp(errno.Success)
	}
	e := errno.Errno{}
	if errors.As(err, &e) {
		return baseResp(e)
	}
	s := errno.ConvertErr(err)
	return baseResp(s)
}

func RespWithErr(w http.ResponseWriter, err error) {
	base := buildBaseResp(err)
	resp := &BaseResp{
		StatusCode: base.StatusCode,
		StatusMsg:  base.StatusMsg,
	}
	_ = json.NewEncoder(w).Encode(resp)
}

func baseResp(err errno.Errno) *BaseResp {
	return &BaseResp{
		StatusCode: err.ErrCode,
		StatusMsg:  err.ErrMsg,
	}
}
