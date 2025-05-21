package conf

import (
	"log"
	"os"
	"path/filepath"

	"github.com/Alf_Grindel/save/pkg/utils/hlog"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	Server        *server
	DB            *db
	runtime_viper = viper.New()
)

func LoadConfig() {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	dir := getPath(path)
	runtime_viper.SetConfigName("config")
	runtime_viper.SetConfigType("yml")
	runtime_viper.AddConfigPath(dir)

	if err = runtime_viper.ReadInConfig(); err != nil {
		hlog.Fatal("read config file failed")
	}
	configMapping()

	runtime_viper.OnConfigChange(func(in fsnotify.Event) {
		hlog.Infof("notice config changed, %v\n", in.String())
		configMapping()
	})

	runtime_viper.WatchConfig()
}

func configMapping() {
	c := &Config{}

	if err := runtime_viper.Unmarshal(&c); err != nil {
		hlog.Fatal("config unmarshal failed, ", err)
	}
	Server = &c.Server
	DB = &c.Db
}

func getPath(path string) string {
	dir := path
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return filepath.Join(dir, "/conf")
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			hlog.Fatal("can not find the config file")
		}
		dir = parent
	}
}
