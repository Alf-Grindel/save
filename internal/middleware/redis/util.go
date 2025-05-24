package redis

import (
	"context"
	"github.com/Alf_Grindel/save/pkg/utils/hlog"
	"github.com/redis/go-redis/v9"
)

// add k & v
func add(c *redis.Client, ctx context.Context, k string, v interface{}) error {
	tx := c.Pipeline()
	tx.Set(ctx, k, v, expireTime)
	_, err := tx.Exec(ctx)
	return err
}

// del k & v
func del(c *redis.Client, ctx context.Context, k string) error {
	tx := c.Pipeline()
	tx.Del(ctx, k)
	_, err := tx.Exec(ctx)
	return err
}

// exist check k is or not exist
func exist(c *redis.Client, ctx context.Context, k string) bool {
	count, err := c.Exists(ctx, k).Result()
	if err != nil {
		hlog.Error("key is not exist, ", err)
		return false
	}
	return count > 0
}

// get using k get the value
func get(c *redis.Client, ctx context.Context, k string) (string, error) {
	return c.Get(ctx, k).Result()
}
