package telegram

import (
	"bytes"
	"cloudflare-dns-updater/cmd/v3/helpers/logging"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

type tg struct {
	ChatID int64
	APIKey string
	Log    logging.Logger
}

type TGMessage struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

type TGSender interface {
	ErrorMsg(msg string)
	InfoMsg(msg string)
}

func NewTG(log logging.Logger) TGSender {
	if os.Getenv("TELEGRAM_CHAT_ID") == "" {
		log.Error("ENV variable TELEGRAM_CHAT_ID is not set")
		return nil
	}
	id := os.Getenv("TELEGRAM_CHAT_ID")
	chatID, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err.Error())
		return nil
	}
	if os.Getenv("TELEGRAM_API_KEY") == "" {
		log.Error("ENV variable TELEGRAM_API_KEY is not set")
		return nil
	}
	apiKey := os.Getenv("TELEGRAM_API_KEY")
	return &tg{
		ChatID: int64(chatID),
		APIKey: apiKey,
		Log:    log,
	}
}

func (t *tg) ErrorMsg(msg string) {
	// Generating multiline message
	msgString := fmt.Sprintf("\xE2\x9D\x8C Error occured \n\r \n\r %s", msg)
	message := TGMessage{
		ChatID: t.ChatID,
		Text:   msgString,
	}
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.APIKey)
	payload, err := json.Marshal(message)
	if err != nil {
		t.Log.Error(err.Error())
		return
	}

	response, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		t.Log.Error(err.Error())
		return
	}
	defer func(body io.ReadCloser) {
		if err := body.Close(); err != nil {
			t.Log.Error(err.Error())
			t.Log.Error("failed to close response body")
		}
	}(response.Body)

	if response.StatusCode != http.StatusOK {
		t.Log.Error(fmt.Sprintf("failed to send successful request. Status was %q", response.Status))
		return
	}
	t.Log.Info(response.Status)
}

func (t *tg) InfoMsg(msg string) {
	// Generating multiline message
	msgString := fmt.Sprintf("\xE2\x9C\x85 Info message \n\r \n\r %s", msg)
	message := TGMessage{
		ChatID: t.ChatID,
		Text:   msgString,
	}
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.APIKey)
	payload, err := json.Marshal(message)
	if err != nil {
		t.Log.Error(err.Error())
		return
	}

	response, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		t.Log.Error(err.Error())
		return
	}
	defer func(body io.ReadCloser) {
		if err := body.Close(); err != nil {
			t.Log.Error(err.Error())
			t.Log.Error("failed to close response body")
		}
	}(response.Body)

	if response.StatusCode != http.StatusOK {
		t.Log.Error(fmt.Sprintf("failed to send successful request. Status was %q", response.Status))
		return
	}
	t.Log.Info(response.Status)
}
