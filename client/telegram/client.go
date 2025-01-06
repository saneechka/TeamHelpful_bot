package telegram

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

const apiHost = "https://api.telegram.org"

type Client struct {
	client  *http.Client
	baseURL string // Changed from host and basePath to a single baseURL
}

func NewClient(token string) *Client {
	return &Client{
		client:  &http.Client{},
		baseURL: apiHost + "/bot" + token,
	}
}

func OurBasePath(token string) string {
	return "bot" + token
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

	var response UpdatesResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}

	return response.Updates, nil
}

func (c *Client) DoRequest(method string, query url.Values) ([]byte, error) {
	// Create the full URL properly
	u, err := url.Parse(c.baseURL + "/" + method)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	u.RawQuery = query.Encode()

	request, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

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
