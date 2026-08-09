package httpdto

import (
	"net/url"
	"strings"
	"testing"
)

// TestOAuthCallbackURL 验证结果页只传递公开状态和转义后的错误码。
func TestOAuthCallbackURL(t *testing.T) {
	success := OAuthCallbackURL(true, "")
	if success != "/#/client-oauth-result?status=success" {
		t.Fatalf("成功结果页地址错误: %s", success)
	}

	failure := OAuthCallbackURL(false, `<ERR-2212>`)
	parts := strings.SplitN(failure, "?", 2)
	if len(parts) != 2 || parts[0] != "/#/client-oauth-result" {
		t.Fatalf("失败结果页地址错误: %s", failure)
	}
	query, err := url.ParseQuery(parts[1])
	if err != nil || query.Get("status") != "error" || query.Get("error_code") != `<ERR-2212>` {
		t.Fatalf("失败结果页参数错误: %s", failure)
	}
}
