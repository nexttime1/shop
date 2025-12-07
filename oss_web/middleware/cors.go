package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 允许的源（生产环境建议指定具体域名，不要用*）
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
		} else {
			c.Header("Access-Control-Allow-Origin", "*") // 本地测试可用*，生产环境替换为具体域名
		}

		// 允许的请求方法
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		// 允许的请求头
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		// 允许携带凭证（前端用了 credentials: "include" 必须设置这个）
		c.Header("Access-Control-Allow-Credentials", "true")
		// 暴露自定义响应头（如果有）
		c.Header("Access-Control-Expose-Headers", "Content-Length")
		// 预检请求缓存时间（避免频繁OPTIONS请求）
		c.Header("Access-Control-Max-Age", "86400")

		// 处理预检请求（OPTIONS）
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent) // 预检请求返回204即可，无需返回body
			return
		}

		// 继续处理请求
		c.Next()
	}
}
