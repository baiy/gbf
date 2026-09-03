package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/baiy/gbf/constant"
	"github.com/baiy/gbf/log"
	"github.com/baiy/gbf/util/cryptor"
	"github.com/baiy/gbf/util/json"
	"github.com/baiy/gbf/util/maps"
	"github.com/baiy/gbf/util/str"
)

var httpClient = &http.Client{
	Timeout: 3 * time.Hour,
	Transport: &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 10 * time.Minute,
		}).DialContext,
		TLSHandshakeTimeout:   5 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 2 * time.Second,

		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 20,
		IdleConnTimeout:     300 * time.Second,
	},
}

type Client struct {
	baseUrl string

	path string

	query url.Values
	body  string

	contentType string

	// 请求超时 默认: 30秒
	timeout time.Duration
	context context.Context

	errs []error
}

func (c *Client) Timeout(timeout time.Duration) *Client {
	c.timeout = timeout
	return c
}

func (c *Client) Path(path string, vars ...map[string]any) *Client {
	if len(vars) > 0 {
		// Replace path variables
		path = str.TemplateReplace(
			path,
			maps.MapValues(
				vars[0],
				func(key string, value any) string {
					return fmt.Sprintf("%v", value)
				},
			),
		)
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	c.path = path
	return c
}

func (c *Client) Body(d string) *Client {
	c.body = d
	return c
}

func (c *Client) ContentType(v string) *Client {
	c.contentType = v
	return c
}

func (c *Client) Queries(params map[string]any) *Client {
	for k, v := range params {
		c.query.Set(k, fmt.Sprintf("%v", v))
	}
	return c
}

func (c *Client) Query(key string, value any) *Client {
	return c.Queries(map[string]any{key: value})
}

func (c *Client) Form(params map[string]any) *Client {
	data := url.Values{}
	for k, v := range params {
		data.Set(k, fmt.Sprintf("%v", v))
	}
	return c.Body(data.Encode()).ContentType("application/x-www-form-urlencoded")
}

func (c *Client) Json[T any](body T) *Client {
	v, err := json.Encode(body)
	if err != nil {
		c.errs = append(c.errs, err)
		return c
	}
	return c.Body(v).ContentType("application/json")
}

// Get 发起GET请求
func (c *Client) Get() *Response {
	return c.doRequest(http.MethodGet)
}

// Post 发起POST请求
func (c *Client) Post() *Response {
	return c.doRequest(http.MethodPost)
}

// Put 发起PUT请求
func (c *Client) Put() *Response {
	return c.doRequest(http.MethodPut)
}

// Delete 发起DELETE请求
func (c *Client) Delete() *Response {
	return c.doRequest(http.MethodDelete)
}

// Patch 发起PATCH请求
func (c *Client) Patch() *Response {
	return c.doRequest(http.MethodPatch)
}

// doRequest 执行HTTP请求
func (c *Client) doRequest(method string) *Response {
	defer log.ElapsedTime(fmt.Sprintf("service client request: %s", c.path))()
	if len(c.errs) > 0 {
		return &Response{err: errors.Join(c.errs...)}
	}

	fullURL := fmt.Sprintf("%s%s", c.baseUrl, c.path)
	if len(c.query) > 0 {
		if strings.Contains(fullURL, "?") {
			fullURL += "&" + c.query.Encode()
		} else {
			fullURL += "?" + c.query.Encode()
		}
	}

	var body io.Reader
	if c.body != "" {
		body = strings.NewReader(c.body)
	}

	ctx, cancel := context.WithTimeout(c.context, c.timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, method, fullURL, body)
	if err != nil {
		return &Response{err: err}
	}
	if c.contentType != "" {
		req.Header.Set("Content-Type", c.contentType)
	}
	req.Header.Set("X-Request-Service-Key", cryptor.Sha256(constant.Secret))
	req.Header.Set("User-Agent", fmt.Sprintf("%s/%s %s/%s Service Request", constant.BaseName, constant.BaseVersion, constant.AppName, constant.AppVersion))

	resp, err := httpClient.Do(req)
	if err != nil {
		return &Response{err: err}
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return &Response{err: err}
	}

	if resp.StatusCode != http.StatusOK {
		return &Response{err: fmt.Errorf("status code: %d, body: %s", resp.StatusCode, string(respBody))}
	}

	return &Response{body: string(respBody)}
}

func New(baseUrl string, ctx context.Context) *Client {
	return &Client{
		baseUrl: baseUrl,
		query:   url.Values{},
		timeout: 30 * time.Second,
		context: ctx,
		path:    "/",
	}
}
