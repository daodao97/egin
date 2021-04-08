package redis

import (
	"encoding/json"
	"testing"

	"github.com/davecgh/go-spew/spew"

	"github.com/daodao97/egin/egin/utils/config"
)

func init() {
	config.Config.Redis = map[string]config.Redis{
		"default": {
			Host: "127.0.0.1",
			Port: 6379,
		},
	}
}

func instance() Redis {
	return New("default")
}

var key = "abc"
var value = struct {
	Id   int
	Name string
}{
	Id:   1,
	Name: "daodao",
}

func TestRedis_Set(t *testing.T) {
	str, _ := json.Marshal(value)
	err := instance().Set(key, string(str), 0)

	if err != nil {
		t.Errorf("redis set %s error", key)
		spew.Dump(err)
	}
}

func TestRedis_Get(t *testing.T) {
	val, err := instance().Get(key)
	if err != nil {
		t.Errorf("redis get %s error", key)
		spew.Dump(err)
	}

	spew.Dump(val)
}
