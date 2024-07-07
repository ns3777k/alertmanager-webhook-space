package space

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-resty/resty/v2"
)

type ClientSettings struct {
	BaseURL          string
	ApplicationToken string
}

type Client struct {
	client   *resty.Client
	settings *ClientSettings
}

type authResponse struct {
	AccessToken string `json:"access_token"`
}

func NewClient(settings *ClientSettings) *Client {
	client := resty.New()
	client.SetRetryCount(0)

	settings.BaseURL = strings.TrimRight(settings.BaseURL, "/")

	return &Client{client: client, settings: settings}
}

func (c *Client) SendMessage(ctx context.Context, message *Message) error {
	b, err := json.Marshal(message)
	if err != nil {
		return err
	}

	resp, err := c.client.R().
		SetBody(b).
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetAuthToken(c.settings.ApplicationToken).
		Post(c.settings.BaseURL + "/api/http/chats/messages/send-message")

	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return errors.New("sending message response code: " + resp.Status())
	}

	return nil
}
