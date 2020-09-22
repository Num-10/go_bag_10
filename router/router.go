package router

import (
	"blog_go/controller"
	"blog_go/middleware"
	"github.com/gin-gonic/gin"
)

func Router(router *gin.Engine) {
	openApi := router.Group("/oo")
	{
		openApi.GET("/reolod/:cid", controller.Index)
		openApi.GET("/list", controller.List)
		openApi.POST("/login", controller.Login)
	}

	authApi := router.Group("/ao")
	authApi.Use(middleware.Verification())
	{
		authApi.GET("/user/:id", controller.User)
	}
}
