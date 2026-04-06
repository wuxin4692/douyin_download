package utils

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type HTTPResponse struct {
	Body    []byte
	Headers http.Header
	Status  int
}

func HttpRequest(requestURL string, options map[string]interface{}) (*HTTPResponse, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	method := "GET"
	if m, ok := options["method"].(string); ok {
		method = m
	}

	headers := make(http.Header)
	if h, ok := options["headers"].(http.Header); ok {
		headers = h
	}

	body := io.NopCloser(strings.NewReader(""))
	if b, ok := options["body"].(string); ok && b != "" {
		body = io.NopCloser(strings.NewReader(b))
		headers.Set("Content-Type", "application/json")
	}

	req, err := http.NewRequest(method, requestURL, body)
	if err != nil {
		return nil, err
	}

	for key, values := range headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// 默认 User-Agent
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	}

	maxRedirects := 10
	for i := 0; i < maxRedirects; i++ {
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode < 300 || resp.StatusCode >= 400 {
			bodyBytes, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				return nil, err
			}
			return &HTTPResponse{
				Body:    bodyBytes,
				Headers: resp.Header,
				Status:  resp.StatusCode,
			}, nil
		}

		// 处理重定向
		location := resp.Header.Get("Location")
		resp.Body.Close()

		if location == "" {
			break
		}

		// 处理相对 URL
		if !strings.HasPrefix(location, "http") {
			base, _ := url.Parse(requestURL)
			location = base.ResolveReference(&url.URL{Path: location}).String()
		}

		req, err = http.NewRequest(method, location, nil)
		if err != nil {
			return nil, err
		}
		for key, values := range headers {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
	}

	return &HTTPResponse{
		Body:    []byte{},
		Headers: http.Header{},
		Status:  200,
	}, nil
}
