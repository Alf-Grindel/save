package errno

import (
	"errors"
	"fmt"
)

type Errno struct {
	ErrCode int32
	ErrMsg  string
}

func (e Errno) Error() string {
	return fmt.Sprintf("[%d] - %s", e.ErrCode, e.ErrMsg)
}

func NewErrno(code int32, msg string) Errno {
	return Errno{
		ErrCode: code,
		ErrMsg:  msg,
	}
}

func (e Errno) WithMessage(msg string) Errno {
	e.ErrMsg = msg
	return e
}

func (e Errno) WithFormat(format string, v ...interface{}) Errno {
	e.ErrMsg = fmt.Sprintf(format, v...)
	return e
}

func ConvertErr(e error) Errno {
	err := Errno{}
	if errors.As(e, &err) {
		return err
	}
	s := SystemErr
	s.ErrMsg = e.Error()
	return s
}
