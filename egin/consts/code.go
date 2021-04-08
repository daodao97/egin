package consts

import (
	"github.com/daodao97/egin/utils/config"
)

type ErrCode int

const (
	ErrorNIL    ErrCode = 404
	ErrorSystem ErrCode = 500
	ErrorParam  ErrCode = 400
)

var msgMap = map[string]map[ErrCode]string{
	"zh-CN": {
		ErrorNIL:    "位置错误",
		ErrorSystem: "服务器内部错误",
		ErrorParam:  "参数错误",
	},
	"en": {
		ErrorNIL:    "error",
		ErrorSystem: "system error",
		ErrorParam:  "param error",
	},
}

func (e ErrCode) String() string {
	lan := "zh-CN"
	if c := config.Config.Lan; c != "" {
		lan = c
	}
	msg := msgMap[lan][e]
	if msg == "" {
		return "未知错误"
	}
	return msg
}
