package router

import (
	"github.com/gin-gonic/gin"
	"user_web/api"
	"user_web/middleware"
)

func UmsRouter(r *gin.Engine) {
	app := api.App.UmsApi
	g := r.Group("/ums/v1").Use(middleware.Trace)
	g.GET("/roles", app.RoleListView)
	g.GET("/roles/:id", app.RoleDetailView)
	g.POST("/roles", app.RoleCreateView)
	g.PUT("/roles/:id", app.RoleUpdateView)
	g.DELETE("/roles/:id", app.RoleDeleteView)
}
