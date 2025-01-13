// ge-logs/pkg/clientlib/loglib/client.go
package loglib

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/PiccoloMondoC/sg-common/logtypes"
)

type Client struct {
	HttpClient *http.Client
	BaseURL    string // Base URL of LogAggregator service
	Token      string
	ApiKey     string
}

func NewClient(baseURL string, token string, apiKey string) *Client {
	return &Client{
		HttpClient: &http.Client{
			Timeout: time.Second * 10,
		},
		BaseURL: baseURL,
		Token:   token,
		ApiKey:  apiKey,
	}
}

func (c *Client) AggregateLogs(ctx context.Context, logEntry logtypes.LogEntry) error {
	// Prepare the request body
	body, err := json.Marshal(logEntry)
	if err != nil {
		return err
	}

	// Prepare the request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/logs", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("X-API-Key", c.ApiKey)

	// Send the request
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Log this error in your service's logs
		return fmt.Errorf("log aggregator returned status: %d", resp.StatusCode)
	}

	return nil
}
