package api

import (
	"encoding/json"
	"testing"
)

// TestRustdeskOIDCQueryResponse 验证轮询接口直接输出客户端需要解析的认证体。
func TestRustdeskOIDCQueryResponse(t *testing.T) {
	inner := `{"access_token":"token","type":"access_token","user":{"name":"tester","status":1,"info":{}}}`
	response := rustdeskOIDCQueryResponse(inner)
	if response.ContentType != "application/json; charset=utf-8" {
		t.Fatalf("unexpected content type: %s", response.ContentType)
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(response.Text), &result); err != nil {
		t.Fatalf("decode auth response: %v", err)
	}
	if result["type"] != "access_token" {
		t.Fatalf("response must be the direct auth body: %s", response.Text)
	}
	if _, exists := result["body"]; exists {
		t.Fatalf("response must not add a second body wrapper: %s", response.Text)
	}
}
