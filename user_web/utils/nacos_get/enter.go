package nacos_get

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"user_web/conf"
)

func GetConfigWithTokenAuth(nacosInfo conf.NacosInfo) (string, error) {
	// 1. 获取access token
	token, err := getAccessToken(nacosInfo)
	if err != nil {
		return "", fmt.Errorf("获取access token失败: %w", err)
	}

	fmt.Printf("成功获取access token: %s\n", token)

	// 2. 使用Bearer Token获取配置
	client := &http.Client{}

	// 构建请求URL（完全按照curl命令的格式）
	baseURL := fmt.Sprintf("http://%s:%d/nacos/v1/cs/configs", nacosInfo.Host, nacosInfo.Port)
	params := url.Values{}
	params.Add("dataId", nacosInfo.DataId)
	params.Add("group", nacosInfo.Group)
	params.Add("namespaceId", nacosInfo.Namespace)

	reqURL := baseURL + "?" + params.Encode()
	fmt.Printf("请求URL: %s\n", reqURL)

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return "", err
	}

	// 设置Bearer Token认证头
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	return string(body), nil
}

// 获取Nacos access token
func getAccessToken(nacosInfo conf.NacosInfo) (string, error) {
	client := &http.Client{}

	// 构建登录请求（完全按照curl命令）
	loginURL := fmt.Sprintf("http://%s:%d/nacos/v1/auth/login", nacosInfo.Host, nacosInfo.Port)
	data := url.Values{}
	data.Set("username", nacosInfo.User)
	data.Set("password", nacosInfo.Password)

	fmt.Printf("登录URL: %s\n", loginURL)
	fmt.Printf("登录参数: username=%s, password=%s\n", nacosInfo.User, nacosInfo.Password)

	req, err := http.NewRequest("POST", loginURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	// 设置Content-Type
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	fmt.Printf("登录响应状态码: %d\n", resp.StatusCode)
	fmt.Printf("登录响应体: %s\n", string(body))

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("登录失败: HTTP %d", resp.StatusCode)
	}

	// 解析响应，提取accessToken
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析登录响应失败: %w", err)
	}

	accessToken, ok := result["accessToken"].(string)
	if !ok {
		return "", fmt.Errorf("响应中未找到accessToken: %s", string(body))
	}

	return accessToken, nil
}
