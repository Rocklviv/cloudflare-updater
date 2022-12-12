package executors

import (
	"cloudflare-dns-updater/cmd/v3/helpers/logging"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type MyIP struct {
	myipUrl    string
	httpClient *http.Client
	log        logging.Logger
}

type IMyIP interface {
	GetIP() (string, error)
}

func NewMyIP(ipInfoProvider string, httpClient *http.Client, log logging.Logger) IMyIP {
	var url string
	if ipInfoProvider == "ipify" {
		url = "https://api.ipify.com"
	} else if ipInfoProvider == "myip" {
		url = "https://api.myip.com"
	}
	log.Info(fmt.Sprintf("Selected IP Info provider is: %s", url))
	return &MyIP{
		myipUrl:    url,
		httpClient: httpClient,
		log:        log,
	}
}

func (m *MyIP) GetIP() (string, error) {
	var myIpResponse MYIPResponse
	resp, err := m.httpClient.Get(m.myipUrl)
	if err != nil {
		m.log.Error(err.Error())
		return "", err
	}
	if resp.StatusCode == 200 {
		body, _ := io.ReadAll(resp.Body)
		if m.isJSON(body) {
			err = json.Unmarshal(body, &myIpResponse)
			if err != nil {
				m.log.Error(err.Error())
				return "", err
			}
			m.log.Info(strings.TrimSuffix(fmt.Sprintln("Current IP:", string(myIpResponse.IP)), "\n"))
			return myIpResponse.IP, nil
		} else {
			m.log.Info(strings.TrimSuffix(fmt.Sprintln("Current IP:", string(body)), "\n"))
			return string(body), nil
		}
	} else {
		m.log.Error(fmt.Sprintf("cannot get current IP. %s returned status code %d", m.myipUrl, resp.StatusCode))
		return "", fmt.Errorf("cannot get current IP. %s returned status code %d", m.myipUrl, resp.StatusCode)
	}
}

func (m *MyIP) isJSON(body []byte) bool {
	var myIpResponse MYIPResponse
	return json.Unmarshal(body, &myIpResponse) == nil
}
