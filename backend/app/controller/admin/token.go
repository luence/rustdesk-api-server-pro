package admin

import (
	"rustdesk-api-server-pro/app/model"
	"rustdesk-api-server-pro/config"
	"rustdesk-api-server-pro/db"
	"rustdesk-api-server-pro/internal/errcode"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"xorm.io/xorm"
)

type TokenController struct {
	basicController
}

func (c *TokenController) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("GET", "/tokens/list", "HandleList")
	b.Handle("POST", "/tokens/kill", "HandleKill")
	b.Handle("POST", "/tokens/clear", "HandleClear")
}

func (c *TokenController) HandleList() mvc.Result {
	currentPage := c.Ctx.URLParamIntDefault("current", 1)
	pageSize := c.Ctx.URLParamIntDefault("size", 10)
	username := c.Ctx.URLParamDefault("username", "")
	filterUserIDs := make([]int, 0)
	if username != "" {
		matchedUsers := make([]model.User, 0)
		if err := c.Db.Where("username LIKE ?", "%"+username+"%").Cols("id").Find(&matchedUsers); err != nil {
			return c.dbError(err)
		}
		for _, user := range matchedUsers {
			filterUserIDs = append(filterUserIDs, user.Id)
		}
	}

	query := func() *xorm.Session {
		q := c.Db.Table(&model.AuthToken{})
		if username != "" {
			if len(filterUserIDs) == 0 {
				q.Where("1 = 0")
			} else {
				q.In("auth_token.user_id", filterUserIDs)
			}
		}
		q.Desc("auth_token.id")
		return q
	}

	pagination := db.NewPagination(currentPage, pageSize)
	tokenList := make([]model.AuthToken, 0)
	if err := pagination.Paginate(query, &model.AuthToken{}, &tokenList); err != nil {
		return c.dbError(err)
	}
	userIDs := make([]int, 0, len(tokenList))
	for _, token := range tokenList {
		userIDs = append(userIDs, token.UserId)
	}
	users := make([]model.User, 0)
	if len(userIDs) > 0 {
		if err := c.Db.In("id", userIDs).Cols("id", "username").Find(&users); err != nil {
			return c.dbError(err)
		}
	}
	usernames := make(map[int]string, len(users))
	for _, user := range users {
		usernames[user.Id] = user.Username
	}

	list := make([]iris.Map, 0, len(tokenList))
	currentToken := c.GetAuthToken()
	for _, t := range tokenList {
		isCurrent := currentToken != nil && currentToken.Id == t.Id
		list = append(list, iris.Map{
			"id":          t.Id,
			"user_id":     t.UserId,
			"username":    usernames[t.UserId],
			"rustdesk_id": t.RustdeskId,
			"uuid":        t.Uuid,
			"device_os":   t.DeviceOs,
			"device_type": t.DeviceType,
			"device_name": t.DeviceName,
			"token_hash":  t.TokenHash,
			"is_admin":    t.IsAdmin,
			"status":      t.Status,
			"expired":     t.Expired.Format(config.TimeFormat),
			"created_at":  t.CreatedAt.Format(config.TimeFormat),
			"is_current":  isCurrent,
		})
	}

	return c.Success(iris.Map{
		"total":   pagination.TotalCount,
		"records": list,
		"current": currentPage,
		"size":    pageSize,
	}, "ok")
}

func (c *TokenController) HandleKill() mvc.Result {
	user := c.GetUser()
	if user == nil || !user.IsAdmin {
		return c.Error(nil, errcode.New(errcode.ERR1010.Code, errcode.ERR1010.Message).Error())
	}
	var params struct {
		Ids []int `json:"ids"`
	}
	if err := c.Ctx.ReadJSON(&params); err != nil {
		return c.dbError(err)
	}
	if len(params.Ids) == 0 {
		return c.Error(nil, errcode.New(errcode.ERR9002.Code, errcode.ERR9002.Message).Error())
	}

	currentToken := c.GetAuthToken()
	if currentToken != nil {
		for i, id := range params.Ids {
			if id == currentToken.Id {
				params.Ids = append(params.Ids[:i], params.Ids[i+1:]...)
				break
			}
		}
	}
	if len(params.Ids) == 0 {
		return c.Error(nil, errcode.New(errcode.ERR9002.Code, errcode.ERR9002.Message).Error())
	}

	affected, err := c.Db.In("id", params.Ids).Delete(&model.AuthToken{})
	if err != nil {
		return c.dbError(err)
	}
	if affected == 0 {
		return c.Error(nil, errcode.New(errcode.ERR9004.Code, errcode.ERR9004.Message).Error())
	}

	return c.Success(nil, "ok")
}

func (c *TokenController) HandleClear() mvc.Result {
	currentUser := c.GetUser()
	if currentUser == nil || !currentUser.IsAdmin {
		return c.Error(nil, errcode.New(errcode.ERR1010.Code, errcode.ERR1010.Message).Error())
	}

	currentToken := c.GetAuthToken()
	if currentToken == nil {
		return c.Error(nil, errcode.New(errcode.ERR9002.Code, errcode.ERR9002.Message).Error())
	}
	affected, err := c.Db.Where("id != ?", currentToken.Id).Delete(&model.AuthToken{})
	if err != nil {
		return c.dbError(err)
	}

	return c.Success(iris.Map{"cleared": affected, "retained": 1}, "ok")
}
