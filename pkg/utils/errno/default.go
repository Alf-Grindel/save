package errno

const (
	SuccessCode = 0

	ParamErrCode     = 40000
	NotLoginErrCode  = 40100
	NoAuthErrCode    = 40101
	NotFoundErrCode  = 40400
	ForbiddenErrCode = 40300
	SystemErrCode    = 50000
)

var (
	Success = NewErrno(SuccessCode, "OK")

	ParamErr     = NewErrno(ParamErrCode, "Invalid parameter")
	NotLoginErr  = NewErrno(NotLoginErrCode, "User not logged in")
	NoAuthErr    = NewErrno(NoAuthErrCode, "Unauthorized access")
	NotFoundErr  = NewErrno(NotFoundErrCode, "Resource not found")
	ForbiddenErr = NewErrno(ForbiddenErrCode, "Access forbidden")
	SystemErr    = NewErrno(SystemErrCode, "Internal system error")
)
