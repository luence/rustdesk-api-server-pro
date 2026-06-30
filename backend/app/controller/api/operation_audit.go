package api

import (
	"encoding/json"
	"rustdesk-api-server-pro/internal/core"
)

func (c *basicController) recordAPIOperationAudit(action string, resourceType string, resourceID string, beforeData interface{}, afterData interface{}, result string, errorMessage string) {
	if c == nil || c.Db == nil || c.Ctx == nil {
		return
	}
	user := c.GetUser()
	cmd := core.OperationAuditCreateCommand{
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		BeforeData:   operationAuditJSON(beforeData),
		AfterData:    operationAuditJSON(afterData),
		IP:           c.Ctx.RemoteAddr(),
		UserAgent:    c.Ctx.GetHeader("User-Agent"),
		Result:       result,
		ErrorMessage: errorMessage,
	}
	if user != nil {
		cmd.ActorUserID = user.Id
		cmd.ActorUsername = user.Username
	}
	_ = c.auditService().CreateOperationAudit(cmd)
}

func operationAuditJSON(value interface{}) string {
	if value == nil {
		return ""
	}
	data, err := json.Marshal(value)
	if err != nil {
		return ""
	}
	return string(data)
}
