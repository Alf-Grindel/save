package constant

// connection information
const (
	MySQLDefaultDsn = "%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local"
	TCP             = "tcp"
	MaxIdleNum      = 10
)

// constants in the project
const (
	UserTableMachineID     = 0
	TeamTableMachineID     = 1
	UserTeamTableMachineID = 2

	Salt           = "saveM814"
	UserLoginState = "User_Login_State"

	SessionKey     = "session-secret"
	CtxUserInfoKey = "user"

	UserRecommendRedisKey = "save:user:recommend:%v:%d:%d"
	PrecacheLockRedisKey  = "save:precache:docache:lock"
	JoinTeamLockRedisKey  = "save:join_team:lock"

	MatchNum = 21

	UserTableName     = "users"
	TeamTableName     = "teams"
	UserTeamTableName = "user_team"

	PageSize    = 10
	CurrentPage = 1
)
