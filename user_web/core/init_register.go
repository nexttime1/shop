package core

import (
	"fmt"
	"github.com/hashicorp/consul/api"

	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"user_web/global"
)

type ClientRegister interface {
	Register() error
	Deregister() error
}

type ConsulRegister struct {
	ServiceID string
}

// Go 编译器对指针类型调用值接收者方法时，会自动做 *ptr 解引用
func (c ConsulRegister) Register() error {
	// 服务注册
	consulConfig := api.DefaultConfig()
	// 改一下默认  这个是 虚拟机Consul 所在的ip
	consulConfig.Address = global.Config.ConsulInfo.GetAddr()

	consulClient, err := api.NewClient(consulConfig)
	//健康检查配置 用于放到 请求的json中  告诉他如何访问  用http 去检查健康 web层
	check := &api.AgentServiceCheck{
		HTTP:                           fmt.Sprintf("http://%s:%d/health", global.Config.System.IP, global.Config.System.Port), // thhp 健康检查地址
		GRPCUseTLS:                     false,                                                                                  // 是否使用 TLS
		Interval:                       "10s",                                                                                  // 检查间隔
		Timeout:                        "5s",                                                                                   // 检查超时时间
		DeregisterCriticalServiceAfter: "10s",                                                                                  // 服务异常后多久注销
	}
	// 服务注册请求体 申请注册表  防止覆盖 id 不一样就可以

	registration := &api.AgentServiceRegistration{
		ID:      c.ServiceID,
		Name:    global.Config.ConsulInfo.Name,
		Tags:    global.Config.ConsulInfo.Tags,
		Address: global.Config.System.IP, // 告诉consul  我这个服务的ip 和 端口
		Port:    global.Config.System.Port,
		Check:   check,
	}
	//  注册服务到 Consul  发送
	err = consulClient.Agent().ServiceRegister(registration)
	if err != nil {
		zap.S().Errorf("注册服务到 Consul错误 %s", err.Error())
		return err
	}
	zap.S().Infof("%s 注册成功", global.Config.ConsulInfo.Name)

	return nil
}
func (c ConsulRegister) Deregister() error {
	// 服务注册
	consulConfig := api.DefaultConfig()
	// 改一下默认  这个是 虚拟机Consul 所在的ip
	consulConfig.Address = global.Config.ConsulInfo.GetAddr()

	consulClient, err := api.NewClient(consulConfig)

	err = consulClient.Agent().ServiceDeregister(c.ServiceID)
	if err != nil {
		zap.S().Errorf("服务注销失败 %s", err.Error())
		return err
	}
	return nil

}

func NewConsulRegister() ClientRegister {

	return &ConsulRegister{
		ServiceID: uuid.NewV4().String(),
	}
}
