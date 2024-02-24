package middleware

import (
	"cloud/tool"
	"encoding/json"
	"github.com/labstack/echo/v4"
)

func Record() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			before := tool.UnixMillisecond()
			defer func() {
				after := tool.UnixMillisecond()
				tool.Logger.Infof("url: %s, time consuming: %d ms", c.Request().RequestURI, after-before)
			}()

			return next(RecordContext{
				Context: c,
			})
		}
	}
}

type RecordContext struct {
	echo.Context
}

func (r RecordContext) JSON(code int, i interface{}) error {
	err := r.Context.JSON(code, i)

	if err == nil {
		data, _ := json.Marshal(i)
		tool.Logger.Infof("url: %s, response: %s", r.Context.Request().RequestURI, data)
	}

	return err
}
