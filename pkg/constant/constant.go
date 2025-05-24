package constant

// connection information
const (
	MySQLDefaultDsn = "%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local"
	TCP             = "tcp"
	MaxIdleNum      = 10
)

// constants in the project
const (
	MachineID      = 0
	Salt           = "saveM814"
	UserLoginState = "User_Login_State"

	SessionKey     = "session-secret"
	CtxUserInfoKey = "user"

	UserRecommendRedisKey = "save:user:recommend:%v:%d:%d"
	PrecacheLockRedisKey  = "save:precache:docache:lock"

	UserTableName = "users"

	PageSize    = 10
	CurrentPage = 1
)
