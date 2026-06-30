package middleware

import "github.com/kataras/iris/v12"

func RequestLogger(debugMode bool) iris.Handler {
	return func(context iris.Context) {
		if debugMode && context.Request().URL.Path != "/api/heartbeat" {
			context.Application().Logger().Infof("▶ %s:%s", context.Method(), context.Request().URL.Path)
		}

		context.Next()
	}
}
