package httpdto

import (
	"fmt"
	"html"
	"strconv"
	"strings"
)

// OAuthCallbackPage 返回与 WebAuth 登录页一致的第三方登录结果页面。
func OAuthCallbackPage(success bool, errorCode, pollToken string) string {
	title := "第三方登录失败"
	statusClass := "error"
	message := "登录失败，请返回客户端后重试。"
	if success {
		title = "第三方登录成功"
		statusClass = "success"
		message = "认证已完成，客户端正在自动获取登录结果。"
	}

	detail := ""
	if !success && strings.TrimSpace(errorCode) != "" {
		detail = `<p class="code">错误码：` + html.EscapeString(errorCode) + `</p>`
	}

	actions := ""
	if success && strings.TrimSpace(pollToken) != "" {
		actions += `<a class="button primary" href="rustdesk://config/">返回 RustDesk 客户端</a>`
	}
	actions += `<button class="button secondary" type="button" onclick="window.close()">关闭页面</button>`

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>%s</title>
<style>
*{box-sizing:border-box}body{min-height:100vh;margin:0;display:flex;align-items:center;justify-content:center;padding:24px;font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,"PingFang SC","Microsoft YaHei",sans-serif;color:#f4f4f5;background:linear-gradient(135deg,#5b5ce2 0%%,#7c83eb 52%%,#a7a9f5 100%%)}
.card{width:min(368px,calc(100vw - 48px));padding:32px 24px;border:1px solid rgba(255,255,255,.08);border-radius:12px;background:#18181c;box-shadow:0 18px 50px rgba(20,20,43,.28);text-align:center}
.brand{display:flex;align-items:center;justify-content:center;gap:10px;margin-bottom:26px;color:#6366f1;font-size:22px;font-weight:600}.logo{width:34px;height:34px;border:8px solid #06a6dd;border-right-color:transparent;border-radius:50%%;transform:rotate(-20deg)}
.status{width:48px;height:48px;margin:0 auto 16px;display:flex;align-items:center;justify-content:center;border-radius:50%%;font-size:26px;font-weight:700}.status.success{color:#22c55e;background:rgba(34,197,94,.14)}.status.error{color:#ef4444;background:rgba(239,68,68,.14)}
h1{margin:0 0 12px;font-size:20px}.message,.code,.tip{margin:0;color:#a1a1aa;font-size:14px;line-height:1.65}.code{margin-top:8px;color:#fca5a5}.actions{display:grid;gap:10px;margin-top:24px}.button{width:100%%;min-height:40px;display:flex;align-items:center;justify-content:center;border:0;border-radius:8px;font:inherit;font-weight:600;text-decoration:none;cursor:pointer}.primary{color:#fff;background:#6366f1}.primary:hover{background:#5558e8}.secondary{color:#d4d4d8;background:#303036}.secondary:hover{background:#3f3f46}.tip{margin-top:18px;font-size:12px}
</style>
</head>
<body>
<main class="card">
  <div class="brand"><span class="logo"></span><span>RustDesk API 服务端</span></div>
  <div class="status %s">%s</div>
  <h1>%s</h1>
  <p class="message">%s</p>
  %s
  <div class="actions">%s</div>
  <p class="tip">若浏览器阻止自动关闭，可使用上方按钮或直接关闭此标签页。</p>
</main>
<script>if(%s){window.setTimeout(function(){window.close()},1200)}</script>
</body>
</html>`, html.EscapeString(title), statusClass, map[bool]string{true: "✓", false: "!"}[success], html.EscapeString(title), html.EscapeString(message), detail, actions, strconv.FormatBool(success))
}
