package middleware

import (
	"context"
	"github.com/Alf_Grindel/save/internal/middleware/redis"
	"github.com/Alf_Grindel/save/pkg/utils/hlog"
	"github.com/go-redsync/redsync/v4"
	"testing"
	"time"
)

var (
	lockKey  = "test:watchdog:lock"
	lockTIme = 30 * time.Second
)

func TestWatchDog(t *testing.T) {
	Init()
	go runNode(redis.RedSync, "node-A")
	time.Sleep(1 * time.Second) // 稍微错开，防止同时抢锁
	go runNode(redis.RedSync, "node-B")

	// 主线程等待足够时间
	time.Sleep(60 * time.Second)
}

func runNode(rs *redsync.Redsync, node string) {
	for {
		ctx, cancel := context.WithCancel(context.Background())

		mu := rs.NewMutex(lockKey, redsync.WithExpiry(lockTIme), redsync.WithTries(1))

		if err := mu.LockContext(ctx); err != nil {
			hlog.Infof("[%s] Failed to acquire lock", node)
			time.Sleep(2 * time.Second)
			cancel()
			continue
		}

		hlog.Infof("[%s] Acquired lock \n", node)

		watchdogCtx, watchdogCancel := context.WithCancel(ctx)

		go func() {
			ticker := time.NewTicker(lockTIme / 2)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					ok, err := mu.Extend()
					if !ok || err != nil {
						hlog.Error("Failed to extend lock", err)
						return
					}
				case <-watchdogCtx.Done():
					return
				}
			}
		}()
		hlog.Infof("[%s] Start long task", node)
		time.Sleep(25 * time.Second)
		hlog.Infof("[%s] Task finished", node)

		watchdogCancel()
		ok, err := mu.UnlockContext(ctx)
		if !ok || err != nil {
			hlog.Errorf("[%s] Failed to unlock: %v", node, err)
		} else {
			hlog.Infof("[%s] Released lock", node)
		}
		cancel()
		break
	}
}
