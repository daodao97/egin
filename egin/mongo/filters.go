package mongo

import (
	"go.mongodb.org/mongo-driver/bson"
)

type E struct {
	Key   string
	Value interface{}
}

type M bson.M

type D bson.D

type A bson.A

//
func BsonD(e ...E) bson.D {
	var be []bson.E
	for _, v := range e {
		be = append(be, bson.E{
			Key:   v.Key,
			Value: v.Value,
		})
	}
	return be
}

func Incr(val interface{}) bson.M {
	return bson.M{
		"$inc": val,
	}
}

// 无需排序的map
func Set(val interface{}) bson.M {
	return bson.M{
		"$set": val,
	}
}

func In(val interface{}) bson.M {
	return bson.M{
		"$in": val,
	}
}

func Lt(val interface{}) bson.M {
	return bson.M{
		"$lt": val,
	}
}

func Gt(val interface{}) bson.M {
	return bson.M{
		"$gt": val,
	}
}

func Match(val interface{}) bson.M {
	return bson.M{
		"$match": val,
	}
}
