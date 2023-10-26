package main

import (
	"fmt"
	"rentcamp/config"
	"rentcamp/controller"
	"rentcamp/model"
	route "rentcamp/routes"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	var config = config.InitConfig()

	db := model.InitModel(*config)
	model.Migrate(db)

	adminModel := model.NewAdminsModel(db)
	ProductModel := model.NewProductsModel(db)

	adminController := controller.NewAdminControlInterface(adminModel)
	ProductController := controller.NewProductControllerInterface(ProductModel)

	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(middleware.CORS())
	e.Use(middleware.LoggerWithConfig(
		middleware.LoggerConfig{
			Format: "method=${method}, uri=${uri}, status=${status}, time=${time_rfc3339}\n",
		}))

	route.RouteAdmin(e, adminController, *config)
	route.RouteProduct(e, ProductController, *config)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.ServerPort)).Error())
}
