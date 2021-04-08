package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"
	goredis "github.com/go-redis/redis/v8"

	"github.com/daodao97/egin/utils/config"
	"github.com/daodao97/egin/utils/logger"
)

var redisLogger = logger.NewLogger("redis")

var ctx = context.Background()

func New(connection string) Redis {
	if connection == "" {
		connection = "default"
	}
	_, ok := config.Config.Database[connection]
	if !ok {
		redisLogger.Error(fmt.Sprintf("redis connection %s not found", connection))
	}

	rdb, ok := getDBInPool(connection)
	if !ok {
		redisLogger.Error(fmt.Sprintf("get %s rdb failed", connection))
	}

	return &redis{
		Connection: connection,
		rdb:        rdb,
	}
}

func NewDefault() Redis {
	return New("default")
}

type Redis interface {
	Get(key string) (string, error)
	Set(key string, value interface{}, expiration int64) error
	Incr(key string) error
	PExpire(key string, expiration time.Duration) error
	HSet(key string, value ...interface{}) error
	HGet(key string, field string) (string, error)
	HDel(key string, field string) error
	Exists(key string) (int64, error)
	GetCache(key string, get func() (string, error), expiration int64) (string, error)
}

type redis struct {
	rdb        *goredis.Client
	Connection string
}

func (r redis) Get(key string) (string, error) {
	return r.rdb.Get(ctx, key).Result()
}

func (r redis) Set(key string, value interface{}, expiration int64) error {
	return r.rdb.Set(ctx, key, value, time.Duration(expiration*1000)).Err()
}

func (r redis) Incr(key string) error {
	return r.rdb.Incr(ctx, key).Err()
}

func (r redis) PExpire(key string, expiration time.Duration) error {
	return r.rdb.PExpire(ctx, key, expiration).Err()
}

// HSet accepts values in following formats:
//   - HSet("myhash", "key1", "value1", "key2", "value2")
//   - HSet("myhash", []string{"key1", "value1", "key2", "value2"})
//   - HSet("myhash", map[string]interface{}{"key1": "value1", "key2": "value2"})
func (r redis) HSet(key string, value ...interface{}) error {
	return r.rdb.HSet(ctx, key, value).Err()
}

func (r redis) HGet(key string, field string) (string, error) {
	return r.rdb.HGet(ctx, key, field).Result()
}

func (r redis) HDel(key string, field string) error {
	return r.rdb.HDel(ctx, key, field).Err()
}

func (r redis) Exists(key string) (int64, error) {
	return r.rdb.Exists(ctx, key).Result()
}

func (r redis) GetCache(key string, get func() (string, error), expiration int64) (string, error) {
	val, err := r.Get(key)
	spew.Dump(val, err, key, "--------------")
	if err != nil {
		return "", err
	}
	if val != "" {
		return val, nil
	}
	val, err = get()
	if err != nil {
		return "", err
	}
	err = r.Set(key, val, expiration)
	if err != nil {
		return "", err
	}
	return val, nil
}
