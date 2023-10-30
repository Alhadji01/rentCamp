package main

import (
	"fmt"
	"rentcamp/config"
	"rentcamp/controller"
	"rentcamp/model"
	route "rentcamp/routes"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

func main() {
	var err = godotenv.Load(".ENV")
	if err != nil {
		logrus.Error("Config : Cannot load config file, ", err.Error())
	}
	e := echo.New()
	var config = config.InitConfig()

	db := model.InitModel(*config)
	model.Migrate(db)

	adminModel := model.NewAdminsModel(db)
	ProductModel := model.NewProductsModel(db)
	userModel := model.NewUsersModel(db)
	cartModel := model.NewCartModel(db)

	adminController := controller.NewAdminControlInterface(adminModel)
	ProductController := controller.NewProductControllerInterface(ProductModel, *config)
	userController := controller.NewUserControlInterface(userModel)
	cartController := controller.NewCartControllerInterface(cartModel)

	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(middleware.CORS())
	e.Use(middleware.LoggerWithConfig(
		middleware.LoggerConfig{
			Format: "method=${method}, uri=${uri}, status=${status}, time=${time_rfc3339}\n",
		}))

	route.RouteAdmin(e, adminController, *config)
	route.RouteProduct(e, ProductController, *config)
	route.RouteUser(e, userController, *config)
	route.RouteCart(e, cartController, *config)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.ServerPort)).Error())
}
