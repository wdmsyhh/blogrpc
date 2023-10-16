package main

import (
	"fmt"
	"github.com/hashicorp/consul/api"
)

var (
	consulIp    = "172.18.0.1"
	consulPort  = 8500
	serviceIp   = "172.18.0.1"
	servicePort = 8080
)

func main() {
	//err := Register(serviceIp, servicePort, "user-web", []string{"web_1", "web_2"}, "user-web")
	//fmt.Println(err)
	//AllServices()
	FilterService()
}

func Register(address string, port int, name string, tags []string, id string) error {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", consulIp, consulPort)

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	// 生成对应的检查对象
	// 本机测试如果 consul 是用 docker 起的，那么这里如果写 127.0.0.1:8080 健康检查失败，因为 docker 的 consul 使用的是 docker 的 127.0.0.1，所以这里应该写服务的实际 ip
	check := &api.AgentServiceCheck{
		HTTP:                           fmt.Sprintf("http://%s:%d/health", serviceIp, servicePort),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "10s",
	}

	//生成注册对象
	registration := new(api.AgentServiceRegistration)
	registration.Name = name
	registration.ID = id
	registration.Port = port
	registration.Tags = tags
	registration.Address = address
	registration.Check = check

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}
	return nil
}

func AllServices() {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", consulIp, consulPort)

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	data, err := client.Agent().Services()
	if err != nil {
		panic(err)
	}
	for key, _ := range data {
		fmt.Println(key)
	}
}

func FilterService() {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", consulIp, consulPort)

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	data, err := client.Agent().ServicesWithFilter(`Service == "user-web"`)
	if err != nil {
		panic(err)
	}
	for key, _ := range data {
		fmt.Println(key)
	}
}
