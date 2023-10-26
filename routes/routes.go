package route

import (
	"rentcamp/config"
	"rentcamp/controller"
	"rentcamp/helper"

	"github.com/labstack/echo/v4"
)

func RouteAdmin(e *echo.Echo, uc controller.AdminControllerInterface, cfg config.Config) {
	var user = e.Group("/admins")
	user.POST("", uc.CreateUser())
	user.POST("/login", uc.Login())
}
func RouteProduct(e *echo.Echo, cpc controller.ProductControllerInterface, cfg config.Config) {
	var product = e.Group("/admins")
	product.Use(helper.Middleware())
	product.POST("/products", cpc.CreateProduct())
	product.GET("/products", cpc.GetAllProduct())
	product.PUT("/products/:id", cpc.UpdateProduct())
	product.DELETE("/products/:id", cpc.DeleteProduct())

	var admin = e.Group("/products")
	admin.PUT("/:id", cpc.UpdateProduct())
	admin.GET("", cpc.GetAllProduct())
	admin.GET("/:id", cpc.GetProductById())
}
