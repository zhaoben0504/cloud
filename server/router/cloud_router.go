package router

func initCloudRouter() {
	list := []middleware.PermissionItem{
		//{
		//	Method:      http.MethodPost,
		//	Handler:     applicationHandler.Add,
		//	URL:         "/lcdp/app",
		//	Permissions: []string{"A"},
		//},
		//{
		//	Method:      http.MethodGet,
		//	Handler:     applicationHandler.Test,
		//	URL:         "/lcdp/app/resources/:appid/:filename",
		//	Permissions: []string{"A"},
		//},
	}

	middleware.GenerateHandler(Echo, list)
}
