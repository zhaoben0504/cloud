package handler

import (
	"cloud/server"
	"cloud/server/dto"
	"cloud/server/service"
	"cloud/server/vo"
	"cloud/tool"
	"github.com/labstack/echo/v4"
	"net/http"
)

// UploadFileHandler UploadFile handler
type UploadFileHandler struct {
	service service.UploadFileService
}

func (h UploadFileHandler) UploadFile(c echo.Context) error {
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
