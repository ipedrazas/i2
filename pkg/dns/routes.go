package dns

import (
	"i2/pkg/models"

	"github.com/gin-gonic/gin"
)

func AddRoutes(api *gin.RouterGroup, config *models.Config) {
	service := NewDNSService(config)

	if config.GCP != (models.GCP{}) {
		service.SetGCPProvider()
		if config.GCP.IsDefault {
			service.defaultProvider = "gcp"
		}
	}
	if config.CloudFlare != (models.CloudFlare{}) {
		service.SetCloudflareProvider()
		if config.CloudFlare.IsDefault {
			service.defaultProvider = "cloudflare"
		}
	}
	// /dns/:zone/entries?provider=gcp
	api.GET("/dns/:zone/entries", service.ListEntriesHandler)
	api.POST("/dns/:zone/records", service.CreateRecordHandler)
	api.GET("/dns/:zone/records/:id", service.ReadRecordHandler)
	api.PUT("/dns/:zone/records/:id", service.UpdateRecordHandler)
	api.DELETE("/dns/:zone/records/:id", service.DeleteRecordHandler)
	api.GET("/dns/ip/:ip", service.CheckIPUsageHandler)
}
