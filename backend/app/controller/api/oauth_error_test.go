package api

import (
	"errors"
	"rustdesk-api-server-pro/internal/errcode"
	"testing"
)

// TestClientOAuthCallbackErrorCode 验证客户端 OAuth 回调只返回标准错误码。
func TestClientOAuthCallbackErrorCode(t *testing.T) {
	tests := map[string]string{
		"NoBindableOauthAccount":                errcode.ERR2208.Code,
		"context deadline exceeded":             errcode.ERR2209.Code,
		"StateInvalidOrExpired":                 errcode.ERR2210.Code,
		"ProviderNotAvailableForClientLogin":    errcode.ERR2211.Code,
		"provider returned an invalid response": errcode.ERR2212.Code,
	}
	for message, expected := range tests {
		if actual := clientOAuthCallbackErrorCode(errors.New(message)); actual != expected {
			t.Fatalf("clientOAuthCallbackErrorCode(%q) = %q, want %q", message, actual, expected)
		}
	}
}
