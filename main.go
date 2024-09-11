package main

import (
	"i2/pkg/dns"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	service := dns.NewDNSService()

	// provider := os.Getenv("PROVIDER")
	// switch provider {
	// case "GCP":
	// 	service.setGCPProvider()
	// case "CF":
	// 	service.setCloudflareProvider()
	// default:
	// 	log.Fatalf("Provider %s not supported", provider)
	// }

	r := gin.Default()

	// Define routes
	r.GET("/domains/:domain/entries", service.ListEntriesHandler)
	r.POST("/domains/:domain/records", service.CreateRecordHandler)
	r.GET("/domains/:domain/records/:id", service.ReadRecordHandler)
	r.PUT("/domains/:domain/records/:id", service.UpdateRecordHandler)
	r.DELETE("/domains/:domain/records/:id", service.DeleteRecordHandler)
	r.GET("/ip-usage/:ip", service.CheckIPUsageHandler)

	log.Fatal(r.Run(":6001"))
}
