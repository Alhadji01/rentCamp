package main

import (
	"fmt"
	"rentcamp/config"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	var config = config.InitConfig()

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.ServerPort)).Error())
}
