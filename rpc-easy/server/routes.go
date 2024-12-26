package server

import "github.com/labstack/echo/v4"

func NewHttpServer() *echo.Echo {
	r := echo.New()
	r.GET("/ping", func(c echo.Context) error {
		return c.String(200, "hello world")
	})
	r.POST("/rpc", handle)
	return r
}
