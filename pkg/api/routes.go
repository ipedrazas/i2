package api

import (
	"i2/pkg/dns"
	"i2/pkg/models"
	"i2/pkg/prxmx"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"

	// swagger embed files
	swaggerFiles "github.com/swaggo/files"
)

func AddRoutes(router *gin.Engine, config *models.Config) {

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/healtz/ready", Ready)
	router.StaticFile("/favicon.ico", "./favicon.ico")

	// Define routes

	api := router.Group("/api/v1")
	api.GET("/", info)
	dns.AddRoutes(api, config)
	prxmx.AddRoutes(api, config)
}
