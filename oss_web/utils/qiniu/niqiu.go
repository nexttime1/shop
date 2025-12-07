package qiniu

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"go.uber.org/zap"
	"io"
	"oss_web/common/res"
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
		zap.S().Error("七牛云上传未启用")
		return nil, errors.New("七牛云上传未启用")
	}
	if q.AccessKey == "" || q.SecretKey == "" || q.Bucket == "" || q.CallbackUrl == "" {
		zap.S().Error("七牛云配置不完整（AK/SK/Bucket/CallbackUrl）")
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
		zap.S().Error(err.Error())
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

// getUploadHostByZone 动态获取对应区域的上传地址（强制 HTTPS 标准域名）
func getUploadHostByZone(zone string) (string, error) {
	// 七牛云各区域标准 HTTPS 上传域名映射
	zoneHostMap := map[string]string{
		"z0":  "https://upload.qiniup.com",     // 华东
		"z1":  "https://upload-z1.qiniup.com",  // 华北（你的配置是 z1，对应这个）
		"z2":  "https://upload-z2.qiniup.com",  // 华南
		"na0": "https://upload-na0.qiniup.com", // 北美
		"as0": "https://upload-as0.qiniup.com", // 东南亚
	}

	host, ok := zoneHostMap[zone]
	if !ok {
		zap.S().Error("不支持的七牛云区域")
		return "", fmt.Errorf("不支持的七牛云区域：%s，可选区域：z0/z1/z2/na0/as0", zone)
	}
	return host, nil
}

// QiNiuCallback 处理七牛云上传回调
// 路由示例：POST /api/qiniu/callback
func QiNiuCallback(c *gin.Context) {
	q := global.Config.QiNiu
	// 1. 读取回调Body（读取后重置，避免后续无法读取）
	callbackBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		zap.S().Errorf("读取回调Body失败：%s", err.Error())
		res.FailWithMsg(c, res.FailServiceCode, "读取回调数据失败")
		return
	}
	// 重置Body，因为gin的Body只能读取一次
	c.Request.Body = io.NopCloser(bytes.NewReader(callbackBody))

	// 2. 解析Authorization头部（格式：QBox <AccessKey>:<Sign>）
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		zap.S().Error("缺少Authorization头部")
		res.FailWithMsg(c, res.FailArgumentCode, "缺少Authorization头部")
		return
	}
	authParts := strings.SplitN(authHeader, " ", 2)
	if len(authParts) != 2 || authParts[0] != "QBox" {
		zap.S().Error("Authorization头部格式错误")
		res.FailWithMsg(c, res.FailArgumentCode, "Authorization头部格式错误")
		return
	}
	authToken := strings.SplitN(authParts[1], ":", 2)
	if len(authToken) != 2 {
		zap.S().Error("Authorization Token格式错误")
		res.FailWithMsg(c, res.FailArgumentCode, "Authorization Token格式错误")
		return
	}
	recvAK := authToken[0] // 回调携带的AccessKey
	//recvSign := authToken[1] // 回调携带的签名

	// 3. 验证AccessKey是否匹配

	if recvAK != q.AccessKey {
		zap.S().Error("AccessKey不匹配")
		res.FailWithMsg(c, res.FailArgumentCode, "AccessKey不匹配")
		return
	}

	// 4. 计算签名（核心：用SecretKey对回调Body做HMAC-SHA1 + Base64编码）
	mac := hmac.New(sha1.New, []byte(q.SecretKey))
	mac.Write(callbackBody)
	// 计算URL安全的签名
	//calcSignURL := base64.URLEncoding.EncodeToString(mac.Sum(nil))
	//// 计算标准Base64签名（兼容七牛云）
	//calcSignStd := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	// 5. 对比签名（同时检查URL安全和标准Base64）
	//recvSignURL := strings.ReplaceAll(strings.ReplaceAll(recvSign, "+", "-"), "/", "_") // 接收签名转URL安全
	//if calcSignURL != recvSign && calcSignURL != recvSignURL && calcSignStd != recvSign {
	//	zap.S().Errorf("签名验证失败：接收签名=%s（URL安全：%s），计算签名(URL)=%s，计算签名(标准)=%s", recvSign, recvSignURL, calcSignURL, calcSignStd)
	//	res.FailWithMsg(c, http.StatusForbidden, "回调签名验证失败")
	//	return
	//}

	// 7. 解析回调数据
	var callbackData struct {
		Key      string  `json:"key"`      // 文件存储Key
		Hash     string  `json:"hash"`     // 文件哈希
		Fname    string  `json:"fname"`    // 原始文件名
		Fsize    float64 `json:"fsize"`    // 文件大小（字节）
		MimeType string  `json:"mimeType"` // 文件类型
	}
	if err := json.Unmarshal(callbackBody, &callbackData); err != nil {
		zap.S().Errorf("解析回调数据失败：%s", err.Error())
		res.FailWithMsg(c, res.FailServiceCode, "回调数据解析失败")
		return
	}

	// 优化：使用 filepath.Join 处理斜杠，再替换为 URL 分隔符
	cdnDomain := strings.TrimSuffix(q.CDN, "/")          // 去除 CDN 域名末尾的 /
	fileKey := strings.TrimPrefix(callbackData.Key, "/") // 去除 Key 开头的 /
	fileUrl := fmt.Sprintf("%s/%s", cdnDomain, fileKey)

	zap.S().Infof("文件上传成功：%s，大小：%.2fMB，类型：%s", fileUrl, callbackData.Fsize/1024/1024, callbackData.MimeType)

	// 9. 响应七牛云（必须返回200，否则七牛会重试回调）
	res.OkWithData(c, map[string]interface{}{
		"code":     res.SuccessCode,
		"msg":      "回调处理成功",
		"file_url": fileUrl,
		"key":      callbackData.Key,
	})

}
