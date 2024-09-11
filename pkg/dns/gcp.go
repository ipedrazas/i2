package dns

import (
	"context"
	"fmt"

	"google.golang.org/api/dns/v1"
	"google.golang.org/api/option"
)

// GCPProvider implements the DNSProvider interface for Google Cloud Platform
type GCPProvider struct {
	client  *dns.Service
	project string
}

func NewGCPProvider(ctx context.Context, projectID, credentialsFile string) (*GCPProvider, error) {
	client, err := dns.NewService(ctx, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		return nil, fmt.Errorf("failed to create DNS client: %v", err)
	}
	return &GCPProvider{
		client:  client,
		project: projectID,
	}, nil
}

func (p *GCPProvider) ListEntries(domain string) ([]DNSEntry, error) {
	zone, err := p.getZone(domain)
	if err != nil {
		return nil, err
	}

	records, err := p.client.ResourceRecordSets.List(p.project, zone.Name).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to list records: %v", err)
	}

	var entries []DNSEntry
	for _, record := range records.Rrsets {
		entries = append(entries, DNSEntry{
			ID:       fmt.Sprintf("gcp-%s-%s-%s", zone.Name, record.Type, record.Name),
			Domain:   domain,
			Type:     record.Type,
			Name:     record.Name,
			Content:  record.Rrdatas[0],
			TTL:      int(record.Ttl),
			Provider: "GCP",
		})
	}

	return entries, nil
}

func (p *GCPProvider) CreateRecord(domain string, record DNSRecord) error {
	zone, err := p.getZone(domain)
	if err != nil {
		return err
	}

	change := &dns.Change{
		Additions: []*dns.ResourceRecordSet{
			{
				Name:    record.Name,
				Type:    record.Type,
				Ttl:     int64(record.TTL),
				Rrdatas: []string{record.Content},
			},
		},
	}

	_, err = p.client.Changes.Create(p.project, zone.Name, change).Do()
	if err != nil {
		return fmt.Errorf("failed to create record: %v", err)
	}

	return nil
}

func (p *GCPProvider) ReadRecord(domain string, recordID string) (DNSRecord, error) {
	entries, err := p.ListEntries(domain)
	if err != nil {
		return DNSRecord{}, err
	}

	for _, entry := range entries {
		if entry.ID == recordID {
			return DNSRecord{
				Type:    entry.Type,
				Name:    entry.Name,
				Content: entry.Content,
				TTL:     entry.TTL,
			}, nil
		}
	}

	return DNSRecord{}, fmt.Errorf("record not found")
}

func (p *GCPProvider) UpdateRecord(domain string, recordID string, record DNSRecord) error {
	zone, err := p.getZone(domain)
	if err != nil {
		return err
	}

	oldRecord, err := p.ReadRecord(domain, recordID)
	if err != nil {
		return err
	}

	change := &dns.Change{
		Deletions: []*dns.ResourceRecordSet{
			{
				Name:    oldRecord.Name,
				Type:    oldRecord.Type,
				Ttl:     int64(oldRecord.TTL),
				Rrdatas: []string{oldRecord.Content},
			},
		},
		Additions: []*dns.ResourceRecordSet{
			{
				Name:    record.Name,
				Type:    record.Type,
				Ttl:     int64(record.TTL),
				Rrdatas: []string{record.Content},
			},
		},
	}

	_, err = p.client.Changes.Create(p.project, zone.Name, change).Do()
	if err != nil {
		return fmt.Errorf("failed to update record: %v", err)
	}

	return nil
}

func (p *GCPProvider) DeleteRecord(domain string, recordID string) error {
	zone, err := p.getZone(domain)
	if err != nil {
		return err
	}

	record, err := p.ReadRecord(domain, recordID)
	if err != nil {
		return err
	}

	change := &dns.Change{
		Deletions: []*dns.ResourceRecordSet{
			{
				Name:    record.Name,
				Type:    record.Type,
				Ttl:     int64(record.TTL),
				Rrdatas: []string{record.Content},
			},
		},
	}

	_, err = p.client.Changes.Create(p.project, zone.Name, change).Do()
	if err != nil {
		return fmt.Errorf("failed to delete record: %v", err)
	}

	return nil
}

func (p *GCPProvider) CheckIPUsage(ip string) (bool, error) {
	zones, err := p.client.ManagedZones.List(p.project).Do()
	if err != nil {
		return false, fmt.Errorf("failed to list zones: %v", err)
	}

	for _, zone := range zones.ManagedZones {
		records, err := p.client.ResourceRecordSets.List(p.project, zone.Name).Do()
		if err != nil {
			return false, fmt.Errorf("failed to list records: %v", err)
		}

		for _, record := range records.Rrsets {
			if record.Type == "A" && len(record.Rrdatas) > 0 && record.Rrdatas[0] == ip {
				return true, nil
			}
		}
	}

	return false, nil
}

func (p *GCPProvider) getZone(domain string) (*dns.ManagedZone, error) {
	zones, err := p.client.ManagedZones.List(p.project).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to list zones: %v", err)
	}

	for _, zone := range zones.ManagedZones {
		if zone.DnsName == domain+"." {
			return zone, nil
		}
	}

	return nil, fmt.Errorf("zone not found for domain: %s", domain)
}
