package v2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type CloudFlareUpdate interface {
	updateDNSRecords() error
	checkPublicIP() error
	getDNSRecords() error
	doRequest(t, url string, data []byte, headers bool) ([]byte, error)
}

// doRequest do request to CloudFlare API
func (c *Configurator) doRequest(t, url string, data []byte, headers bool) ([]byte, error) {
	req, _ := http.NewRequest(t, url, bytes.NewBuffer(data))
	if headers {
		log.Println("Added autorization header")
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.cloudFlareAPIToken))
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	return body, nil
}

// checkPublicIP check current public IP and check if it is changed
func (c *Configurator) checkPublicIP() error {
	resp, err := c.client.Get(c.ipifyURL)
	if err != nil {
		log.Println(err)
		return err
	}
	if resp.StatusCode == 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Println(strings.TrimSuffix(fmt.Sprintln("Current IP:", string(body)), "\n"))
		c.currentPublicIP = string(body)
	} else {
		return fmt.Errorf("Cannot get current IP. %s return status code: %d", c.ipifyURL, resp.StatusCode)
	}
	return nil
}

// getCloudDNSRecords get all DNS records from CloudFlare
func (c *Configurator) getDNSRecords() error {
	var results CloudDNSResult
	url := fmt.Sprintf("%s/zones/%s/%s?type=A", c.cloudFlareAPI, c.cloudFlareZONEID, "dns_records")
	log.Println(strings.TrimSuffix(fmt.Sprintln("Request URL: ", url), "\n"))

	res, err := c.doRequest("GET", url, nil, true)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	err = json.Unmarshal(res, &results)
	if err != nil {
		return fmt.Errorf("GetDNS Records. Unable to unmarshal json. %s", err)
	}

	for k := range results.Result {
		c.ipInCloudDNS = results.Result[k].Content
		c.domainATypeID = results.Result[k].ID
	}
	return nil
}

// updateDNSRecords update DNS records in CloudFlare
func (c *Configurator) updateDNSRecords() error {
	var results CloudDNSResult
	data, err := json.Marshal(UpdateDNSRequest{
		Type:    "A",
		Content: c.currentPublicIP,
		Name:    c.domainName,
		TTL:     1,
		Proxied: true,
	})

	if err != nil {
		return fmt.Errorf("Cannot marshal object")
	}

	url := fmt.Sprintf("%s/zones/%s/dns_records/%s", c.cloudFlareAPI, c.cloudFlareZONEID, c.domainATypeID)
	log.Println(strings.TrimSuffix(fmt.Sprintln("Request URL: ", url), "\n"))

	res, err := c.doRequest("PUT", url, data, true)
	err = json.Unmarshal(res, &results)
	if err != nil {
		return fmt.Errorf("Unable to update DNS record: %s", err)
	}
	if results.Success {
		return nil
	} else if !results.Success && results.Errors[0].Code != 0 {
		b, err := json.Marshal(results)
		if err != nil {
			log.Println(fmt.Errorf(err.Error()))
		}
		log.Println(fmt.Sprintf("Response from cloudflare: %s", string(b)))
		return fmt.Errorf("Nothing to update")
	}
	return nil
}

// Run main function
func Start() {
	cfg, err := setVariables()
	if err != nil {
		log.Println(err)
		return
	}
	for {
		err = cfg.getDNSRecords()
		if err != nil {
			log.Println(err)
			log.Println("Timeout: 300 seconds")
			time.Sleep(300 * time.Second)
			return
		}
		err = cfg.checkPublicIP()
		if err != nil {
			log.Println(err)
			log.Println("Timeout: 300 seconds")
			time.Sleep(300 * time.Second)
			return
		}

		if cfg.currentPublicIP != cfg.ipInCloudDNS {
			log.Println(cfg.currentPublicIP)
			log.Println(cfg.ipInCloudDNS)
			log.Println(fmt.Sprintln("Public IP will be changed to:", cfg.currentPublicIP))
			err = cfg.updateDNSRecords()
			if err != nil {
				log.Println(err)
				log.Println("Timeout: 300 seconds")
				time.Sleep(300 * time.Second)
				return
			}
		} else {
			log.Println("IP Address haven't changed. Nothing to update")
		}

		log.Println("Timeout: 300 seconds")
		time.Sleep(300 * time.Second)
	}
}

// Set variables from ENV Vars
func setVariables() (*Configurator, error) {
	log.Println("Set variables from ENV Vars")
	cfg := &Configurator{}
	cfg.cloudFlareAPI = os.Getenv("CLOUDFLARE_API")
	cfg.cloudFlareAPIToken = os.Getenv("CLOUDFLARE_TOKEN")
	cfg.cloudFlareZONEID = os.Getenv("CLOUDFLARE_ZONE_ID")
	cfg.domainName = os.Getenv("DOMAIN_NAME")
	cfg.ipifyURL = "https://api.ipify.org"
	client := &http.Client{}
	cfg.client = client

	if cfg.cloudFlareAPIToken == "" {
		log.Println("API Token is not set")
		return nil, fmt.Errorf("API Token is not set")
	}
	if cfg.cloudFlareZONEID == "" {
		log.Println("Cloud ZONE ID is not set")
		return nil, fmt.Errorf("Cloud ZONE ID is not set")
	}
	return cfg, nil
}
