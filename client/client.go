package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type SalesServiceClient struct {
	serviceURL string
	httpClient *http.Client
}

func NewSalesClient(baseURL string) *SalesServiceClient {
	return &SalesServiceClient{
		serviceURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *SalesServiceClient) Healthcheck() (*http.Response, error) {
	return c.httpClient.Get(c.serviceURL + "/healthz")
}

func (c *SalesServiceClient) CreateSale(req CreateSaleRequest) (*http.Response, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	endpoint := c.serviceURL + "/sales"
	httpReq, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	return c.httpClient.Do(httpReq)
}

func (c *SalesServiceClient) CreateSaleAndParse(req CreateSaleRequest) (*CreateSaleResponse, int, error) {
	resp, err := c.CreateSale(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, nil
	}

	var result CreateSaleResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, resp.StatusCode, err
	}

	return &result, resp.StatusCode, nil
}

func (c *SalesServiceClient) GetSales() (*http.Response, error) {
	return c.httpClient.Get(c.serviceURL + "/sales")
}

func (c *SalesServiceClient) GetSaleByID(id string) (*http.Response, error) {
	return c.httpClient.Get(fmt.Sprintf("/sales/%s", id))
}

func (c *SalesServiceClient) IsSaleActive(id string, atTime string) (*http.Response, error) {
	rawURL := fmt.Sprintf("%s/sales/%s/is-active", c.serviceURL, id)
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("не удалось распарсить URL: %w", err)
	}
	query := parsedURL.Query()
	query.Set("at", atTime)
	parsedURL.RawQuery = query.Encode()

	req, err := http.NewRequest("GET", parsedURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать HTTP-запрос: %w", err)
	}
	return c.httpClient.Do(req)
}

func (c *SalesServiceClient) IsSaleActiveAndParse(id string, targetTime string) (*IsActiveResponse, int, error) {
	resp, err := c.IsSaleActive(id, targetTime)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, nil
	}

	var result IsActiveResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, resp.StatusCode, err
	}

	return &result, resp.StatusCode, nil
}
