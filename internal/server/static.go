package server

import (
	"github.com/labstack/echo"
)

func serveStaticFiles(echo *echo.Echo) {
	echo.Static("/", "./public")
}
