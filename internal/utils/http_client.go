package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type HTTPClientConfig struct {
	Timeout     time.Duration
	AccessToken string
	ContentType string
	Accept      string
	UserAgent   string
}

func DefaultHTTPConfig() *HTTPClientConfig {
	return &HTTPClientConfig{
		Timeout:     30 * time.Second,
		ContentType: "application/json",
		Accept:      "application/json, text/plain, */*",
		UserAgent:   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
	}
}

type APIRequestConfig struct {
	Method      string
	URL         string
	Body        interface{}
	Headers     map[string]string
	AccessToken string
	Timeout     time.Duration
}

func CreateHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
	}
}

func SetCommonHeaders(req *http.Request, accessToken string) {
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Sec-Ch-Ua", `"Google Chrome";v="119", "Chromium";v="119", "Not?A_Brand";v="24"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", `"macOS"`)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36")

	if accessToken != "" {
		req.Header.Set("Authorization", accessToken)
	}
}

func SetAdditionalHeaders(req *http.Request, headers map[string]string) {
	for key, value := range headers {
		req.Header.Set(key, value)
	}
}

func MakeAPIRequest(ctx context.Context, config APIRequestConfig) (*http.Response, error) {
	var body io.Reader

	if config.Body != nil {
		bodyBytes, err := json.Marshal(config.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		body = bytes.NewBuffer(bodyBytes)
	}

	req, err := http.NewRequestWithContext(ctx, config.Method, config.URL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	SetCommonHeaders(req, config.AccessToken)

	if config.Headers != nil {
		SetAdditionalHeaders(req, config.Headers)
	}

	timeout := config.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	client := CreateHTTPClient(timeout)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return resp, nil
}

func CheckResponseStatus(resp *http.Response, acceptableStatuses ...int) error {
	if len(acceptableStatuses) == 0 {
		acceptableStatuses = []int{http.StatusOK}
	}

	for _, status := range acceptableStatuses {
		if resp.StatusCode == status {
			return nil
		}
	}

	body, _ := io.ReadAll(resp.Body)
	return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
}

func DecodeJSONResponse(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("failed to decode JSON response: %w", err)
	}

	return nil
}

func MakeJSONRequest(ctx context.Context, config APIRequestConfig, target interface{}, acceptableStatuses ...int) error {

	if config.Headers == nil {
		config.Headers = make(map[string]string)
	}
	if _, exists := config.Headers["Content-Type"]; !exists && config.Body != nil {
		config.Headers["Content-Type"] = "application/json"
	}

	resp, err := MakeAPIRequest(ctx, config)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := CheckResponseStatus(resp, acceptableStatuses...); err != nil {
		return err
	}

	if target != nil {
		return DecodeJSONResponse(resp, target)
	}

	return nil
}
