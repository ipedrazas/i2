package dns

import (
	"context"
	"fmt"

	"github.com/cloudflare/cloudflare-go"
)

type CloudflareProvider struct {
	api *cloudflare.API
}

func NewCloudflareProvider(apiToken string) (*CloudflareProvider, error) {
	api, err := cloudflare.NewWithAPIToken(apiToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create Cloudflare client: %v", err)
	}
	return &CloudflareProvider{api: api}, nil
}

func (p *CloudflareProvider) ListEntries(domain string) ([]DNSEntry, error) {
	zoneID, err := p.getZoneID(domain)
	if err != nil {
		return nil, err
	}

	records, _, err := p.api.ListDNSRecords(context.Background(), cloudflare.ZoneIdentifier(zoneID), cloudflare.ListDNSRecordsParams{})
	if err != nil {
		return nil, fmt.Errorf("failed to list records: %v", err)
	}

	var entries []DNSEntry
	for _, record := range records {
		entries = append(entries, DNSEntry{
			ID:       record.ID,
			Domain:   domain,
			Type:     record.Type,
			Name:     record.Name,
			Content:  record.Content,
			TTL:      record.TTL,
			Provider: "Cloudflare",
		})
	}

	return entries, nil
}

func (p *CloudflareProvider) CreateRecord(domain string, record DNSRecord) error {
	zoneID, err := p.getZoneID(domain)
	if err != nil {
		return err
	}

	_, err = p.api.CreateDNSRecord(context.Background(), cloudflare.ZoneIdentifier(zoneID), cloudflare.CreateDNSRecordParams{
		Type:    record.Type,
		Name:    record.Name,
		Content: record.Content,
		TTL:     record.TTL,
	})
	if err != nil {
		return fmt.Errorf("failed to create record: %v", err)
	}

	return nil
}

func (p *CloudflareProvider) ReadRecord(domain string, recordID string) (DNSRecord, error) {
	zoneID, err := p.getZoneID(domain)
	if err != nil {
		return DNSRecord{}, err
	}

	record, err := p.api.GetDNSRecord(context.Background(), cloudflare.ZoneIdentifier(zoneID), recordID)
	if err != nil {
		return DNSRecord{}, fmt.Errorf("failed to read record: %v", err)
	}

	return DNSRecord{
		Type:    record.Type,
		Name:    record.Name,
		Content: record.Content,
		TTL:     record.TTL,
	}, nil
}

func (p *CloudflareProvider) UpdateRecord(domain string, recordID string, record DNSRecord) error {
	zoneID, err := p.getZoneID(domain)
	if err != nil {
		return err
	}

	_, err = p.api.UpdateDNSRecord(context.Background(), cloudflare.ZoneIdentifier(zoneID), cloudflare.UpdateDNSRecordParams{
		ID:      recordID,
		Type:    record.Type,
		Name:    record.Name,
		Content: record.Content,
		TTL:     record.TTL,
	})
	if err != nil {
		return fmt.Errorf("failed to update record: %v", err)
	}

	return nil
}

func (p *CloudflareProvider) DeleteRecord(domain string, recordID string) error {
	zoneID, err := p.getZoneID(domain)
	if err != nil {
		return err
	}

	err = p.api.DeleteDNSRecord(context.Background(), cloudflare.ZoneIdentifier(zoneID), recordID)
	if err != nil {
		return fmt.Errorf("failed to delete record: %v", err)
	}

	return nil
}

func (p *CloudflareProvider) CheckIPUsage(ip string) (bool, error) {
	zones, err := p.api.ListZones(context.Background())
	if err != nil {
		return false, fmt.Errorf("failed to list zones: %v", err)
	}

	for _, zone := range zones {
		records, _, err := p.api.ListDNSRecords(context.Background(), cloudflare.ZoneIdentifier(zone.ID), cloudflare.ListDNSRecordsParams{})
		if err != nil {
			return false, fmt.Errorf("failed to list records: %v", err)
		}

		for _, record := range records {
			if record.Type == "A" && record.Content == ip {
				return true, nil
			}
		}
	}

	return false, nil
}

func (p *CloudflareProvider) getZoneID(domain string) (string, error) {
	zones, err := p.api.ListZones(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to list zones: %v", err)
	}

	for _, zone := range zones {
		if zone.Name == domain {
			return zone.ID, nil
		}
	}

	return "", fmt.Errorf("zone not found for domain: %s", domain)
}
