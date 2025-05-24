package redis

import (
	"github.com/Alf_Grindel/save/conf"
	"github.com/go-redsync/redsync/v4"

	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
	"time"
)

var (
	expireTime   = time.Hour * 1
	rdbRecommend *redis.Client

	// mutex lock
	rdbRedSync *redis.Client
	RedSync    *redsync.Redsync
)

func Init() {
	rdbRecommend = redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Addr,
		Password: "",
		DB:       0,
	})

	rdbRedSync = redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Addr,
		Password: "",
		DB:       3,
	})
	pool := goredis.NewPool(rdbRedSync)
	RedSync = redsync.New(pool)
}
