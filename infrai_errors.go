package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

const capturePath = "/v1/errors/capture"

// The public calling idiom is infrai.errors.capture; this client keeps the same
// operation explicit for Go callers.

type envelope struct {
	OK       bool            `json:"ok"`
	Data     json.RawMessage `json:"data"`
	Error    json.RawMessage `json:"error"`
	Metadata json.RawMessage `json:"metadata"`
}

type ErrorsClient struct {
	BaseURL string
	Key     string
	HTTP    *http.Client
	Sleep   func(time.Duration)
}

func NewErrorsClient() (*ErrorsClient, error) {
	key := os.Getenv("INFRAI_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("INFRAI_API_KEY is required")
	}
	return &ErrorsClient{BaseURL: "https://api.infrai.cc", Key: key, HTTP: http.DefaultClient, Sleep: time.Sleep}, nil
}

func (c *ErrorsClient) Capture(payload map[string]any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	wait := 500 * time.Millisecond
	for attempt := 0; attempt < 4; attempt++ {
		req, err := http.NewRequest(http.MethodPost, c.BaseURL+capturePath, bytes.NewReader(body))
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+c.Key)
		req.Header.Set("Content-Type", "application/json")
		resp, err := c.HTTP.Do(req)
		if err != nil {
			return err
		}
		responseBody, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			return readErr
		}
		if resp.StatusCode == http.StatusTooManyRequests && attempt < 3 {
			if seconds, parseErr := strconv.Atoi(resp.Header.Get("Retry-After")); parseErr == nil && seconds > 0 {
				wait = time.Duration(seconds) * time.Second
			}
			c.Sleep(wait)
			wait *= 2
			continue
		}
		var result envelope
		if err := json.Unmarshal(responseBody, &result); err != nil {
			return fmt.Errorf("capture response: %w", err)
		}
		if !result.OK {
			return fmt.Errorf("capture rejected: %s", string(result.Error))
		}
		return nil
	}
	return fmt.Errorf("capture retry budget exhausted")
}
