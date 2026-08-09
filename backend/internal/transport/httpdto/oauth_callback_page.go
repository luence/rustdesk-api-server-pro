package httpdto

import "net/url"

// OAuthCallbackURL 返回复用 Web 前端主题的客户端第三方登录结果页地址。
func OAuthCallbackURL(success bool, errorCode string) string {
	values := url.Values{}
	if success {
		values.Set("status", "success")
	} else {
		values.Set("status", "error")
		if errorCode != "" {
			values.Set("error_code", errorCode)
		}
	}
	return "/#/client-oauth-result?" + values.Encode()
}
