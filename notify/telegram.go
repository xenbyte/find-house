package notify

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type TelegramResponse struct {
	Ok     bool `json:"ok"`
	Result []struct {
		UpdateID int `json:"update_id"`
		Message  struct {
			MessageID int `json:"message_id"`
			From      struct {
				ID           int    `json:"id"`
				IsBot        bool   `json:"is_bot"`
				FirstName    string `json:"first_name"`
				Username     string `json:"username"`
				LanguageCode string `json:"language_code"`
				IsPremium    bool   `json:"is_premium"`
			} `json:"from"`
			Chat struct {
				ID        int    `json:"id"`
				FirstName string `json:"first_name"`
				Username  string `json:"username"`
				Type      string `json:"type"`
			} `json:"chat"`
			Date int    `json:"date"`
			Text string `json:"text"`
		} `json:"message"`
	} `json:"result"`
}

func GetChatID(token, username string) (string, error) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates", token)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var telegramResponse TelegramResponse
	if err := json.NewDecoder(resp.Body).Decode(&telegramResponse); err != nil {
		return "", err
	}

	for _, update := range telegramResponse.Result {
		if update.Message.Chat.Username == username {
			return strconv.Itoa(update.Message.Chat.ID), nil
		}
	}
	return "", fmt.Errorf("username %s not found", username)
}

func SendTelegramMessage(token, chatID, message string) error {
	telegramAPI := "https://api.telegram.org/bot" + token + "/sendMessage"
	fmt.Println("url: ", telegramAPI)
	values := url.Values{}
	values.Add("chat_id", chatID)
	values.Add("text", message)

	resp, err := http.PostForm(telegramAPI, values)
	if err != nil {
		return err
	}
	fmt.Println("status code: ", resp.StatusCode)
	defer resp.Body.Close()
	return nil
}
