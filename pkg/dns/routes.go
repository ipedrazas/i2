package dns

import (
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func AddRoutes(api *gin.RouterGroup) {
	service := NewDNSService()
	Providers := os.Getenv("PROVIDERS")
	providers := strings.Split(Providers, ",")
	for _, provider := range providers {
		switch provider {
		case "gcp":
			service.SetGCPProvider()
		case "cloudflare":
			service.SetCloudflareProvider()
		default:
			log.Fatalf("Invalid provider: %s", provider)
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
