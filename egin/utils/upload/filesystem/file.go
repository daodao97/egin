package filesystem

import (
	"errors"

	"github.com/daodao97/egin/utils/config"
	"github.com/daodao97/egin/utils/upload/aliyunoss"
	"github.com/daodao97/egin/utils/upload/base"
)

func New(name string) (base.Interface, error) {
	conf, exist := config.Config.Oss[name]
	if !exist {
		return nil, errors.New("oss name not found")
	}
	obj, err := aliyunoss.New(aliyunoss.Conf{
		Endpoint:        conf.Endpoint,
		AccessKeyId:     conf.AccessKeyId,
		AccessKeySecret: conf.AccessKeySecret,
		Host:            conf.Host,
	})
	return obj, err
}

func NewDefault() (base.Interface, error) {
	return New("default")
}
