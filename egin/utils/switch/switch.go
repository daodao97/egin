package _switch

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/daodao97/egin/utils/consul"
)

var (
	prefix = "/mt_logicswitch"
)

var s *Switch
var once sync.Once

type Switch struct {
	ConsulAddress string
}

type SwitchValue struct {
	Status  int         `json:"status"`
	Enabled int         `json:"enabled"`
	Filter  interface{} `json:"filter"`
}

func NewSwitch(consulAddress string) *Switch {
	once.Do(func() {
		s = &Switch{
			ConsulAddress: consulAddress,
		}
	})
	return s
}

func (s *Switch) key(path string) string {
	return fmt.Sprintf("%s/%s", prefix, path)
}

func (s *Switch) IsOn(path string) bool {
	kv, err := consul.ConsulKV(s.ConsulAddress)
	if err != nil {
		return false
	}
	pair, _, err := kv.Get(s.key(path), nil)
	if err != nil || pair == nil {
		return false
	}
	value := &SwitchValue{}
	if err := json.Unmarshal(pair.Value, value); err != nil {
		return false
	}

	// todo, filter

	return value.Enabled == 1 && value.Status == 1
}
