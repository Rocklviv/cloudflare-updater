package executors

import (
	"bytes"
	"cloudflare-dns-updater/cmd/v3/helpers/logging"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type cloudFlare struct {
	cloudFlareAPIURL   string
	cloudFlareAPIToken string
	cloudFlareZoneID   string
	domainName         string
	httpClient         *http.Client
	// Logger
	log logging.Logger
	//
	ipInCloudDNS  string
	domainATypeID string
}

type CloudExecutor interface {
	CheckForUpdates(string) (bool, error)
	Update(string) error
}

func NewCloudFlareExecutor(apiURL, apiToken, zoneID, domainName string, client *http.Client, log logging.Logger) CloudExecutor {
	return &cloudFlare{
		cloudFlareAPIURL:   apiURL,
		cloudFlareAPIToken: apiToken,
		cloudFlareZoneID:   zoneID,
		domainName:         domainName,
		httpClient:         client,
		log:                log,
	}
}

func (cl *cloudFlare) doRequest(reqType, url string, data []byte, headers bool) ([]byte, error) {
	req, _ := http.NewRequest(reqType, url, bytes.NewBuffer(data))
	if headers {
		cl.log.Info("Added autorization header")
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", cl.cloudFlareAPIToken))
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := cl.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	return body, nil
}

func (cl *cloudFlare) CheckForUpdates(ip string) (bool, error) {
	var results CloudDNSResult
	url := fmt.Sprintf("%s/zones/%s/%s?type=A", cl.cloudFlareAPIURL, cl.cloudFlareZoneID, "dns_records")
	cl.log.Info(strings.TrimSuffix(fmt.Sprintln("Request URL: ", url), "\n"))

	res, err := cl.doRequest("GET", url, nil, true)
	if err != nil {
		return false, fmt.Errorf(err.Error())
	}

	err = json.Unmarshal(res, &results)
	if err != nil {
		return false, fmt.Errorf("GetDNS Records. Unable to unmarshal json. %s", err)
	}

	for k := range results.Result {
		cl.ipInCloudDNS = results.Result[k].Content
		cl.domainATypeID = results.Result[k].ID
	}

    cl.log.Info(fmt.Sprintf("IP in CloudFlare DNS A record %s", cl.ipInCloudDNS))
	cl.log.Info("Compare current public IP with IP in CloudFlare DNS")
	if cl.ipInCloudDNS == ip {
        cl.log.Info("IP address still same, nothing to change")
		return false, nil
	}
	return true, nil
}

func (cl *cloudFlare) Update(ip string) error {
    var results CloudDNSResult
    data, err := json.Marshal(UpdateDNSRequest{
        Type:    "A",
        Content: ip,
        Name:    cl.domainName,
        TTL:     1,
        Proxied: true,
        })

    if err != nil {
        return fmt.Errorf("Cannot marshal object")
    }

    url := fmt.Sprintf("%s/zones/%s/dns_records/%s", cl.cloudFlareAPIURL, cl.cloudFlareZoneID, cl.domainATypeID)
    cl.log.Info(strings.TrimSuffix(fmt.Sprintln("Request URL: ", url), "\n"))

    res, err := cl.doRequest("PUT", url, data, true)
    err = json.Unmarshal(res, &results)
    if err != nil {
        return fmt.Errorf("Unable to update DNS record: %s", err)
    }
    if results.Success {
        return nil
    } else if !results.Success && results.Errors[0].Code != 0 {
        b, err := json.Marshal(results)
        if err != nil {
           cl. log.Error(err.Error())
        }
        cl.log.Info(fmt.Sprintf("Response from cloudflare: %s", string(b)))
        return fmt.Errorf("Nothing to update")
    }
    return nil
}
