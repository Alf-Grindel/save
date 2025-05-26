package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Alf_Grindel/save/internal/dal/db"
	"github.com/Alf_Grindel/save/internal/middleware/redis"
	"github.com/Alf_Grindel/save/internal/model/basic/user"
	"github.com/Alf_Grindel/save/internal/service"
	"github.com/Alf_Grindel/save/pkg/constant"
	"github.com/Alf_Grindel/save/pkg/utils/hlog"
	"github.com/go-redsync/redsync/v4"
	"github.com/robfig/cron/v3"
	"math/rand"
	"sync"
	"time"
)

func getMainUsers() []int {
	var mainUsers []int
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	accounts := []string{"save"}
	for _, account := range accounts {
		info, err := db.QueryUserByAccount(ctx, account)
		if err != nil {
			hlog.Error(err)
			continue
		}
		u := service.GetSafeUser(info)
		mainUsers = append(mainUsers, int(u.Id))
	}
	return mainUsers
}

func DoCacheRecommendUser(taskCtx context.Context) {
	mainUsers := getMainUsers()
	var rdbRecommend redis.Recommend
	c := cron.New()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	lockTime := 30 * time.Second

	c.AddFunc("0 0 0 * * *", func() {
		mx := redis.RedSync.NewMutex(
			constant.PrecacheLockRedisKey,
			redsync.WithExpiry(lockTime),
			redsync.WithTries(1),
		)
		err := mx.LockContext(ctx)
		if err != nil {
			hlog.Error("Lock failed", err)
		}
		defer mx.UnlockContext(ctx)

		watchdogCtx, watchdogCancel := context.WithCancel(ctx)
		defer watchdogCancel()

		go func() {
			ticker := time.NewTicker(lockTime / 2)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					ok, err := mx.Extend()
					if !ok || err != nil {
						hlog.Error("Failed to extend lock", err)
						continue
					}
				case <-watchdogCtx.Done():
					return
				}
			}
		}()

		for _, userId := range mainUsers {
			wg.Add(1)
			go func(userId int) {
				defer wg.Done()
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				users := make([]user.UserVo, 0)
				page := rand.Intn(20)
				currents, err := db.QueryUserByList(ctx, int64(page), constant.PageSize)
				if err != nil {
					hlog.Error(userId, "db query error: ", err)
				}
				for _, current := range currents {
					users = append(users, *service.GetSafeUser(&current))
				}
				redisKey := fmt.Sprintf(constant.UserRecommendRedisKey, userId, page, constant.PageSize)
				data, err := json.Marshal(users)
				if err != nil {
					hlog.Error(userId, "json marshal error: ", err)
				}
				err = rdbRecommend.AddRecommend(ctx, redisKey, data)
				if err != nil {
					hlog.Error(userId, "redis add error: ", err)
				}
			}(userId)
		}
		wg.Wait()
	})

	go func() {
		<-taskCtx.Done()
		c.Stop()
		cancel()
	}()

	c.Start()

}
