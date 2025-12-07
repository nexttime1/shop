package z_test

import (
	"errors"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"strings"

	"oss_web/global"
	"time"

	"github.com/google/uuid"
)

// GetUploadTokenForBrowser 生成浏览器直传的上传凭证
// prefix: 文件前缀（如goods/images）, filename: 原始文件名
func GetUploadTokenForBrowser(prefix, filename string) (map[string]interface{}, error) {
	q := global.Config.QiNiu
	if !q.Enable {
		return nil, errors.New("七牛云上传未启用")
	}
	if q.AccessKey == "" || q.SecretKey == "" || q.Bucket == "" || q.CallbackUrl == "" {
		return nil, errors.New("七牛云配置不完整（AK/SK/Bucket/CallbackUrl）")
	}

	// 生成唯一Key（防止文件覆盖）：前缀/时间_随机串_原始文件名
	now := time.Now().Format("20060102150405")
	randomStr := uuid.New().String()[:8] // 8位随机串
	key := fmt.Sprintf("%s/%s_%s_%s", prefix, now, randomStr, filename)

	// 构建上传策略（核心：添加回调配置）
	putPolicy := storage.PutPolicy{
		Scope:            fmt.Sprintf("%s:%s", q.Bucket, key),                                                              // 指定Key，避免覆盖
		CallbackURL:      q.CallbackUrl,                                                                                    // 七牛回调后端的地址
		CallbackBody:     `{"key":"$(key)","hash":"$(etag)","fname":"$(fname)","fsize":$(fsize),"mimeType":"$(mimeType)"}`, // 回调参数（魔法变量）
		CallbackBodyType: q.CallbackBodyType,                                                                               // 回调数据格式（application/json）
		Expires:          3600,                                                                                             // 凭证有效期（1小时）
		FsizeLimit:       int64(q.Size * 1024 * 1024),                                                                      // 文件大小限制（字节）
	}

	// 生成上传Token
	mac := qbox.NewMac(q.AccessKey, q.SecretKey)
	upToken := putPolicy.UploadToken(mac)

	// 获取七牛云对应区域的上传地址
	uploadHost, err := getUploadHostByZone(q.Zone)
	if err != nil {
		return nil, err
	}

	// 返回前端所需参数
	return map[string]interface{}{
		"token":      upToken,
		"key":        key,
		"uploadHost": uploadHost,
		"cdnDomain":  q.CDN,
	}, nil
}

// getUploadHostByZone 动态获取对应区域的上传地址
func getUploadHostByZone(zone string) (string, error) {
	// 转换为 SDK 识别的 RegionID 类型
	regionID := storage.RegionID(zone)

	// 通过 SDK 获取区域信息（替代手动匹配）
	region, ok := storage.GetRegionByID(regionID)
	if !ok {
		return "", fmt.Errorf("识别区域失败")
	}

	// 优先选择 HTTPS 上传地址（更安全）
	for _, host := range region.SrcUpHosts {
		if strings.HasPrefix(host, "https") {
			return host, nil
		}
	}

	// 无 HTTPS 地址则返回第一个 HTTP 地址
	return region.SrcUpHosts[0], nil
}
