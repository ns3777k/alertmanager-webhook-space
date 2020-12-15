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
	BaseURL      string
	ClientID     string
	ClientSecret string
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

func (c *Client) auth(ctx context.Context) (string, error) {
	resp, err := c.client.R().
		SetBody("grant_type=client_credentials").
		SetContext(ctx).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetBasicAuth(c.settings.ClientID, c.settings.ClientSecret).
		Post(c.settings.BaseURL + "/oauth/token")

	if err != nil {
		return "", err
	}

	defer resp.RawResponse.Body.Close()

	if resp.StatusCode() != http.StatusOK {
		return "", errors.New("auth response code: " + resp.Status())
	}

	var response authResponse
	if err := json.Unmarshal(resp.Body(), &response); err != nil {
		return "", err
	}

	return response.AccessToken, nil
}

func (c *Client) SendMessage(ctx context.Context, message *Message) error {
	token, err := c.auth(ctx)
	if err != nil {
		return err
	}

	b, err := json.Marshal(message)
	if err != nil {
		return err
	}

	resp, err := c.client.R().
		SetBody(b).
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetAuthToken(token).
		Post(c.settings.BaseURL + "/api/http/chats/messages/send-message")

	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return errors.New("sending message response code: " + resp.Status())
	}

	return nil
}
