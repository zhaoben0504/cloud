package handler

import (
	"cloud/server"
	"cloud/server/dto"
	"cloud/server/service"
	"cloud/tool"
	"github.com/labstack/echo/v4"
	"net/http"
)

// UserHandler user handler
type UserHandler struct {
	service service.UserService
}

func (h UserHandler) Register(c echo.Context) error {
	var req dto.UserRegisterDTO
	lang := tool.GetHeaderLanguage(c.Request().Header)
	if err := c.Bind(&req); err != nil {
		tool.Logger.Error(err.Error())
		return c.JSON(http.StatusOK, server.NewError(lang, server.ParamErrCode))
	}
	if err := c.Validate(req); err != nil {
		return err
	}
	resp, code := h.service.Register(&req)
	if code != server.OkCode {
		return c.JSON(http.StatusOK, server.NewError(lang, code))
	}
	return c.JSON(http.StatusOK, server.NewOK(lang, resp))
}

func (h UserHandler) Login(c echo.Context) error {
	var req dto.UserLoginDTO
	lang := tool.GetHeaderLanguage(c.Request().Header)
	if err := c.Bind(&req); err != nil {
		tool.Logger.Error(err.Error())
		return c.JSON(http.StatusOK, server.NewError(lang, server.ParamErrCode))
	}
	if err := c.Validate(req); err != nil {
		return err
	}
	resp, code := h.service.Login(&req)
	if code != server.OkCode {
		return c.JSON(http.StatusOK, server.NewError(lang, code))
	}
	return c.JSON(http.StatusOK, server.NewOK(lang, resp))
}

func (h UserHandler) EmailCode(c echo.Context) error {
	var req dto.EmailCodeDTO
	lang := tool.GetHeaderLanguage(c.Request().Header)
	if err := c.Bind(&req); err != nil {
		tool.Logger.Error(err.Error())
		return c.JSON(http.StatusOK, server.NewError(lang, server.ParamErrCode))
	}
	if err := c.Validate(req); err != nil {
		return err
	}
	resp, code := h.service.EmailCode(&req)
	if code != server.OkCode {
		return c.JSON(http.StatusOK, server.NewError(lang, code))
	}
	return c.JSON(http.StatusOK, server.NewOK(lang, resp))
}

func (h UserHandler) UserInfo(c echo.Context) error {
	var req dto.UserInfoDTO
	lang := tool.GetHeaderLanguage(c.Request().Header)
	if err := c.Bind(&req); err != nil {
		tool.Logger.Error(err.Error())
		return c.JSON(http.StatusOK, server.NewError(lang, server.ParamErrCode))
	}
	if err := c.Validate(req); err != nil {
		return err
	}
	result, code := h.service.UserInfo(&req)
	if code != server.OkCode {
		return c.JSON(http.StatusOK, server.NewError(lang, code))
	}
	return c.JSON(http.StatusOK, server.NewOK(lang, result))
}
