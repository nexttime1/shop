package main

import (
	"encoding/json"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"shop_api/conf"
)

func main() {
	zapLogger, _ := zap.NewDevelopment()
	defer zapLogger.Sync()
	logger := zapLogger.Sugar()

	// =================== 配置参数 ===================
	host := "192.168.163.50"
	port := uint64(8848)
	namespaceId := "b5dc39e2-0639-4fbc-a591-d84a0903381a"
	dataId := "user_web_dev.yaml"
	group := "dev"
	username := "nacos"
	password := "nacos"

	// 创建缓存和日志目录
	homeDir, _ := os.UserHomeDir()
	cacheDir := filepath.Join(homeDir, ".nacos", "cache")
	logDir := filepath.Join(homeDir, ".nacos", "log")
	_ = os.MkdirAll(cacheDir, 0755)
	_ = os.MkdirAll(logDir, 0755)

	fmt.Printf("🔧 使用参数:\n")
	fmt.Printf("   Host: %s:%d\n", host, port)
	fmt.Printf("   NamespaceId: %s\n", namespaceId)
	fmt.Printf("   DataId: %s\n", dataId)
	fmt.Printf("   Group: %s\n", group)
	fmt.Printf("   CacheDir: %s\n", cacheDir)
	fmt.Printf("   LogDir: %s\n", logDir)

	// === Step 1: 先用原生 HTTP 调试是否能访问 Nacos ===
	url := fmt.Sprintf("http://%s:%d/nacos/v1/cs/configs?dataId=%s&group=%s&namespaceId=%s",
		host, port, dataId, group, namespaceId)

	req, _ := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(username, password)
	req.Header.Set("Accept", "application/json,text/plain")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("❌ HTTP 请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("🌐 Nacos HTTP 响应 [%d]:\n%s\n", resp.StatusCode, string(body))

	if resp.StatusCode != 200 {
		fmt.Printf("🛑 状态码错误，可能是：\n")
		fmt.Printf("   - 权限不足（检查用户名密码）\n")
		fmt.Printf("   - 命名空间不存在或拼写错误\n")
		fmt.Printf("   - dataId/group 不匹配\n")
		return
	}

	// === Step 2: 使用 Nacos SDK 获取配置 ===
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      host,
			Port:        port,
			ContextPath: "nacos",
			Scheme:      "http",
		},
	}

	clientConfig := constant.ClientConfig{
		NamespaceId:         namespaceId,
		TimeoutMs:           15000,
		NotLoadCacheAtStart: true,
		CacheDir:            cacheDir,
		LogDir:              logDir,
		LogLevel:            "info",
		Username:            username,
		Password:            password,
	}

	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		})
	if err != nil {
		logger.Errorf("创建 Nacos 客户端失败: %v", err)
		return
	}

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
	})
	if err != nil {
		logger.Errorf("SDK 获取配置失败: %v", err)

		// 尝试搜索配置
		page, err := configClient.SearchConfig(vo.SearchConfigParam{
			Search:   "blur",
			DataId:   dataId,
			Group:    group,
			PageNo:   1,
			PageSize: 10,
		})
		if err != nil {
			fmt.Printf("❌ 搜索配置也失败: %v\n", err)
		} else {
			fmt.Printf("🔍 搜索结果 (%d 条):\n", page.TotalCount)
			for _, item := range page.PageItems {
				fmt.Printf(" - DataId: %s | Group: %s | ContentSize: %d\n", item.DataId, item.Group, len(item.Content))
			}
		}
		return
	}

	fmt.Printf("✅✅✅ 成功通过 SDK 获取配置！内容如下：\n%s\n", content)

	// === Step 3: 解析 YAML 到结构体 ===
	var cfg conf.Config
	err = yaml.Unmarshal([]byte(content), &cfg)
	if err != nil {
		logger.Errorf("解析 YAML 失败: %v", err)
		return
	}

	// 打印部分字段验证
	fmt.Printf("🎯 解析成功！%v,", cfg)

	// 可选：打印完整结构
	pretty, _ := json.MarshalIndent(cfg, "", "  ")
	fmt.Printf("📄 完整配置结构:\n%s\n", string(pretty))
}
