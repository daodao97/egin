package json

import (
	"encoding/json"
	"strconv"
	"strings"
)

func NewJsonObj(config string) *JsonObj {
	var conf interface{}
	err := json.Unmarshal([]byte(config), &conf)
	if err != nil {
		panic(err)
	}
	return &JsonObj{Data: conf}
}

type JsonObj struct {
	Data interface{}
}

func (j *JsonObj) Get(key string, def interface{}) interface{} {
	keyNodes := strings.Split(key, ".")
	var val interface{}
	for i, key := range keyNodes {
		if i == 0 {
			val = get(key, j.Data)
		} else {
			val = get(key, val)
		}
	}
	if val != nil {
		return val
	}
	return def
}

func (j *JsonObj) GetBool(key string, def bool) bool {
	keyNodes := strings.Split(key, ".")
	var val interface{}
	for i, key := range keyNodes {
		if i == 0 {
			val = get(key, j.Data)
		} else {
			val = get(key, val)
		}
	}
	if val, ok := val.(bool); ok {
		return val
	}
	if val, ok := val.(float64); ok {
		return int(val) != 0
	}
	if val, ok := val.(string); ok {
		return val == "true"
	}
	return def
}

func get(key string, data interface{}) interface{} {
	if v, ok := data.(map[string]interface{}); ok {
		result, ok := v[key]
		if !ok {
			return nil
		}
		return result
	}
	if v, ok := data.([]interface{}); ok {
		k, err := strconv.Atoi(key)
		if err != nil {
			return nil
		}
		return v[k]
	}
	return nil
}
