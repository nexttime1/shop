package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 允许的来源域
		// 为了安全，指定具体的域名，例如 "http://localhost:8080"
		// "*" 表示允许所有来源
		allowOrigins := []string{"http://localhost:8080"}

		origin := c.Request.Header.Get("Origin")

		// 检查请求来源是否在允许的列表中
		if origin != "" {
			for _, o := range allowOrigins {
				if o == "*" || o == origin {
					c.Header("Access-Control-Allow-Origin", origin)
					break
				}
			}
		}

		// 允许的请求头
		//这里要包含 JWT Token 所在的请求头
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, Token")

		// 允许的 HTTP 方法
		c.Header("Access-Control-Allow-Methods", "POST, GET, PUT, PATCH, DELETE, OPTIONS")

		// 是否允许发送 Cookie
		c.Header("Access-Control-Allow-Credentials", "true")

		// 预检请求（OPTIONS 请求）的有效期，单位为秒
		c.Header("Access-Control-Max-Age", "86400")

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		// 将请求传递给下一个中间件或处理器
		c.Next()
	}
}
