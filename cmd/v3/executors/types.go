package executors

import "net/http"

// CloudDNSResult represent API response from CloudFlareDNS API
type CloudDNSResult struct {
	Success    bool       `json:"success"`
	Errors     []Errors   `json:"errors"`
	Messages   []Messages `json:"messages"`
	ResultInfo struct{}   `json:"result_info"`
	Result     []Result   `json:"result"`
}

// Errors from API
type Errors struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Messages from API
type Messages struct{}

// Result represents response from CloudFlare API
type Result struct {
	ID         string   `json:"id"`
	ZoneID     string   `json:"zone_id"`
	ZoneName   string   `json:"zone_name"`
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	Content    string   `json:"content"`
	Proxiable  bool     `json:"proxiable"`
	Proxied    bool     `json:"proxied"`
	TTL        int32    `json:"ttl"`
	Locked     bool     `json:"locked"`
	Meta       struct{} `json:"meta"`
	CreatedOn  string   `json:"created_on"`
	ModifiedOn string   `json:"modified_on"`
}

// UpdateDNSRequest represents a API request body
type UpdateDNSRequest struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	TTL     int    `json:"ttl"`
	Proxied bool   `json:"proxied"`
}

// Configurator represents general application configuration
type Configurator struct {
	client             *http.Client
	ipifyURL           string
	ipInCloudDNS       string
	currentPublicIP    string
	cloudFlareAPI      string
	cloudFlareAPIToken string
	cloudFlareZONEID   string
	domainName         string
	domainATypeID      string
}
