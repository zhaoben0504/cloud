package handler

import (
	"bytes"
	"cloud/server"
	"cloud/server/dto"
	"cloud/server/service"
	"cloud/server/vo"
	"cloud/tool"
	"github.com/labstack/echo/v4"
	"net/http"
)

// UploadFileHandler UploadFile handler
type FileHandler struct {
	service service.FileService
}

func (h FileHandler) UploadFile(c echo.Context) error {
	var req dto.UploadFileDTO
	lang := tool.GetHeaderLanguage(c.Request().Header)
	if err := c.Bind(&req); err != nil {
		tool.Logger.Error(err.Error())
		return c.JSON(http.StatusOK, server.NewError(lang, server.ParamErrCode))
	}
	if err := c.Validate(req); err != nil {
		tool.Logger.Error(err.Error())
		return c.JSON(http.StatusOK, server.NewError(lang, server.ParamErrCode))
	}
	addr, code := h.service.UploadFile(req)
	if code != server.OkCode {
		tool.Logger.Error(server.GetMsgByCode(lang, code))
		return c.JSON(http.StatusOK, server.NewError(lang, code))
	}

	return c.JSON(http.StatusOK, server.NewOK(lang, vo.UploadFileVO{FileAddr: addr}))
}

func (h FileHandler) DownloadFile(c echo.Context) error {
	var req dto.DownloadFileDTO
	lang := tool.GetHeaderLanguage(c.Request().Header)
	if err := c.Bind(&req); err != nil {
		tool.Logger.Error(err.Error())
		return c.JSON(http.StatusOK, server.NewError(lang, server.ParamErrCode))
	}
	if err := c.Validate(req); err != nil {
		tool.Logger.Error(err.Error())
		return c.JSON(http.StatusOK, server.NewError(lang, server.ParamErrCode))
	}
	file, code := h.service.DownloadFile(req)
	if code != server.OkCode {
		tool.Logger.Error(server.GetMsgByCode(lang, code))
		return c.JSON(http.StatusOK, server.NewError(lang, code))
	}

	c.Response().Header().Set("md5", tool.MD5(string(file)))
	return c.Stream(http.StatusOK, echo.MIMEOctetStream, bytes.NewReader(file))
}
