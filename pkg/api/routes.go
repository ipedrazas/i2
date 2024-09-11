package api

import (
	"i2/pkg/dns"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"

	// swagger embed files
	swaggerFiles "github.com/swaggo/files"
)

func AddRoutes(router *gin.Engine) {

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/healtz/ready", Ready)
	router.StaticFile("/favicon.ico", "./favicon.ico")

	// Define routes

	service := dns.NewDNSService()

	api := router.Group("/api/v1")
	api.GET("/", info)
	// /dns/:zone/entries?provider=gcp
	api.GET("/dns/:zone/entries", service.ListEntriesHandler)
	api.POST("/dns/:zone/records", service.CreateRecordHandler)
	api.GET("/dns/:zone/records/:id", service.ReadRecordHandler)
	api.PUT("/dns/:zone/records/:id", service.UpdateRecordHandler)
	api.DELETE("/dns/:zone/records/:id", service.DeleteRecordHandler)
	api.GET("/dns/ip/:ip", service.CheckIPUsageHandler)

}
