package fakturoid

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

const baseURL = "https://app.fakturoid.cz/api/v3"
const tokenURL = "https://app.fakturoid.cz/api/v3/oauth/token"
const userAgent = "fakturoid-mcp (github.com/tedyno/fakturoid-mcp)"

type Client struct {
	clientID     string
	clientSecret string
	slug         string
	httpClient   *http.Client

	mu          sync.Mutex
	accessToken string
	tokenExpiry time.Time
}

func NewClient(clientID, clientSecret, slug string) *Client {
	return &Client{
		clientID:     clientID,
		clientSecret: clientSecret,
		slug:         slug,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func (c *Client) authenticate() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Token still valid (with 5 min buffer)
	if c.accessToken != "" && time.Now().Add(5*time.Minute).Before(c.tokenExpiry) {
		return nil
	}

	payload, _ := json.Marshal(map[string]string{"grant_type": "client_credentials"})
	req, err := http.NewRequest("POST", tokenURL, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("create token request: %w", err)
	}
	req.SetBasicAuth(c.clientID, c.clientSecret)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read token response: %w", err)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("oauth token error (%d): %s", resp.StatusCode, string(body))
	}

	var tok tokenResponse
	if err := json.Unmarshal(body, &tok); err != nil {
		return fmt.Errorf("unmarshal token: %w", err)
	}

	c.accessToken = tok.AccessToken
	c.tokenExpiry = time.Now().Add(time.Duration(tok.ExpiresIn) * time.Second)
	return nil
}

func (c *Client) do(method, endpoint string, body any, result any) error {
	if err := c.authenticate(); err != nil {
		return err
	}

	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	fullURL := fmt.Sprintf("%s/accounts/%s%s", baseURL, c.slug, endpoint)
	req, err := http.NewRequest(method, fullURL, bodyReader)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	c.mu.Lock()
	token := c.accessToken
	c.mu.Unlock()

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode == 429 {
		return fmt.Errorf("fakturoid rate limit exceeded, try again later")
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("fakturoid API error (%d): %s", resp.StatusCode, string(respBody))
	}

	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("unmarshal response: %w", err)
		}
	}

	return nil
}
