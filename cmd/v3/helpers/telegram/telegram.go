package telegram

import (
	"bytes"
	"cloudflare-dns-updater/cmd/v3/helpers/logging"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	SendMessage(msg string)
}

func NewTG(chatID int64, apiKey string, log logging.Logger) TGSender {
	return &tg{
		ChatID: chatID,
		APIKey: apiKey,
		Log:    log,
	}
}

func (t *tg) SendMessage(msg string) {
	message := TGMessage{
		ChatID: t.ChatID,
		Text:   msg,
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
	return
}
