package mongo

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/daodao97/egin/egin/utils/config"
)

// D：一个BSON文档。这种类型应该在顺序重要的情况下使用，比如MongoDB命令。
// M：一张无序的map。它和D是一样的，只是它不保持顺序。
// A：一个BSON数组。
// E：D里面的一个元素。

func init() {
	config.Config.Mongo = make(map[string]config.Mongo)
	config.Config.Mongo["default"] = config.Mongo{Url: "mongodb://localhost:27017"}
	InitDb()
	_, _ = client().getCollection().DeleteMany(context.Background(), bson.D{})
}

func client() *mgo {
	return NewMgo("default", "test", "test")
}

func TestNewMgo(t *testing.T) {
	client := client()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	assert.Nil(t, client.db.Mongo.Ping(ctx, readpref.Primary()))
}

func TestMgo_InsertOne(t *testing.T) {
	client := client()
	doc := map[string]interface{}{
		"name": "张三",
		"age":  12,
	}
	result := client.InsertOne(doc)
	var a interface{}
	assert.NotEqual(t, a, result.InsertedID)
}

func TestMgo_InsertMany(t *testing.T) {
	client := client()
	docs := []interface{}{
		map[string]interface{}{
			"name": "王二",
			"age":  1,
		},
		map[string]interface{}{
			"name": "李四",
			"age":  30,
		},
	}
	result := client.InsertMany(docs)
	var a []interface{}
	assert.NotEqual(t, a, result.InsertedIDs)
}

func TestMgo_FindOne(t *testing.T) {
	client := client()
	key := "name"
	val := "张三"
	result := client.FindOne(key, val)
	var r bson.M
	assert.Equal(t, nil, result.Err())
	_ = result.Decode(&r)
	assert.Equal(t, val, r[key])
}

func TestMgo_FindMany(t *testing.T) {
	client := client()
	filter := bson.D{}
	cur, err := client.FindMany(filter)
	assert.Equal(t, nil, err)
	for cur.Next(context.TODO()) {
		var result bson.M
		err := cur.Decode(&result)
		assert.Equal(t, nil, err)
		fmt.Println(result)
	}
}

func TestMgo_FindManyByFilters(t *testing.T) {
	client := client()
	var filter = []bson.M{
		{"name": "张三"},
		{"age": bson.M{"$gte": 20}},
	}
	cur, err := client.FindManyByFilters(filter)
	assert.Equal(t, nil, err)
	i := 0
	for cur.Next(context.TODO()) {
		i++
	}
	assert.Equal(t, 0, i)
}
