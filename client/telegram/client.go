package telegram

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	_"path"
	"strconv"
	"crypto/tls"
)

const apiHost = "https://api.telegram.org"

type Client struct {
	client  *http.Client
	baseURL string
}

func NewClient(token string) *Client {
	
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}

	return &Client{
		client:  &http.Client{Transport: tr},
		baseURL: apiHost + "/bot" + token,
	}
}

type UpdatesResponse struct {
	Ok      bool     `json:"ok"`
	Updates []Update `json:"result"`
}

func (c *Client) Updates(offset, limit int) ([]Update, error) {
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.DoRequest("getUpdates", q)
	if err != nil {
		return nil, err
	}
//test commit
	var response UpdatesResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}

	return response.Updates, nil
}

func (c *Client) DoRequest(method string, query url.Values) ([]byte, error) {
    u := fmt.Sprintf("%s/%s", c.baseURL, method)
    
    request, err := http.NewRequest(http.MethodGet, u, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }

    request.URL.RawQuery = query.Encode()

    response, err := c.client.Do(request)
    if err != nil {
        return nil, fmt.Errorf("failed to send request: %w", err)
    }
    defer response.Body.Close()

    body, err := io.ReadAll(response.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response: %w", err)
    }

    if response.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("unexpected status code: %d, body: %s", response.StatusCode, body)
    }

    return body, nil
}

func (c *Client) SendMessage(chatID int64, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.FormatInt(chatID, 10))
	q.Add("text", text)

	_, err := c.DoRequest("sendMessage", q)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

func (c *Client) SendMessageWithKeyboard(chatID int64, text string, keyboard ReplyKeyboardMarkup) error {
    q := url.Values{}
    q.Add("chat_id", strconv.FormatInt(chatID, 10))
    q.Add("text", text)
    
    keyboardJSON, err := json.Marshal(keyboard)
    if err != nil {
        return fmt.Errorf("failed to marshal keyboard: %w", err)
    }
    q.Add("reply_markup", string(keyboardJSON))

    _, err = c.DoRequest("sendMessage", q)
    if err != nil {
        return fmt.Errorf("failed to send message with keyboard: %w", err)
    }

    return nil
}
