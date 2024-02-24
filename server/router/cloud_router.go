package router

import (
	"cloud/middleware"
	"net/http"
)

func initCloudRouter() {
	list := []middleware.PermissionItem{
		{
			Method:      http.MethodPost,
			Handler:     userHandler.Login,
			URL:         "/user/login",
			Permissions: []string{"A"},
		},
		{
			Method:      http.MethodPost,
			Handler:     userHandler.Register,
			URL:         "/user/register",
			Permissions: []string{"A"},
		},
		{
			Method:      http.MethodPost,
			Handler:     userHandler.EmailCode,
			URL:         "/user/code",
			Permissions: []string{"A"},
		},
		//{
		//	Method:      http.MethodGet,
		//	Handler:     ,
		//	URL:         "",
		//	Permissions: []string{"A"},
		//},
	}

	middleware.GenerateHandler(Echo, list)
}
