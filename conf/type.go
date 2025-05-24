package conf

type server struct {
	Port    string
	Name    string
	Version string
}

type db struct {
	Addr     string
	User     string
	Password string
	Database string
}

type redis struct {
	Addr string
}

type Config struct {
	Server server
	Db     db
	Redis  redis
}
