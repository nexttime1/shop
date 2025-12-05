package middleware

import (
	"github.com/gin-gonic/gin"
	"goods_api/common/enum"
	"goods_api/common/res"
	"goods_api/service/redis_service/redis_jwt"
	"goods_api/utils/jwts"
)

func AuthMiddleware(c *gin.Context) {
	claims, err := jwts.ParseTokenByGin(c)
	if err != nil {
		res.FailWithMsg(c, res.FailLoginCode, err.Error())
		c.Abort()
		return
	}
	//判断是否在黑名单里
	ok, blackType := redis_jwt.HasTokenBlackByGin(c)
	if ok { //ok = true  的话 在黑名单 不能再走了
		res.FailWithMsg(c, res.FailLoginCode, blackType.Msg())
		c.Abort() //后面请求响应都不走  但	c.Set  要走  所以要return
		return
	}
	c.Set("claims", claims)
	c.Next()
	return
}

func AdminMiddleware(c *gin.Context) {
	claims, err := jwts.ParseTokenByGin(c)
	if err != nil {
		res.FailWithMsg(c, res.FailLoginCode, err.Error())
		c.Abort()
		return
	}
	//判断是否在黑名单里
	ok, blackType := redis_jwt.HasTokenBlackByGin(c)
	if ok { //ok = true  的话 在黑名单 不能再走了
		res.FailWithMsg(c, res.FailLoginCode, blackType.Msg())
		c.Abort() //后面请求响应都不走  但	c.Set  要走  所以要return
		return
	}

	if claims.Role != enum.AdminRole {
		res.FailWithMsg(c, res.FailLoginCode, "权限错误")
		c.Abort()
		return
	}
	c.Set("claims", claims)
	c.Next()
}
