package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	defaultTimeout    = 10 * time.Second
	defaultMaxRetries = 1
	maxErrorBodyBytes = 1024
)

// MasterClient 封装了与 Master 服务的 HTTP 通信逻辑。
type MasterClient struct {
	baseURL    string
	apiToken   string
	httpClient *http.Client
	maxRetries int
}

// NewMasterClient 创建一个配置了默认超时与重试策略的客户端实例。
func NewMasterClient(baseURL, apiToken string) *MasterClient {
	return &MasterClient{
		baseURL:    strings.TrimRight(strings.TrimSpace(baseURL), "/"),
		apiToken:   strings.TrimSpace(apiToken),
		httpClient: &http.Client{Timeout: defaultTimeout},
		maxRetries: defaultMaxRetries,
	}
}

// SetHTTPClient 允许在测试或自定义环境中注入自定义 http.Client。
func (c *MasterClient) SetHTTPClient(client *http.Client) {
	if client == nil {
		return
	}
	c.httpClient = client
}

// SetMaxRetries 设置网络错误时的最大重试次数。
func (c *MasterClient) SetMaxRetries(max int) {
	if max < 0 {
		max = 0
	}
	c.maxRetries = max
}

// Get 发送 GET 请求并返回响应体。
func (c *MasterClient) Get(endpoint string) ([]byte, error) {
	return c.doRequest(http.MethodGet, endpoint, nil)
}

// Post 发送带 JSON Body 的 POST 请求并返回响应体。
func (c *MasterClient) Post(endpoint string, data interface{}) ([]byte, error) {
	var payload []byte
	var err error
	if data != nil {
		payload, err = json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
	}
	return c.doRequest(http.MethodPost, endpoint, payload)
}

func (c *MasterClient) doRequest(method, endpoint string, payload []byte) ([]byte, error) {
	if c.httpClient == nil {
		c.httpClient = &http.Client{Timeout: defaultTimeout}
	}

	url := c.buildURL(endpoint)
	if url == "" {
		return nil, fmt.Errorf("baseURL is not configured")
	}

	retries := c.maxRetries
	if retries < 0 {
		retries = 0
	}

	var lastErr error
	for attempt := 0; attempt <= retries; attempt++ {
		var body io.Reader
		if payload != nil {
			body = bytes.NewReader(payload)
		}

		req, err := http.NewRequest(method, url, body)
		if err != nil {
			return nil, fmt.Errorf("create request: %w", err)
		}
		req.Header.Set("Accept", "application/json")
		if payload != nil {
			req.Header.Set("Content-Type", "application/json")
		}
		c.addAuthHeader(req)

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			if attempt < retries {
				time.Sleep(time.Duration(attempt+1) * 200 * time.Millisecond)
				continue
			}
			return nil, fmt.Errorf("request %s %s failed: %w", method, url, lastErr)
		}

		data, err := c.handleResponse(resp)
		if err != nil {
			return nil, err
		}
		return data, nil
	}

	if lastErr == nil {
		lastErr = fmt.Errorf("unknown request error")
	}
	return nil, lastErr
}

func (c *MasterClient) buildURL(endpoint string) string {
	endpoint = strings.TrimSpace(endpoint)
	if strings.HasPrefix(endpoint, "http://") || strings.HasPrefix(endpoint, "https://") {
		return endpoint
	}

	if c.baseURL == "" {
		return ""
	}

	if !strings.HasPrefix(endpoint, "/") {
		endpoint = "/" + endpoint
	}

	return c.baseURL + endpoint
}

func (c *MasterClient) addAuthHeader(req *http.Request) {
	if req == nil || c.apiToken == "" {
		return
	}
	req.Header.Set("X-API-Token", c.apiToken)
}

func (c *MasterClient) handleResponse(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()

	data, err := io.ReadAll(io.LimitReader(resp.Body, maxErrorBodyBytes))
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &HTTPError{StatusCode: resp.StatusCode, Body: string(data)}
	}

	return data, nil
}

// HTTPError 描述 HTTP 状态码异常的错误。
type HTTPError struct {
	StatusCode int
	Body       string
}

func (e *HTTPError) Error() string {
	if e == nil {
		return "http error"
	}
	return fmt.Sprintf("http request failed: status=%d body=%s", e.StatusCode, strings.TrimSpace(e.Body))
}
