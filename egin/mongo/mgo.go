package mongo

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMgo(connection, database, collection string) *mgo {
	return &mgo{
		connection,
		database,
		collection,
		GetConnect(connection),
	}
}

type mgo struct {
	connection string
	database   string
	collection string
	db         *Database
}

// 插入单个文档
func (m *mgo) InsertOne(value interface{}) *mongo.InsertOneResult {
	collection := m.getCollection()
	insertResult, err := collection.InsertOne(context.TODO(), value)
	if err != nil {
		log.Fatal(err)
	}
	return insertResult
}

// 插入多个文档
func (m *mgo) InsertMany(values []interface{}) *mongo.InsertManyResult {
	collection := m.getCollection()
	result, err := collection.InsertMany(context.TODO(), values)
	if err != nil {
		log.Fatal(err)
	}
	return result
}

// 删除文档
func (m *mgo) Delete(key string, value interface{}) int64 {
	collection := m.getCollection()
	filter := bson.D{{key, value}}
	count, err := collection.DeleteOne(context.TODO(), filter, nil)
	if err != nil {
		log.Fatal(err)
	}
	return count.DeletedCount
}

// 删除多个文档
func (m *mgo) DeleteMany(key string, value interface{}) int64 {
	collection := m.getCollection()
	filter := bson.D{{key, value}}

	count, err := collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	return count.DeletedCount
}

// 更新单个文档
func (m *mgo) UpdateOne(filter, update interface{}) int64 {
	collection := m.getCollection()
	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	return result.UpsertedCount
}

// 更新多个文档
func (m *mgo) UpdateMany(filter, update interface{}) int64 {
	collection := m.getCollection()
	result, err := collection.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	return result.UpsertedCount
}

// 查询单个文档
func (m *mgo) FindOne(key string, value interface{}) *mongo.SingleResult {
	client := m.db.Mongo
	collection, e := client.Database(m.database).Collection(m.collection).Clone()
	if e != nil {
		log.Fatal(e)
	}
	filter := bson.D{{key, value}}
	singleResult := collection.FindOne(context.TODO(), filter)
	return singleResult
}

// 查询多个文档
func (m *mgo) FindMany(filter interface{}) (*mongo.Cursor, error) {
	client := m.db.Mongo
	collection, e := client.Database(m.database).Collection(m.collection).Clone()
	if e != nil {
		log.Fatal(e)
	}
	return collection.Find(context.TODO(), filter)
}

// 多条件查询
func (m *mgo) FindManyByFilters(filter interface{}) (*mongo.Cursor, error) {
	client := m.db.Mongo
	collection, e := client.Database(m.database).Collection(m.collection).Clone()
	if e != nil {
		log.Fatal(e)
	}
	return collection.Find(context.TODO(), bson.M{"$and": filter})
}

// 查询集合里有多少数据
func (m *mgo) CollectionCount() (string, int64) {
	collection := m.getCollection()
	name := collection.Name()
	size, _ := collection.EstimatedDocumentCount(context.TODO())
	return name, size
}

// 按选项查询集合
// Skip 跳过
// Limit 读取数量
// sort 1 ，-1 . 1 为升序 ， -1 为降序
func (m *mgo) CollectionDocuments(Skip, Limit int64, sort int, key string, value interface{}) *mongo.Cursor {
	collection := m.getCollection()
	SORT := bson.D{{"_id", sort}}
	filter := bson.D{{key, value}}
	findOptions := options.Find().SetSort(SORT).SetLimit(Limit).SetSkip(Skip)
	temp, _ := collection.Find(context.Background(), filter, findOptions)
	return temp
}

func (m *mgo) getCollection() *mongo.Collection {
	return m.db.Mongo.Database(m.database).Collection(m.collection)
}
