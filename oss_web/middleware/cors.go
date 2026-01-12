package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 修正允許的來源，加入前端的 8090 埠
		allowOrigins := []string{"http://localhost:8090", "http://localhost:8080", "http://127.0.0.1:8090"}

		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			for _, o := range allowOrigins {
				if o == "*" || o == origin {
					c.Header("Access-Control-Allow-Origin", origin)
					break
				}
			}
		}

		// 2. 這裡已經包含了 "Token"，是正確的
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}
