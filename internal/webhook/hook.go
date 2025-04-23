package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

var WebHookURL string

type Embed struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Color       int    `json:"color"`
	Thumbnail   struct {
		URL string `json:"url"`
	} `json:"thumbnail"`
}

type WebhookPayload struct {
	Content string  `json:"content"`
	Embeds  []Embed `json:"embeds,omitempty"`
}

func SendInfo(content string) error {
	embed := Embed{
		Title:       "🔥 Streak Maintained!",
		Description: "Another day, another streak maintained! ✅",
		Color:       0x00ff00, // Green
	}
	embed.Thumbnail.URL = "https://tryhackme.com/img/THMlogo.png"

	err := sendWebhook("@everyone Daily streak notification from the bot: "+content, embed)
	if err != nil {
		return err
	}
	return nil
}

func SendError(content string, err error) error {
	errorEmbed := Embed{
		Title:       "❌ Error Occurred - " + content,
		Description: fmt.Sprintf("```\n%s\n```", err.Error()),
		Color:       0xff0000, // Red
	}
	er := sendWebhook("@here Something went wrong!", errorEmbed)
	if er != nil {
		return er
	}
	return nil
}

func sendWebhook(content string, embed Embed) error {
	payload := WebhookPayload{
		Content: content,
		Embeds:  []Embed{embed},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	resp, err := http.Post(WebHookURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error sending webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned non-200 status: %s", resp.Status)
	}

	return nil
}
