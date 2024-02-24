package router

import (
	"cloud/middleware"
	"cloud/server"
	"cloud/server/handler"
	"cloud/tool"
	"context"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"os"
)

var (
	Echo         = echo.New()
	cloudHandler = handler.UserHandler{}
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return err
	}
	return nil
}

func InitRouter() {
	Echo.Validator = &CustomValidator{validator: validator.New()}
	Echo.Use(middleware.Record())
	Echo.Use(middleware.RecoverWithReturnMsg(server.NewError(tool.GetHeaderLanguage(nil), server.InternalErrCode)))
	cors := os.Getenv("CORS")
	if cors == "true" {
		Echo.Use(echoMiddleware.CORS())
	}

	Echo.Use(middleware.PermissionWithConfig(middleware.PermissionMiddlewareConfig{
		Key: "token",
		IgnoreURLs: []string{
			"/cloud/about",
		},
		GetPermissionList: func(k string) []string {
			client := server.GetRedisClient()
			val, err := client.Get(context.Background(), k).Result()
			if err != nil {
				tool.Logger.Errorf("get token %s error: %v", k, err)
				return nil
			}
			if len(val) == 0 {
				return nil
			}
			var user server.RedisUserInfo
			err = json.Unmarshal([]byte(val), &user)
			if err != nil {
				tool.Logger.Error(err)
				return nil
			}
			return user.Permissions
		},
		GetContext: func(k string) map[string]interface{} {
			client := server.GetRedisClient()
			val, err := client.Get(context.Background(), k).Result()
			if err != nil {
				tool.Logger.Error(err)
				return nil
			}
			var info server.RedisUserInfo
			err = json.Unmarshal([]byte(val), &info)
			if err != nil {
				tool.Logger.Error(err)
				return nil
			}
			return map[string]interface{}{
				// handler 需要用的值
			}
		},
		InternalErrFunc: func(lang string) interface{} {
			return server.NewError(lang, server.InternalErrCode)
		},
		TokenNotExistErrFunc: func(lang string) interface{} {
			return server.NewError(lang, server.TokenInvalidErrCode)
		},
		TokenInvalidErrFunc: func(lang string) interface{} {
			return server.NewError(lang, server.TokenInvalidErrCode)
		},
		PermissionErrFunc: func(lang string) interface{} {
			return server.NewError(lang, server.PermissionErrCode)
		},
	}))

	initCloudRouter()
}
