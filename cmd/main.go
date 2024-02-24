package main

import (
	"cloud/server"
	"cloud/server/router"
	"cloud/tool"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"path/filepath"
)

var (
	DepMode   string
	BuildTime string
	GoVersion string
	Version   string
)

func main() {
	DepMode = os.Getenv("DEP_MODE")
	tool.Logger.SetFormatter(&logrus.JSONFormatter{})
	tool.Logger.AddHook(tool.NewContextHook())

	tool.Logger.Info("mode:", DepMode)
	if len(DepMode) == 0 {
		tool.Logger.Fatal("env DEP_MODE is absent")
	}

	err := server.NewServer("./config/", DepMode)
	if nil != err {
		tool.Logger.Fatal(err.Error())
	}
	server.LoadMessageFile([]string{"./i18n/cloud.en.yaml", "./i18n/cloud.zh.yaml"})
	router.InitRouter()
	router.Echo.GET("/cloud/about", about)
	router.Echo.Logger.Fatal(router.Echo.Start(fmt.Sprintf(":%d", server.GetPort())))

}

func about(c echo.Context) error {
	type profile struct {
		DepMode   string
		BuildTime string
		GoVersion string
		Version   string
		Desc      string
	}
	p := profile{
		DepMode:   DepMode,
		BuildTime: BuildTime,
		GoVersion: GoVersion,
		Version:   Version,
		Desc:      fmt.Sprintf("%s application.", filepath.Base(os.Args[0])),
	}
	return c.JSON(http.StatusOK, p)
}
