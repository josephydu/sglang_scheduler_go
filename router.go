package main

import (
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	// 注册路由
	r.POST("/register_nodes", registerNodes)
	r.POST("/generate", generate)
	r.POST("/v1/completions", v1Completions)
	r.GET("/get_model_info", getModelInfo)

	return r
}
