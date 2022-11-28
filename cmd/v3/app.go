package v3

import (
	"cloudflare-dns-updater/cmd/v3/executors"
	"cloudflare-dns-updater/cmd/v3/helpers/logging"
	"cloudflare-dns-updater/cmd/v3/helpers/notificators/telegram"
	"fmt"
	"net/http"
	"os"
	"time"
)

var (
	log                logging.Logger
	cloudFlareExecutor executors.CloudExecutor
	ipifyExecutor      executors.IPiFY

	cloudFlareAPIToken string
	cloudFlareAPIURL   = "https://api.cloudflare.com/client/v4"
	cloudFlareZoneID   string
	domainName         string
	ipifyURL           = "https://api.ipify.org"
	telegramAPIKey     string
	telegramChatID     int64
	telegramBotName    string
	tgSender           telegram.TGSender
)

func getSetVars(envVarName string) string {
	if os.Getenv(envVarName) != "" {
		return os.Getenv(envVarName)
	}
	log.Error(fmt.Sprintf("ENV variable %s not found", envVarName))
	os.Exit(1)
	return ""
}

func prepare() {
	log.Info("Getting env variables")
	cloudFlareAPIURL = getSetVars("CLOUDFLARE_API")
	cloudFlareAPIToken = getSetVars("CLOUDFLARE_TOKEN")
	cloudFlareZoneID = getSetVars("CLOUDFLARE_ZONE_ID")
	domainName = getSetVars("DOMAIN_NAME")

	log.Info("Create http.Client")
	client := &http.Client{}
	log.Info("Creating instance of CloudFlare executor")
	// Creates instance of CloudFlare executor
	cloudFlareExecutor = executors.NewCloudFlareExecutor(cloudFlareAPIURL, cloudFlareAPIToken, cloudFlareZoneID, domainName, client, log)
	log.Info("Creating instance of IPiFY executor")
	// Creates instance of IPiFY executor
	ipifyExecutor = executors.NewIPiFY(ipifyURL, client, log)
	log.Info("Creating instance of Telegram bot message sender")
	tgSender = telegram.NewTG(log)
	if tgSender == nil {
		log.Error("Failed to init notificator/Telegram messenger")
		os.Exit(1)
	}
}

func Start() {
	log = logging.NewLogger()
	prepare()
	env := os.Getenv("ENV")
	if env == "" {
		env = "prod like"
	}
	tgSender.InfoMsg(fmt.Sprintf("CloudFlare DNS Updater started - %s", env))
	for {
		ip, err := ipifyExecutor.GetIP()
		if err != nil {
			log.Error("failed to get current Public IP")
			tgSender.ErrorMsg(fmt.Sprintf("Failed to get current Public IP \n %s", err))
			time.Sleep(300 * time.Second)
			return
		}

		ok, err := cloudFlareExecutor.CheckForUpdates(ip)
		if err != nil {
			log.Error(err.Error())
			tgSender.ErrorMsg(err.Error())
			time.Sleep(300 * time.Second)
			return
		}
		if ok {
			err = cloudFlareExecutor.Update(ip)
			if err != nil {
				log.Error(err.Error())
				tgSender.ErrorMsg(err.Error())
				time.Sleep(300 * time.Second)
				return
			}
			tgSender.InfoMsg(fmt.Sprintf("Changed IP address in A DNS entry for domain %s to IP %s", domainName, ip))
		}
		log.Info("Timeout for 300 seconds")
		time.Sleep(300 * time.Second)
	}
}
