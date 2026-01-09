package middleware

import (
	"github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/gin-gonic/gin"
	"goods_api/global"
	"net/http"
)

func CurrentLimiting(c *gin.Context) {
	// top_category_id（用于热点参数限流）
	categoryId := c.Query("top_category_id")
	// 申请Sentinel令牌，带热点参数
	entry, blockErr := api.Entry(
		global.Config.Sentinel.LimitResourceName,
		api.WithArgs(categoryId),          // 把分类ID传入，触发热点参数限流
		api.WithTrafficType(base.Inbound), // 标记为入站流量，固定值
	)

	// 3. 触发限流：返回兜底数据，不返回错误
	if blockErr != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "商品列表加载中，当前访问人数较多，请稍后重试~",
		})
		c.Abort() // 终止请求，不再向下执行
		return
	}

	//执行后续逻辑，请求结束后释放令牌
	defer entry.Exit()
	c.Next()
}
