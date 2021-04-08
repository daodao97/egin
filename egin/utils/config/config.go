package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/joho/godotenv"

	"github.com/daodao97/egin/egin/lib"
	"github.com/daodao97/egin/egin/utils/consul"
)

type ConfigStruct struct {
	Name     string
	Address  string
	Mode     string
	Custom   interface{}
	Database Databases
	Redis    map[string]Redis
	Mongo    map[string]Mongo
	Logger   LoggerStruct
	Lan      string
	Auth     struct {
		Cors struct {
			Enable           bool
			AllowOrigins     []string // 允许源列表
			AllowMethods     []string // 允许的方法列表
			AllowHeaders     []string // 允许的头部信息
			AllowCredentials bool     // 允许暴露请求的响应
		}
		IpAuth struct {
			Enable        bool
			AllowedIpList []string
		}
		IpLimiter struct {
			Enable  bool
			IPLimit map[string]int
		}
		AKSK struct {
			Enable  bool
			Allowed map[string]string
		}
	}
	Jwt struct {
		Secret      string
		TokenExpire int64
		OpenApi     []string
	}
	RabbitMQ map[string]Rabbitmq
	Kafka    map[string]string
	Nsq      map[string]struct {
		LookupAddress string
		NsqAddress    []string
	}
	Consul string
	Oss    map[string]OssConf
	WXConf WorkWechatConf
}

type OssConf struct {
	Endpoint        string
	AccessKeyId     string
	AccessKeySecret string
	Host            string
}

type Database struct {
	Host     string
	Port     int
	User     string
	Passwd   string
	Database string
	Driver   string
	Options  map[string]string
	Pool     struct {
		MaxOpenConns int
		MaxIdleConns int
	}
}

type Databases map[string]Database

type LoggerStruct struct {
	Type      string // stdout|file
	FileName  string
	Formatter string
	Level     int // 0 PanicLevel 5 InfoLevel 6 DebugLevel
}

type Redis struct {
	Host     string
	Port     int
	DB       int
	Password string
}

type Rabbitmq struct {
	Host   string
	Port   int
	User   string
	Passwd string
	Vhost  string
}

type Mongo struct {
	Url string
}

type WorkWechatConf struct {
	AgentId int
	Secret  string
	CorpId  string
}

var Config ConfigStruct

var defaultConfig = `
{
    "address":"127.0.0.1:8080",
    "mode":"debug",
    "logger":{
        "type":"stdout",
        "fileName":"tmp/egin_app.log",
        "level":5
    }
}
`

func init() {
	if err := godotenv.Load(".env"); err != nil {
		// log.Printf("load .env fail: %s", err)
	}

	// TODO 支持命令行参数
	data, err := ioutil.ReadFile("app.json")

	if err != nil {
		// log.Printf("load app.json fail: %s, will use default config", err)
	}

	str := string(data)

	if str == "" {
		str = defaultConfig
	}

	re, _ := regexp.Compile("<.*>")

	all := re.FindAllString(str, -1)

	for i := range all {
		s := all[i]
		factory := lib.String{Str: s}
		r := os.Getenv(factory.TrimLeft("<").TrimRight(">").Done())
		str = strings.Replace(str, s, r, -1)
	}

	err = json.Unmarshal([]byte(str), &Config)
	if err != nil {
		return
	}

	kv, err := consul.ConsulKV(Config.Consul)
	if err != nil {
		return
	}
	// 远程配置只能覆盖式的, 不支持删除某个配置
	remoteConfKey := fmt.Sprintf("%s/%s", Config.Name, Config.Mode)
	kp, _, err := kv.Get(remoteConfKey, nil)
	if err != nil {
		return
	}
	err = json.Unmarshal(kp.Value, &Config)
	if err != nil {
		return
	}
	// 配置校验 零值, 完整度
	go func() {
		for range time.Tick(time.Second * 2) {
			kp, _, err := kv.Get(remoteConfKey, nil)
			if err != nil {
				log.Fatal(err)
				return
			}
			err = json.Unmarshal(kp.Value, &Config)
			fmt.Println(Config)
		}
	}()
}
