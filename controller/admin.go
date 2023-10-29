package controller

import (
	"net/http"
	"rentcamp/config"
	"rentcamp/helper"
	"rentcamp/model"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type AdminControllerInterface interface {
	CreateUser() echo.HandlerFunc
	Login() echo.HandlerFunc
}

type AdminController struct {
	config config.Config
	model  model.AdminModelInterface
}

func NewAdminControlInterface(m model.AdminModelInterface) AdminControllerInterface {
	return &AdminController{
		model: m,
	}
}

func (uc *AdminController) CreateUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input = model.Admin{}
		if err := c.Bind(&input); err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("invalid user input", nil))
		}

		var res = uc.model.Insert(input)
		if res == nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Cannot process data, something happend", nil))
		}

		return c.JSON(http.StatusCreated, helper.FormatResponse("success create user", res))
	}
}

func (uc *AdminController) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input = model.Login{}

		if err := c.Bind(&input); err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Invalid user input", nil))
		}

		var res = uc.model.Login(input.Username, input.Password)

		if res == nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Cannot process data, something happend", nil))
		}

		if res.Id == 0 {
			return c.JSON(http.StatusNotFound, helper.FormatResponse("Data not found", nil))
		}

		var jwtToken = helper.GenerateJWT(uc.config.Secret, uc.config.RefreshSecret, res.Id, res.Username, res.Role)

		if jwtToken == nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Cannot process data, something happend", nil))
		}

		var info = map[string]any{}
		info["username"] = res.Username
		info["role"] = res.Role

		jwtToken["info"] = info

		return c.JSON(http.StatusOK, helper.FormatResponse("login success", jwtToken))
	}
}

func SomeSecureHandler(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["username"].(string)
	return c.String(http.StatusOK, "Welcome, "+username+"!")
}
