package utils

import (
	"fmt"
	"net/http"
	"net/url"
)

func SendTelegram(BOT_TOKEN string, CHAT_ID string, message string) error {

	apiURL := fmt.Sprintf(
		"https://api.telegram.org/bot%s/sendMessage",
		BOT_TOKEN,
	)

	data := url.Values{}
	data.Set("chat_id", CHAT_ID)
	data.Set("text", message)
	data.Set("parse_mode", "HTML") // ⬅️ WAJIB

	resp, err := http.PostForm(apiURL, data)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}

func SendTelegramPhoto(BOT_TOKEN string, CHAT_ID string, imageURL string) error {

	apiURL := fmt.Sprintf(
		"https://api.telegram.org/bot%s/sendPhoto",
		BOT_TOKEN,
	)

	data := url.Values{}
	data.Set("chat_id", CHAT_ID)
	data.Set("photo", imageURL)
	data.Set("parse_mode", "HTML") // ⬅️ WAJIB

	resp, err := http.PostForm(apiURL, data)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
