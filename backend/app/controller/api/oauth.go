package api

import (
	"net/url"
	apiform "rustdesk-api-server-pro/app/form/api"
	"rustdesk-api-server-pro/app/model"
	"rustdesk-api-server-pro/config"
	"rustdesk-api-server-pro/internal/core"
	"rustdesk-api-server-pro/internal/errcode"
	v2service "rustdesk-api-server-pro/internal/service"
	"rustdesk-api-server-pro/internal/transport/httpdto"
	"rustdesk-api-server-pro/util"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

// OAuthController exposes /api/oauth/* endpoints so RustDesk desktop and mobile
// clients can complete third-party login through a server-side callback plus
// client polling flow.
type OAuthController struct {
	basicController
}

type oauthStartForm struct {
	Provider   string             `json:"provider"`
	RustdeskId string             `json:"id"`
	Uuid       string             `json:"uuid"`
	DeviceInfo apiform.DeviceInfo `json:"deviceInfo"`
}

type oauthPollForm struct {
	PollToken string `json:"poll_token"`
}

type oauthExchangeForm struct {
	Ticket string `json:"ticket"`
}

func (c *OAuthController) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("GET", "/oauth/providers", "HandleProviders")
	b.Handle("POST", "/oauth/start", "HandleStart")
	b.Handle("GET", "/oauth/{provider:string}/callback", "HandleCallback")
	b.Handle("POST", "/oauth/poll", "HandlePoll")
	b.Handle("POST", "/oauth/exchange", "HandleExchange")
}

func (c *OAuthController) HandleProviders() mvc.Result {
	service := v2service.NewOAuthProviderService(config.GetServerConfig(), c.Db)
	return mvc.Response{Object: service.ListClientProviders()}
}

func (c *OAuthController) HandleStart() mvc.Result {
	var form oauthStartForm
	if err := c.readJSONBody(&form); err != nil {
		return c.fail(err)
	}
	if strings.TrimSpace(form.Provider) == "" {
		return c.fail(errcode.New(errcode.ERR2201.Code, errcode.ERR2201.Message))
	}
	if strings.TrimSpace(form.RustdeskId) == "" || strings.TrimSpace(form.Uuid) == "" {
		return c.fail(errcode.New(errcode.ERR2202.Code, errcode.ERR2202.Message))
	}

	if form.Provider == "webauth" {
		service := v2service.NewOAuthProviderService(config.GetServerConfig(), c.Db)
		loginURL, pollToken, err := service.StartWebauthLogin(c.currentBaseURL(), form.RustdeskId, form.Uuid, form.DeviceInfo.OS, form.DeviceInfo.Type, form.DeviceInfo.Name)
		if err != nil {
			c.recordClientOAuthAudit("", false, "webauth_start: "+err.Error())
			return c.fail(err)
		}
		c.recordClientOAuthAudit("", true, "webauth_start")
		return mvc.Response{Object: iris.Map{
			"enabled":    true,
			"url":        loginURL,
			"poll_token": pollToken,
		}}
	}

	service := v2service.NewOAuthProviderService(config.GetServerConfig(), c.Db)
	authURL, pollToken, enabled, err := service.BuildClientAuthURL(form.Provider, c.currentBaseURL(), form.RustdeskId, form.Uuid, form.DeviceInfo.OS, form.DeviceInfo.Type, form.DeviceInfo.Name)
	if err != nil {
		c.recordClientOAuthAudit("", false, "oauth_start: "+err.Error())
		return c.fail(err)
	}
	if !enabled {
		c.recordClientOAuthAudit("", false, "oauth_start: provider disabled")
		return mvc.Response{Object: iris.Map{"enabled": false}}
	}
	c.recordClientOAuthAudit("", true, "oauth_start: "+form.Provider)
	return mvc.Response{Object: iris.Map{
		"enabled":    true,
		"url":        authURL,
		"poll_token": pollToken,
	}}
}

func (c *OAuthController) HandleCallback() mvc.Result {
	service := v2service.NewOAuthProviderService(config.GetServerConfig(), c.Db)
	provider := c.Ctx.Params().Get("provider")
	code := c.Ctx.URLParamDefault("code", "")
	state := c.Ctx.URLParamDefault("state", "")

	pollToken, bindingTicket, _, err := service.ConsumeUnifiedCallback(provider, code, state)
	if err != nil {
		if bindingTicket != "" && strings.HasPrefix(err.Error(), errcode.ERR2023.Code) {
			query := url.Values{"oauth_bind_ticket": {bindingTicket}, "oauth_provider": {provider}}
			if clientOAuthProviderAllowsAutoCreate(provider) {
				query.Set("oauth_allow_create", "1")
			}
			c.Ctx.Redirect("/#/login?"+query.Encode(), iris.StatusFound)
			return mvc.Response{}
		}
		c.recordClientOAuthAudit("", false, "client_oauth_callback: "+provider+": "+err.Error())
		return c.oauthCallbackPage(false, clientOAuthCallbackErrorCode(err), "")
	}
	c.recordClientOAuthAudit("", true, "client_oauth_callback: "+provider+": poll_token")
	return c.oauthCallbackPage(true, "", pollToken)
}

func clientOAuthProviderAllowsAutoCreate(providerName string) bool {
	for _, provider := range config.GetServerConfig().OAuthProviders() {
		if provider.Name == providerName {
			return provider.AutoCreateUser
		}
	}
	return false
}

func (c *OAuthController) HandlePoll() mvc.Result {
	var form oauthPollForm
	if err := c.readJSONBody(&form); err != nil {
		return c.fail(err)
	}
	service := v2service.NewOAuthProviderService(config.GetServerConfig(), c.Db)
	ticket, ready, err := service.PollClientTicket(form.PollToken)
	if err != nil {
		return c.fail(err)
	}
	if !ready {
		return mvc.Response{Object: iris.Map{"ready": false}}
	}
	return mvc.Response{Object: iris.Map{
		"ready":  true,
		"ticket": ticket,
	}}
}

func (c *OAuthController) HandleExchange() mvc.Result {
	var form oauthExchangeForm
	if err := c.readJSONBody(&form); err != nil {
		return c.fail(err)
	}
	service := v2service.NewOAuthProviderService(config.GetServerConfig(), c.Db)
	token, user, err := service.ExchangeClientTicket(form.Ticket)
	if err != nil {
		c.recordClientOAuthAudit("", false, "client_oauth_exchange: "+err.Error())
		return c.fail(err)
	}
	c.recordClientOAuthAudit(user.Username, true, "client_oauth_exchange: access_token")
	return mvc.Response{Object: iris.Map{
		"access_token": token,
		"type":         "access_token",
		"user": iris.Map{
			"name":         user.Name,
			"display_name": user.Name,
			"email":        user.Email,
			"note":         user.Note,
			"status":       user.Status,
			"is_admin":     false,
		},
	}}
}

type webauthConfirmForm struct {
	PollToken string `json:"poll_token"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

func (c *OAuthController) HandleWebauthConfirm() mvc.Result {
	var form webauthConfirmForm
	if err := c.readJSONBody(&form); err != nil {
		return c.fail(err)
	}
	if strings.TrimSpace(form.PollToken) == "" {
		return c.fail(errcode.New(errcode.ERR2205.Code, errcode.ERR2205.Message))
	}
	if form.Username == "" || form.Password == "" {
		return c.fail(errcode.New(errcode.ERR2203.Code, errcode.ERR2203.Message))
	}

	var user model.User
	has, err := c.Db.Where("username = ? and status > 0", form.Username).Get(&user)
	if err != nil || !has {
		c.recordClientOAuthAudit(form.Username, false, "webauth_confirm: user not found")
		return c.fail(errcode.New(errcode.ERR1002.Code, errcode.ERR1002.Message))
	}
	if !util.PasswordVerify(form.Password, user.Password) {
		c.recordClientOAuthAudit(form.Username, false, "webauth_confirm: password error")
		return c.fail(errcode.New(errcode.ERR1003.Code, errcode.ERR1003.Message))
	}

	service := v2service.NewOAuthProviderService(config.GetServerConfig(), c.Db)
	if err := service.ConfirmWebauthLogin(form.PollToken, user.Id); err != nil {
		c.recordClientOAuthAudit("", false, "webauth_confirm: "+err.Error())
		return c.fail(err)
	}

	c.recordClientOAuthAudit(form.Username, true, "webauth_confirm: success")
	return mvc.Response{Object: iris.Map{"ok": true}}
}

func (c *OAuthController) HandleWebauthLoginPage() mvc.Result {
	pollToken := c.Ctx.URLParam("poll_token")

	html := `<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>登录</title>
<style>
* { margin: 0; padding: 0; box-sizing: border-box; }
body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; display: flex; justify-content: center; align-items: center; min-height: 100vh; padding: 16px; background-color: #07172f; background-image: linear-gradient(rgba(5,14,31,.34),rgba(5,14,31,.48)), var(--web-background, url("/login-background.jpg")); background-position: center; background-size: cover; background-repeat: no-repeat; }
.card { background: rgba(255,255,255,.9); border-radius: 12px; padding: 40px; max-width: 380px; width: 90%; box-shadow: 0 20px 55px rgba(0,0,0,.28); backdrop-filter: blur(14px); }
h2 { text-align: center; margin-bottom: 24px; color: #333; }
input { width: 100%; padding: 12px 16px; border: 1px solid #d9d9d9; border-radius: 8px; font-size: 15px; margin-bottom: 16px; transition: border-color .2s; }
input:focus { outline: none; border-color: #1890ff; }
button { width: 100%; padding: 12px; border: none; border-radius: 8px; font-size: 15px; font-weight: 600; cursor: pointer; transition: all .2s; }
.btn-primary { background: #1890ff; color: #fff; }
.btn-primary:hover { background: #40a9ff; }
#status { text-align: center; margin-top: 16px; font-size: 14px; min-height: 20px; }
.error { color: #ff4d4f; }
.success { color: #52c41a; }
.launch-btn { display: block; width: 100%; padding: 14px; background: #52c41a; color: #fff; border: none; border-radius: 8px; font-size: 16px; font-weight: 600; text-decoration: none; text-align: center; margin-top: 16px; }
.launch-btn:hover { background: #73d13d; }
</style>
</head>
<body>
<script>(function(){try{var m=localStorage.getItem("rustdesk-background-mode"),u="/login-background.jpg";if(m==="bing")u="/api/background/bing";else if(m==="upload"){var x=localStorage.getItem("rustdesk-background-upload");if(x&&x.indexOf("data:image/")===0)u=x}document.documentElement.style.setProperty("--web-background",'url("'+u.replace(/["\\]/g,"")+'")')}catch(e){}})();</script>
<div class="card">
  <h2>RustDesk 登录</h2>
  <div id="login-form">
    <input type="text" id="username" placeholder="用户名" autocomplete="username">
    <input type="password" id="password" placeholder="密码" autocomplete="current-password">
    <button class="btn-primary" onclick="doLogin()">登录</button>
  </div>
  <div id="status"></div>
  <div id="success-section" style="display:none">
    <p class="success" style="text-align:center;margin-bottom:8px">登录成功！</p>
    <a href="#" class="launch-btn" id="launch-btn">返回 RustDesk 客户端</a>
    <p style="color:#999;font-size:13px;text-align:center;margin-top:12px">请点击上方按钮返回客户端<br>如按钮无效，请手动切回客户端</p>
  </div>
</div>
<script>
const pollToken = "` + pollToken + `";

function setStatus(msg, cls) {
  const el = document.getElementById("status");
  el.textContent = msg;
  el.className = cls || "";
}

async function doLogin() {
  const username = document.getElementById("username").value.trim();
  const password = document.getElementById("password").value;
  if (!username || !password) { setStatus("请输入用户名和密码", "error"); return; }

  setStatus("正在登录...");
  try {
    const resp = await fetch("/api/oauth/webauth/confirm", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ poll_token: pollToken, username, password })
    });
    const data = await resp.json();
    if (data.ok) {
      document.getElementById("login-form").style.display = "none";
      document.getElementById("success-section").style.display = "block";
      document.getElementById("launch-btn").href = "rustdesk://config/";
      setStatus("登录成功！", "success");
    } else {
      setStatus(data.error || "登录失败", "error");
    }
  } catch (e) {
    setStatus("登录出错: " + (e.message || e), "error");
  }
}
</script>
</body>
</html>`

	c.Ctx.ContentType("text/html; charset=utf-8")
	_, _ = c.Ctx.WriteString(html)
	return mvc.Response{}
}

func (c *OAuthController) currentBaseURL() string {
	scheme := "http"
	if c.Ctx.Request().TLS != nil {
		scheme = "https"
	}
	if forwardedProto := strings.TrimSpace(c.Ctx.GetHeader("X-Forwarded-Proto")); forwardedProto == "https" {
		scheme = "https"
	}
	host := strings.TrimSpace(c.Ctx.Host())
	if host == "" {
		host = strings.TrimSpace(c.Ctx.Request().Host)
	}
	return scheme + "://" + host
}

func (c *OAuthController) oauthCallbackPage(success bool, errorCode, pollToken string) mvc.Result {
	c.Ctx.Redirect(httpdto.OAuthCallbackURL(success, errorCode), iris.StatusFound)
	return mvc.Response{}
}

func (c *OAuthController) recordClientOAuthAudit(username string, success bool, reason string) {
	_ = c.auditService().CreateSecurityAudit(core.SecurityAuditCreateCommand{
		Username:  username,
		Event:     "client_oauth",
		IP:        c.Ctx.RemoteAddr(),
		UserAgent: c.Ctx.GetHeader("User-Agent"),
		Success:   success,
		Reason:    reason,
	})
}

func clientOAuthCallbackErrorCode(err error) string {
	if err == nil {
		return errcode.ERR2212.Code
	}
	msg := err.Error()
	if strings.HasPrefix(msg, "ERR-") {
		code := strings.SplitN(msg, ":", 2)[0]
		switch code {
		case errcode.ERR2023.Code, errcode.ERR2022.Code:
			return errcode.ERR2208.Code
		case errcode.ERR2004.Code, errcode.ERR2029.Code:
			return errcode.ERR2210.Code
		case errcode.ERR2203.Code, errcode.ERR2211.Code:
			return errcode.ERR2211.Code
		case errcode.ERR2030.Code, errcode.ERR2031.Code, errcode.ERR2034.Code:
			return errcode.ERR2209.Code
		default:
			return errcode.ERR2212.Code
		}
	}
	switch {
	case strings.Contains(msg, "NoBindableOauthAccount"), strings.Contains(msg, "BoundAdminUserNotAvailable"):
		return errcode.ERR2208.Code
	case strings.Contains(msg, "timeout"), strings.Contains(msg, "deadline exceeded"), strings.Contains(msg, "connection refused"):
		return errcode.ERR2209.Code
	case strings.Contains(msg, "StateInvalidOrExpired"), strings.Contains(msg, "StateExpired"):
		return errcode.ERR2210.Code
	case strings.Contains(msg, "ProviderNotAvailableForClientLogin"):
		return errcode.ERR2211.Code
	default:
		return errcode.ERR2212.Code
	}
}
