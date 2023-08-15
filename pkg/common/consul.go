package common

import (
	"fmt"
	"github.com/hashicorp/consul/api"
)

type Register struct {
	Host string
	Port int
}
type RegisterClient interface {
	Register(address string, port int, name string, tags []string, id string) error
	DeRegister(serviceId string) error
}

func NewRegister(host string, port int) *Register {
	return &Register{
		Host: host,
		Port: port,
	}
}

func (r Register) Register(address string, port int, name string, tags []string, id string) error {
	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%d", r.Host, r.Port)
	client, err := api.NewClient(config)
	if err != nil {
		panic(err)
	}

	//consul的检查工具，检查服务健康和分布式一致的检查对象
	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", address, port),
		Timeout:                        "30s",
		Interval:                       "6s",
		DeregisterCriticalServiceAfter: "15s", // 严重状态超过此时间，自动取消注册
	}

	// 生成注册对象
	registration := &api.AgentServiceRegistration{
		Name:    name,
		ID:      id,
		Port:    port,
		Address: address,
		Tags:    tags,
		Check:   check,
	}

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}
	return nil
}

func (r Register) DeRegister(serviceId string) error {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", r.Host, r.Port)

	client, err := api.NewClient(cfg)
	if err != nil {
		return err
	}
	err = client.Agent().ServiceDeregister(serviceId)
	return err
}
