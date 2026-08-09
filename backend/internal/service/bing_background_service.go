package service

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"rustdesk-api-server-pro/internal/errcode"
	"strings"
	"sync"
	"time"
)

const defaultBingArchiveURL = "https://www.bing.com/HPImageArchive.aspx?format=js&idx=0&n=1&mkt=zh-CN"

// BingBackgroundService 获取并缓存 Bing 每日背景地址。
type BingBackgroundService struct {
	client   *http.Client
	endpoint string
	mu       sync.Mutex
	imageURL string
	expires  time.Time
}

// NewBingBackgroundService 创建 Bing 背景服务。
func NewBingBackgroundService() *BingBackgroundService {
	return newBingBackgroundService(&http.Client{Timeout: 8 * time.Second}, defaultBingArchiveURL)
}

func newBingBackgroundService(client *http.Client, endpoint string) *BingBackgroundService {
	return &BingBackgroundService{client: client, endpoint: endpoint}
}

// Resolve 返回经过安全校验的 Bing 图片地址。
func (s *BingBackgroundService) Resolve(ctx context.Context) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.imageURL != "" && time.Now().Before(s.expires) {
		return s.imageURL, nil
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.endpoint, nil)
	if err != nil {
		return "", errcode.New(errcode.ERR1011.Code, errcode.ERR1011.Message)
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return "", errcode.New(errcode.ERR1011.Code, errcode.ERR1011.Message)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", errcode.New(errcode.ERR1011.Code, errcode.ERR1011.Message)
	}
	var payload struct {
		Images []struct {
			URL string `json:"url"`
		} `json:"images"`
	}
	if err = json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&payload); err != nil || len(payload.Images) == 0 {
		return "", errcode.New(errcode.ERR1011.Code, errcode.ERR1011.Message)
	}
	imageURL, err := url.Parse(payload.Images[0].URL)
	if err != nil {
		return "", errcode.New(errcode.ERR1011.Code, errcode.ERR1011.Message)
	}
	if !imageURL.IsAbs() {
		imageURL = (&url.URL{Scheme: "https", Host: "www.bing.com"}).ResolveReference(imageURL)
	}
	host := strings.ToLower(imageURL.Hostname())
	if imageURL.Scheme != "https" || (host != "bing.com" && !strings.HasSuffix(host, ".bing.com")) {
		return "", errcode.New(errcode.ERR1011.Code, errcode.ERR1011.Message)
	}
	s.imageURL = imageURL.String()
	s.expires = time.Now().Add(6 * time.Hour)
	return s.imageURL, nil
}
