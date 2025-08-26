package utils

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type TokenProvider interface {
	GetAccessToken() string
}

type RepositoryConfig struct {
	BaseURL       string
	TokenProvider TokenProvider
	Timeout       time.Duration
}

func DefaultRepositoryConfig(baseURL string, tokenProvider TokenProvider) *RepositoryConfig {
	return &RepositoryConfig{
		BaseURL:       baseURL,
		TokenProvider: tokenProvider,
		Timeout:       10 * time.Second,
	}
}

// Example usage:
// cfg := config.Init()
// repoConfig := DefaultRepositoryConfig("https://api.example.com", cfg)
// client := NewRepositoryClient(repoConfig)

type RepositoryClient struct {
	config *RepositoryConfig
}

func NewRepositoryClient(config *RepositoryConfig) *RepositoryClient {
	return &RepositoryClient{
		config: config,
	}
}

func (rc *RepositoryClient) Get(ctx context.Context, endpoint string, target interface{}) error {
	url := fmt.Sprintf("%s%s", rc.config.BaseURL, endpoint)

	apiConfig := APIRequestConfig{
		Method:      "GET",
		URL:         url,
		AccessToken: rc.config.TokenProvider.GetAccessToken(),
		Timeout:     rc.config.Timeout,
	}

	return MakeJSONRequest(ctx, apiConfig, target)
}

func (rc *RepositoryClient) Post(ctx context.Context, endpoint string, body interface{}, target interface{}, acceptableStatuses ...int) error {
	url := fmt.Sprintf("%s%s", rc.config.BaseURL, endpoint)

	apiConfig := APIRequestConfig{
		Method:      "POST",
		URL:         url,
		Body:        body,
		AccessToken: rc.config.TokenProvider.GetAccessToken(),
		Timeout:     rc.config.Timeout,
		Headers: map[string]string{
			"Content-Type": "application/json",
			"accept":       "application/json",
		},
	}

	if len(acceptableStatuses) == 0 {
		acceptableStatuses = []int{http.StatusOK}
	}

	return MakeJSONRequest(ctx, apiConfig, target, acceptableStatuses...)
}

func (rc *RepositoryClient) Put(ctx context.Context, endpoint string, body interface{}, target interface{}) error {
	url := fmt.Sprintf("%s%s", rc.config.BaseURL, endpoint)

	apiConfig := APIRequestConfig{
		Method:      "PUT",
		URL:         url,
		Body:        body,
		AccessToken: rc.config.TokenProvider.GetAccessToken(),
		Timeout:     rc.config.Timeout,
		Headers: map[string]string{
			"Content-Type": "application/json",
			"accept":       "application/json",
		},
	}

	return MakeJSONRequest(ctx, apiConfig, target)
}

func (rc *RepositoryClient) Delete(ctx context.Context, endpoint string) error {
	url := fmt.Sprintf("%s%s", rc.config.BaseURL, endpoint)

	apiConfig := APIRequestConfig{
		Method:      "DELETE",
		URL:         url,
		AccessToken: rc.config.TokenProvider.GetAccessToken(),
		Timeout:     rc.config.Timeout,
	}

	return MakeJSONRequest(ctx, apiConfig, nil)
}

func (rc *RepositoryClient) PostWithConflictHandling(ctx context.Context, endpoint string, body interface{}, target interface{}) (bool, error) {
	url := fmt.Sprintf("%s%s", rc.config.BaseURL, endpoint)

	apiConfig := APIRequestConfig{
		Method:      "POST",
		URL:         url,
		Body:        body,
		AccessToken: rc.config.TokenProvider.GetAccessToken(),
		Timeout:     rc.config.Timeout,
		Headers: map[string]string{
			"Content-Type": "application/json",
			"accept":       "application/json",
		},
	}

	resp, err := MakeAPIRequest(ctx, apiConfig)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	isConflict := resp.StatusCode == http.StatusConflict
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusConflict {
		return isConflict, CheckResponseStatus(resp, http.StatusOK, http.StatusConflict)
	}

	if target != nil {
		if err := DecodeJSONResponse(resp, target); err != nil {
			return isConflict, err
		}
	}

	return isConflict, nil
}

type BatchRequest struct {
	operations []func(context.Context) error
}

func NewBatchRequest() *BatchRequest {
	return &BatchRequest{
		operations: make([]func(context.Context) error, 0),
	}
}

func (br *BatchRequest) AddOperation(operation func(context.Context) error) {
	br.operations = append(br.operations, operation)
}

func (br *BatchRequest) Execute(ctx context.Context) error {
	for i, operation := range br.operations {
		if err := operation(ctx); err != nil {
			return WrapErrorf(err, "batch operation %d failed", i)
		}
	}
	return nil
}

type GenericIDResponse struct {
	ID int `json:"id"`
}

func CreateIDMapper[T any](client *RepositoryClient, endpoint string, nameMapper func(T) map[string]interface{}) func(context.Context, T) (int, error) {
	return func(ctx context.Context, item T) (int, error) {
		body := nameMapper(item)
		var response GenericIDResponse

		err := client.Post(ctx, endpoint, body, &response)
		if err != nil {
			return 0, err
		}

		return response.ID, nil
	}
}

func CreateBatchIDMapper[T any](client *RepositoryClient, endpoint string, nameMapper func(T) map[string]interface{}) func(context.Context, []T) ([]int, error) {
	singleMapper := CreateIDMapper(client, endpoint, nameMapper)

	return func(ctx context.Context, items []T) ([]int, error) {
		ids := make([]int, len(items))
		batch := NewBatchRequest()

		for i, item := range items {

			index := i
			currentItem := item

			batch.AddOperation(func(ctx context.Context) error {
				id, err := singleMapper(ctx, currentItem)
				if err != nil {
					return err
				}
				ids[index] = id
				return nil
			})
		}

		if err := batch.Execute(ctx); err != nil {
			return nil, err
		}

		return ids, nil
	}
}
