package v3

import (
	"cloudflare-dns-updater/cmd/v3/executors"
	"cloudflare-dns-updater/cmd/v3/helpers/logging"
	"cloudflare-dns-updater/cmd/v3/helpers/telegram"
	"fmt"
	"net/http"
	"os"
	"strconv"
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

func prepare() {
	log.Info("Getting env variables")
	cloudFlareAPIURL = os.Getenv("CLOUDFLARE_API")
	cloudFlareAPIToken = os.Getenv("CLOUDFLARE_TOKEN")
	cloudFlareZoneID = os.Getenv("CLOUDFLARE_ZONE_ID")
	domainName = os.Getenv("DOMAIN_NAME")
	telegramAPIKey = os.Getenv("TELEGRAM_API_KEY")
	telegramBotName = os.Getenv("TELEGRAM_BOT_NAME")
	chatID, err := strconv.Atoi(os.Getenv("TELEGRAM_CHAT_ID"))
	if err != nil {
		log.Error(err.Error())
	}
	telegramChatID = int64(chatID)
	log.Info("Create http.Client")
	client := &http.Client{}
	log.Info("Creating instance of CloudFlare executor")
	// Creates instance of CloudFlare executor
	cloudFlareExecutor = executors.NewCloudFlareExecutor(cloudFlareAPIURL, cloudFlareAPIToken, cloudFlareZoneID, domainName, client, log)
	log.Info("Creating instance of IPiFY executor")
	// Creates instance of IPiFY executor
	ipifyExecutor = executors.NewIPiFY(ipifyURL, client, log)
	log.Info("Creating instance of Telegram bot message sender")
	tgSender = telegram.NewTG(telegramChatID, telegramAPIKey, log)
}

func Start() {
	log = logging.NewLogger()
	prepare()

    tgSender.SendMessage("CloudFlare DNS Updater started")
	for {
		ip := ipifyExecutor.GetIP()
		if ip == "" {
			log.Error("failed to get current Public IP")
			tgSender.SendMessage("failed to get current Public IP")
			time.Sleep(300 * time.Second)
			return
		}
		ok, err := cloudFlareExecutor.CheckForUpdates(ip)
		if err != nil {
			log.Error(err.Error())
			tgSender.SendMessage(err.Error())
            time.Sleep(300 * time.Second)
			return
		}
		if ok {
			err = cloudFlareExecutor.Update(ip)
			if err != nil {
				log.Error(err.Error())
				tgSender.SendMessage(err.Error())
                time.Sleep(300 * time.Second)
				return
			}
			tgSender.SendMessage(fmt.Sprintf("Changed IP address in A DNS entry for domain %s to IP %s", domainName, ip))
		}
		log.Info("Timeout for 300 seconds")
        time.Sleep(300 * time.Second)
	}
}
