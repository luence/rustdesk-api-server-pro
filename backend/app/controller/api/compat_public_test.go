package api

import (
	"encoding/json"
	"testing"
)

// TestRustdeskOIDCQueryResponse 验证轮询响应始终保持官方客户端要求的字符串 body。
func TestRustdeskOIDCQueryResponse(t *testing.T) {
	inner := `{"access_token":"token","type":"access_token","user":{"name":"tester","status":1,"info":{}}}`
	response := rustdeskOIDCQueryResponse(inner)
	if response.ContentType != "application/json; charset=utf-8" {
		t.Fatalf("unexpected content type: %s", response.ContentType)
	}

	var outer struct {
		Body string `json:"body"`
	}
	if err := json.Unmarshal([]byte(response.Text), &outer); err != nil {
		t.Fatalf("decode outer response: %v", err)
	}
	if outer.Body != inner {
		t.Fatalf("body must remain an encoded JSON string: %s", response.Text)
	}
}
