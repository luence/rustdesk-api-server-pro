package api

import (
	"rustdesk-api-server-pro/app/form/api"
	"rustdesk-api-server-pro/internal/core"
	v2service "rustdesk-api-server-pro/internal/service"
	"rustdesk-api-server-pro/internal/transport/httpdto"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

var bingBackgroundService = v2service.NewBingBackgroundService()

type SystemController struct {
	basicController
}

// GetBackgroundBing 将 Bing 每日背景重定向到经过校验的官方图片地址。
func (c *SystemController) GetBackgroundBing() mvc.Result {
	imageURL, err := bingBackgroundService.Resolve(c.Ctx.Request().Context())
	if err != nil {
		c.Ctx.Redirect("/login-background.jpg", iris.StatusTemporaryRedirect)
		return mvc.Response{}
	}
	c.Ctx.Header("Cache-Control", "public, max-age=21600")
	c.Ctx.Redirect(imageURL, iris.StatusTemporaryRedirect)
	return mvc.Response{}
}

func (c *SystemController) PostHeartbeat() mvc.Result {
	// {"conns":[762],"id":"182921366","modified_at":1725698100,"uuid":"xxx","ver":1002070}
	var form api.HeartbeatForm
	err := c.readJSONBody(&form)
	if err != nil {
		return c.fail(err)
	}

	result, err := c.systemService().HandleHeartbeat(core.HeartbeatCommand{
		RustdeskID: form.RustdeskId,
		UUID:       form.Uuid,
		ConnCount:  len(form.Conns),
		Conns:      form.Conns,
		ModifiedAt: form.ModifiedAt,
	})
	if err != nil {
		return c.fail(err)
	}

	return mvc.Response{
		Object: httpdto.NewHeartbeatResponse(result),
	}
}

func (c *SystemController) PostSysinfo() mvc.Result {
	var form api.DeviceForm
	err := c.readJSONBody(&form)
	if err != nil {
		return c.fail(err)
	}

	result, err := c.systemService().UpdateSysinfo(core.SysinfoUpdateCommand{
		RustdeskID: form.RustdeskId,
		CPU:        form.Cpu,
		Hostname:   form.Hostname,
		Memory:     form.Memory,
		OS:         form.Os,
		Username:   form.Username,
		UUID:       form.Uuid,
		Version:    form.Version,
	})
	if err != nil {
		return c.fail(err)
	}
	if !result.Updated {
		return mvc.Response{
			Text: "ID_NOT_FOUND",
		}
	}

	return mvc.Response{
		Text: "SYSINFO_UPDATED",
	}
}
