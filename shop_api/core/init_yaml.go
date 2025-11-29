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
)

func ReadConf() *conf.Config {

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
