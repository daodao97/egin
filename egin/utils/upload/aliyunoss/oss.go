package aliyunoss

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/davecgh/go-spew/spew"

	"github.com/daodao97/egin/egin/utils/upload/base"
)

type Conf struct {
	Endpoint        string
	AccessKeyId     string
	AccessKeySecret string
	Host            string
}

func New(conf Conf) (base.Interface, error) {
	client, err := newOssClient(conf)
	if err != nil {
		return nil, err
	}

	return &aliYunOss{
		instance: client,
		conf: conf,
	}, nil
}

func newOssClient(conf Conf) (*oss.Client, error) {
	// Endpoint以杭州为例，其它Region请按实际情况填写。
	endpoint := conf.Endpoint
	// 阿里云主账号AccessKey拥有所有API的访问权限，风险很高。强烈建议您创建并使用RAM账号进行API访问或日常运维，请登录 https://ram.console.aliyun.com 创建RAM账号。
	accessKeyId := conf.AccessKeyId
	accessKeySecret := conf.AccessKeySecret
	// 创建OSSClient实例。
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	return client, err
}

type aliYunOss struct {
	instance *oss.Client
	conf Conf
}

func (a aliYunOss) Save(bucketName string, localFile string, saveFile string) (string, error) {
	bucket, err := a.instance.Bucket(bucketName)
	if err != nil {
		return "", err
	}
	err = bucket.PutObjectFromFile(saveFile, localFile)
	if err != nil {
		return "", err
	}
	info, _ := bucket.GetObjectDetailedMeta(saveFile)
	spew.Dump(info)
	return a.conf.Host + "/" + saveFile, nil
}
