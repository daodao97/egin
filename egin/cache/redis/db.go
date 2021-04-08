package redis

import (
	"fmt"
	"sync"

	goredis "github.com/go-redis/redis/v8"

	"github.com/daodao97/egin/egin/utils/config"
)

var pool sync.Map

func init() {
	fmt.Println(22222222)
	dbConf := config.Config.Redis
	for key, conf := range dbConf {
		db := makeDb(conf)
		pool.Store(key, db)
	}
}

func makeDb(conf config.Redis) *goredis.Client {
	rdb := goredis.NewClient(&goredis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Password: conf.Password,
		DB:       conf.DB,
	})

	rdb.AddHook(&loggerHook{})

	return rdb
}

func getDBInPool(key string) (*goredis.Client, bool) {
	val, ok := pool.Load(key)
	if val == nil {
		db := makeDb(config.Config.Redis[key])
		pool.Store(key, db)
		val, ok = pool.Load(key)
	}
	return val.(*goredis.Client), ok
}
