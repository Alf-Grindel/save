package constant

// connection information
const (
	MySQLDefaultDsn = "%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local"
)

// constants in the project
const (
	MachineID      = 0
	Salt           = "saveM814"
	UserLoginState = "User_Login_State"

	SessionKey     = "session"
	CtxUserInfoKey = "user"

	UserTableName = "users"
)
