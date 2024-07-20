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
		} `json:"message,omitempty"`
		MyChatMember struct {
			Chat struct {
				ID    int64  `json:"id"`
				Title string `json:"title"`
				Type  string `json:"type"`
			} `json:"chat"`
			From struct {
				ID           int    `json:"id"`
				IsBot        bool   `json:"is_bot"`
				FirstName    string `json:"first_name"`
				Username     string `json:"username"`
				LanguageCode string `json:"language_code"`
				IsPremium    bool   `json:"is_premium"`
			} `json:"from"`
			Date          int `json:"date"`
			OldChatMember struct {
				User struct {
					ID        int64  `json:"id"`
					IsBot     bool   `json:"is_bot"`
					FirstName string `json:"first_name"`
					Username  string `json:"username"`
				} `json:"user"`
				Status string `json:"status"`
			} `json:"old_chat_member"`
			NewChatMember struct {
				User struct {
					ID        int64  `json:"id"`
					IsBot     bool   `json:"is_bot"`
					FirstName string `json:"first_name"`
					Username  string `json:"username"`
				} `json:"user"`
				Status              string `json:"status"`
				CanBeEdited         bool   `json:"can_be_edited"`
				CanManageChat       bool   `json:"can_manage_chat"`
				CanChangeInfo       bool   `json:"can_change_info"`
				CanPostMessages     bool   `json:"can_post_messages"`
				CanEditMessages     bool   `json:"can_edit_messages"`
				CanDeleteMessages   bool   `json:"can_delete_messages"`
				CanInviteUsers      bool   `json:"can_invite_users"`
				CanRestrictMembers  bool   `json:"can_restrict_members"`
				CanPromoteMembers   bool   `json:"can_promote_members"`
				CanManageVideoChats bool   `json:"can_manage_video_chats"`
				CanPostStories      bool   `json:"can_post_stories"`
				CanEditStories      bool   `json:"can_edit_stories"`
				CanDeleteStories    bool   `json:"can_delete_stories"`
				IsAnonymous         bool   `json:"is_anonymous"`
				CanManageVoiceChats bool   `json:"can_manage_voice_chats"`
			} `json:"new_chat_member"`
		} `json:"my_chat_member,omitempty"`
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

func GetChannelID(token, channelName string) (string, error) {
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

	fmt.Println("DEBUG")
	for _, update := range telegramResponse.Result {
		if update.MyChatMember.Chat.Title == channelName {
			return strconv.Itoa(int(update.MyChatMember.Chat.ID)), nil
		}
	}
	return "", fmt.Errorf("channel %s not found", channelName)
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
