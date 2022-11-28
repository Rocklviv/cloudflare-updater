package executors

import (
	"cloudflare-dns-updater/cmd/v3/helpers/logging"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type ipify struct {
	ipifyURL   string
	httpClient *http.Client
	log        logging.Logger
}

type IPiFY interface {
	GetIP() (string, error)
}

func NewIPiFY(ipifyURL string, client *http.Client, log logging.Logger) IPiFY {
	return &ipify{
		ipifyURL:   ipifyURL,
		httpClient: client,
		log:        log,
	}
}

func (i *ipify) GetIP() (string, error) {
	resp, err := i.httpClient.Get(i.ipifyURL)
	if err != nil {
		i.log.Error(err.Error())
		return "", err
	}
	if resp.StatusCode == 200 {
		body, _ := io.ReadAll(resp.Body)
		i.log.Info(strings.TrimSuffix(fmt.Sprintln("Current IP:", string(body)), "\n"))
		return string(body), nil
	} else {
		//return fmt.Errorf("Cannot get current IP. %s return status code: %d", c.ipifyURL, resp.StatusCode)
		i.log.Error(fmt.Sprintf("cannot get current IP. %s returned status code %d", i.ipifyURL, resp.StatusCode))
		return "", errors.New(fmt.Sprintf("cannot get current IP. %s returned status code %d", i.ipifyURL, resp.StatusCode))
	}
}
