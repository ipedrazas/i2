package dns

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// DNSProvider interface defines methods that each cloud provider must implement
type DNSProvider interface {
	ListEntries(domain string) ([]DNSEntry, error)
	CreateRecord(domain string, record DNSRecord) error
	ReadRecord(domain string, recordID string) (DNSRecord, error)
	UpdateRecord(domain string, recordID string, record DNSRecord) error
	DeleteRecord(domain string, recordID string) error
	CheckIPUsage(ip string) (bool, error)
}

// DNSEntry represents a DNS entry
type DNSEntry struct {
	ID       string `json:"id"`
	Domain   string `json:"domain"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	Content  string `json:"content"`
	TTL      int    `json:"ttl"`
	Provider string `json:"provider"`
}

// DNSRecord represents the structure for creating or updating a DNS record
type DNSRecord struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	TTL     int    `json:"ttl"`
}

// DNSService manages multiple DNS providers
type DNSService struct {
	providers map[string]DNSProvider
}

func NewDNSService() *DNSService {
	return &DNSService{
		providers: make(map[string]DNSProvider),
	}
}

func (s *DNSService) AddProvider(name string, provider DNSProvider) {
	s.providers[name] = provider
}

// ListEntriesHandler godoc
// @Summary      List DNS entries
// @Accept		 json
// @Produce      json
// @Success      200  {object}  dns.DNSEntry
// @Failure      500  {object}	interface{}
// @Router       /dns/:zone/entries [get]
func (s *DNSService) ListEntriesHandler(c *gin.Context) {
	domain := c.Param("domain")

	var allEntries []DNSEntry
	for _, provider := range s.providers {
		entries, err := provider.ListEntries(domain)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error listing entries: %v", err)})
			return
		}
		allEntries = append(allEntries, entries...)
	}

	c.JSON(http.StatusOK, allEntries)
}

// CreateRecordHandler godoc
// @Summary      Create a DNS record
// @Accept		 json
// @Produce      json
// @Success      200  {object}  dns.DNSRecord
// @Failure      500  {object}	interface{}
// @Router       /dns/:zone/records [post]
func (s *DNSService) CreateRecordHandler(c *gin.Context) {
	domain := c.Param("domain")

	var record DNSRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Error decoding request body: %v", err)})
		return
	}

	for _, provider := range s.providers {
		err := provider.CreateRecord(domain, record)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error creating record: %v", err)})
			return
		}
		break
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Record created successfully"})
}

// ReadRecordHandler godoc
// @Summary      Read a DNS record
// @Accept		 json
// @Produce      json
// @Success      200  {object}  dns.DNSRecord
// @Failure      500  {object}	interface{}
// @Router       /dns/:zone/records/:id [get]
func (s *DNSService) ReadRecordHandler(c *gin.Context) {
	domain := c.Param("domain")
	id := c.Param("id")

	for _, provider := range s.providers {
		record, err := provider.ReadRecord(domain, id)
		if err == nil {
			c.JSON(http.StatusOK, record)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
}

// UpdateRecordHandler godoc
// @Summary      Update a DNS record
// @Accept		 json
// @Produce      json
// @Success      200  {object}  dns.DNSRecord
// @Failure      500  {object}	interface{}
// @Router       /dns/:zone/records/:id [put]
func (s *DNSService) UpdateRecordHandler(c *gin.Context) {
	domain := c.Param("domain")
	id := c.Param("id")

	var record DNSRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Error decoding request body: %v", err)})
		return
	}

	for _, provider := range s.providers {
		err := provider.UpdateRecord(domain, id, record)
		if err == nil {
			c.JSON(http.StatusOK, gin.H{"message": "Record updated successfully"})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Record not found or error updating record"})
}

// DeleteRecordHandler godoc
// @Summary      Delete a DNS record
// @Accept		 json
// @Produce      json
// @Success      200  {object}  dns.DNSRecord
// @Failure      500  {object}	interface{}
// @Router       /dns/:zone/records/:id [delete]
func (s *DNSService) DeleteRecordHandler(c *gin.Context) {
	domain := c.Param("domain")
	id := c.Param("id")

	for _, provider := range s.providers {
		err := provider.DeleteRecord(domain, id)
		if err == nil {
			c.JSON(http.StatusOK, gin.H{"message": "Record deleted successfully"})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Record not found or error deleting record"})
}

// CheckIPUsageHandler godoc
// @Summary      Check if an IP is in use
// @Accept		 json
// @Produce      json
// @Success      200  {object}  dns.DNSRecord
// @Failure      500  {object}	interface{}
// @Router       /dns/:zone/records/:id [delete]
func (s *DNSService) CheckIPUsageHandler(c *gin.Context) {
	ip := c.Param("ip")

	for _, provider := range s.providers {
		used, err := provider.CheckIPUsage(ip)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error checking IP usage: %v", err)})
			return
		}
		if used {
			c.JSON(http.StatusOK, gin.H{"message": "IP is in use"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "IP is not in use"})
}

func (s *DNSService) setGCPProvider() {
	ctx := context.Background()
	projectId := os.Getenv("GCP_PROJECT_ID")
	credentialsPath := os.Getenv("GCP_CREDENTIALS_PATH")
	gcpProvider, err := NewGCPProvider(ctx, projectId, credentialsPath)

	if err != nil {
		log.Fatalf("Failed to create GCP provider: %v", err)
	}
	s.AddProvider("GCP", gcpProvider)
}

func (s *DNSService) setCloudflareProvider() {
	apiToken := os.Getenv("CF_API_TOKEN")
	cloudflareProvider, err := NewCloudflareProvider(apiToken)

	if err != nil {
		log.Fatalf("Failed to create Cloudflare provider: %v", err)
	}
	s.AddProvider("Cloudflare", cloudflareProvider)
}
