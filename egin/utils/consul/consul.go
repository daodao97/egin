package consul

import (
	"fmt"
	"log"

	consulapi "github.com/hashicorp/consul/api"
)

// var address = "127.0.0.1:8500"
var serviceId = "egin-api"

func ConsulClient(address string) (*consulapi.Client, error) {
	config := consulapi.DefaultConfig()
	config.Address = address
	return consulapi.NewClient(config)
}

// 注册服务到consul
func ConsulRegister(address string) {
	client, err := ConsulClient(address)
	if err != nil {
		log.Fatal("consul client error : ", err)
	}

	// 创建注册到consul的服务到
	registration := new(consulapi.AgentServiceRegistration)
	registration.ID = serviceId
	registration.Name = "go-gin-api"
	registration.Port = 8080
	registration.Tags = []string{"gin"}
	registration.Address = "127.0.0.1"

	// 增加consul健康检查回调函数
	check := new(consulapi.AgentServiceCheck)
	check.HTTP = fmt.Sprintf("http://%s:%d/consul", registration.Address, registration.Port)
	check.Timeout = "5s"
	check.Interval = "5s"
	check.DeregisterCriticalServiceAfter = "30s" // 故障检查失败30s后 consul自动将注册服务删除
	registration.Check = check

	// 注册服务到consul
	err = client.Agent().ServiceRegister(registration)
}

// 取消consul注册的服务
func ConsulDeRegister(address string, serviceId string) error {
	client, err := ConsulClient(address)
	if err != nil {
		return err
	}

	err = client.Agent().ServiceDeregister(serviceId)
	if err != nil {
		return err
	}

	return nil
}

// 所有服务
func ConsulServices(address string) (services map[string]*consulapi.AgentService, err error) {
	client, err := ConsulClient(address)
	if err != nil {
		return services, err
	}
	return client.Agent().Services()
}

// 从consul中发现服务
func ConsulFindServer(address string, serviceId string) (service *consulapi.AgentService, err error) {
	client, err := ConsulClient(address)
	if err != nil {
		return service, err
	}
	// 获取指定service
	service, _, err = client.Agent().Service(serviceId, nil)
	return service, err
}

func ConsulCheckHeath(address string, serviceId string) error {
	client, err := ConsulClient(address)
	if err != nil {
		return err
	}
	// 健康检查
	_, _, err = client.Agent().AgentHealthServiceByID(serviceId)
	if err != nil {
		return err
	}
	return nil
}

func ConsulKV(address string) (kv *consulapi.KV, err error) {
	client, err := ConsulClient(address)
	if err != nil {
		return kv, err
	}
	kv = client.KV()
	return kv, nil
}

func ConsulKVPut(address string, key string, value string) error {
	client, err := ConsulClient(address)
	if err != nil {
		return err
	}
	_, err = client.KV().Put(&consulapi.KVPair{Key: key, Flags: 0, Value: []byte(value)}, nil)
	if err != nil {
		return err
	}
	return nil
}

// client.KV().Get(key, nil)
// client.KV().List(prefix, nil)
// client.KV().Keys(prefix, nil)
