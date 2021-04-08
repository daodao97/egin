package mongo

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/daodao97/egin/utils/config"
)

type Database struct {
	Mongo *mongo.Client
}

var pool sync.Map

// 初始化
func init() {
	InitDb()
}

func InitDb() {
	list := config.Config.Mongo
	for k, v := range list {
		pool.Store(k, &Database{
			Mongo: SetConnect(v.Url),
		})
	}
}

// 连接设置
func SetConnect(mongoUrl string) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// 连接池
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUrl).SetMaxPoolSize(20))
	if err != nil {
		log.Println(err)
	}
	return client
}

// 获取连接
func GetConnect(key string) *Database {
	db, ok := pool.Load(key)
	if !ok {
		log.Println(fmt.Sprintf("mongo connect [%s] not found", key))
	}
	return db.(*Database)
}
