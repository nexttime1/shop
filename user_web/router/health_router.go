package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func HealthRouter(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		// 返回 200 表示服务健康
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "service is healthy",
		})
	})
}
