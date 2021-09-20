package server

import (
	echo "github.com/labstack/echo/v4"
)

func serveStaticFiles(echo *echo.Echo) {
	echo.Static("/", "./public")
}
