package core

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"shop_api/conf"
	"shop_api/flags"
	"shop_api/utils/nacos_get"
)

func ReadConf() *conf.Config {
	v := viper.New()
	v.SetConfigFile(flags.FileOption.File)
	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}
	var c = new(conf.Config)
	err = v.Unmarshal(&c)
	if err != nil {
		panic(fmt.Sprintf("配置文件格式错误 ,%s", err))
	}

	fmt.Printf("正在从Nacos获取配置...\n")
	fmt.Printf("Nacos地址: %s:%d\n", c.NacosInfo.Host, c.NacosInfo.Port)
	fmt.Printf("命名空间: %s\n", c.NacosInfo.Namespace)
	fmt.Printf("DataId: %s, Group: %s\n", c.NacosInfo.DataId, c.NacosInfo.Group)

	// 使用HTTP客户端直接调用Nacos API
	configContent, err := nacos_get.GetConfigWithTokenAuth(c.NacosInfo)
	if err != nil {
		panic("从Nacos获取配置失败: " + err.Error())
	}

	fmt.Printf("成功获取配置，内容长度: %d\n", len(configContent))
	fmt.Printf("配置内容预览:\n%s\n", configContent[:min(200, len(configContent))])

	// 解析配置
	err = yaml.Unmarshal([]byte(configContent), c)
	if err != nil {
		panic("解析Nacos配置失败: " + err.Error())
	}

	fmt.Printf("Nacos配置解析成功\n")

	// 配置文件监听
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		zap.S().Infof("配置文件发生变化: %s", e.Name)
	})

	return c
}

// ReadConf11 第一版
func ReadConf11() *conf.Config {

	homeDir, _ := os.UserHomeDir()
	cacheDir := filepath.Join(homeDir, ".nacos", "cache")
	logDir := filepath.Join(homeDir, ".nacos", "log")

	// 创建目录
	_ = os.MkdirAll(cacheDir, 0755)
	_ = os.MkdirAll(logDir, 0755)

	v := viper.New()
	v.SetConfigFile(flags.FileOption.File)
	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}
	var c = new(conf.Config)
	err = v.Unmarshal(&c)

	if err != nil {
		panic(fmt.Sprintf("配置文件格式错误 ,%s", err))
	}

	fmt.Printf("Nacos连接信息: host=%s, port=%d, namespace=%s, dataId=%s, group=%s\n, user= %s\n, password= %s\n,",
		c.NacosInfo.Host, c.NacosInfo.Port, c.NacosInfo.Namespace, c.NacosInfo.DataId, c.NacosInfo.Group, c.NacosInfo.User, c.NacosInfo.Password)

	// 读取 nacos 里的配置
	//  Nacos服务器配置
	serverConfigs := []constant.ServerConfig{
		{IpAddr: c.NacosInfo.Host,
			Port:        c.NacosInfo.Port,
			ContextPath: "/nacos",
			Scheme:      "http"},
	}
	// 客户端配置
	clientConfig := constant.ClientConfig{
		Username:            c.NacosInfo.User,
		Password:            c.NacosInfo.Password,
		NamespaceId:         c.NacosInfo.Namespace, // id
		TimeoutMs:           15000,                 // 超时时间5秒
		NotLoadCacheAtStart: true,                  // 启动不加载缓存
		LogDir:              logDir,
		CacheDir:            cacheDir,
		LogLevel:            "info",
	}

	client, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		zap.L().Error(err.Error())
		panic(err)

	}
	// 获得配置
	config, err := client.GetConfig(vo.ConfigParam{
		DataId: c.NacosInfo.DataId,
		Group:  c.NacosInfo.Group,
	})
	if err != nil {
		zap.L().Error(err.Error())
		// 尝试调试：列出可用的配置
		fmt.Printf("尝试列出命名空间中的配置...\n")
		page, _ := client.SearchConfig(vo.SearchConfigParam{
			Search:   "blur",
			DataId:   c.NacosInfo.DataId,
			Group:    c.NacosInfo.Group,
			PageNo:   1,
			PageSize: 10,
		})
		if page != nil {
			fmt.Printf("搜索到 %d 个配置\n", page.TotalCount)
			for _, item := range page.PageItems {
				fmt.Printf(" - DataId: %s, Group: %s\n", item.DataId, item.Group)
			}
		}

		panic(err)
	}
	err = yaml.Unmarshal([]byte(config), c)
	if err != nil {
		zap.L().Error(err.Error())
		panic(err)
	}

	fmt.Printf("读取配置文件 %s 成功\n", flags.FileOption.File)
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		zap.S().Infof("配置文件发生变化: %s", e.Name)
	})

	return c
}
