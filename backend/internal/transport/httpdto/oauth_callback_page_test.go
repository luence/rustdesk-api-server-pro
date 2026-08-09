package httpdto

import (
	"strings"
	"testing"
)

// TestOAuthCallbackPage 验证统一页面不泄露 poll token，并转义错误信息。
func TestOAuthCallbackPage(t *testing.T) {
	page := OAuthCallbackPage(true, "", "sensitive-poll-token")
	if strings.Contains(page, "sensitive-poll-token") {
		t.Fatal("回调页面不得输出 poll token")
	}
	if !strings.Contains(page, "RustDesk API 服务端") || !strings.Contains(page, "rustdesk://config/") {
		t.Fatal("成功页面缺少统一品牌或客户端返回入口")
	}

	failure := OAuthCallbackPage(false, `<ERR-2212>`, "")
	if strings.Contains(failure, `<ERR-2212>`) || !strings.Contains(failure, `&lt;ERR-2212&gt;`) {
		t.Fatal("错误码必须经过 HTML 转义")
	}
}
