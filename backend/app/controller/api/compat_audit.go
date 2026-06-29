package api

import "rustdesk-api-server-pro/internal/core"

func (c *basicController) recordCompatAPIAudit(isStub bool, statusCode int, result string, errorMessage string, body []byte) {
	if c == nil || c.Db == nil || c.Ctx == nil {
		return
	}
	if statusCode == 0 {
		statusCode = 200
	}
	if result == "" {
		result = "ok"
	}

	path := c.Ctx.Path()
	if path == "" && c.Ctx.Request() != nil && c.Ctx.Request().URL != nil {
		path = c.Ctx.Request().URL.Path
	}

	_ = c.auditService().CreateCompatAPIAudit(core.CompatAPIAuditCreateCommand{
		Method:        c.Ctx.Method(),
		Path:          path,
		ClientVersion: firstNonEmptyString(c.Ctx.GetHeader("X-RustDesk-Version"), c.Ctx.GetHeader("RustDesk-Version"), c.Ctx.GetHeader("X-Client-Version"), c.Ctx.GetHeader("Client-Version")),
		RustdeskID:    compatRustdeskID(c, body),
		IsStub:        isStub,
		StatusCode:    statusCode,
		IP:            c.Ctx.RemoteAddr(),
		UserAgent:     c.Ctx.GetHeader("User-Agent"),
		Result:        result,
		ErrorMessage:  errorMessage,
	})
}

func compatRustdeskID(c *basicController, body []byte) string {
	if c != nil && c.Ctx != nil {
		if id := firstNonEmptyString(c.Ctx.URLParam("id"), c.Ctx.URLParam("rustdesk_id"), c.Ctx.URLParam("rustdeskId")); id != "" {
			return id
		}
	}
	if len(body) == 0 {
		return ""
	}
	return firstJSONValue(body, "id", "rustdesk_id", "rustdeskId")
}

func firstNonEmptyString(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
