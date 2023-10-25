package route

import (
	"rentcamp/config"
	"rentcamp/controller"

	"github.com/labstack/echo/v4"
)

func RouteAdmin(e *echo.Echo, uc controller.AdminControllerInterface, cfg config.Config) {
	var user = e.Group("/users")

	user.POST("", uc.CreateUser())
	user.POST("/login", uc.Login())
}
