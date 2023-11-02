package controller

import (
	"errors"
	"net/http"
	"rentcamp/config"
	"rentcamp/helper"
	"rentcamp/model"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserControllerInterface interface {
	CreateUser() echo.HandlerFunc
	Login() echo.HandlerFunc
	GetAllUsers() echo.HandlerFunc
	GetUserById() echo.HandlerFunc
	UpdateUser() echo.HandlerFunc
	DeleteUser() echo.HandlerFunc
}

type UserController struct {
	config config.Config
	model  model.UserModelInterface
}

func NewUserControlInterface(m model.UserModelInterface) UserControllerInterface {
	return &UserController{
		model: m,
	}
}

func (uc *UserController) CreateUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input = model.User{}
		if err := c.Bind(&input); err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("invalid user input", nil))
		}
		hashpwd, err := helper.HashPassword(input.Password)
		if err != nil {
			return errors.New("Gagal mengenkripsi kata sandi: " + err.Error())
		}
		var newUser = model.User{}
		newUser.Name = input.Name
		newUser.Username = input.Username
		newUser.Password = hashpwd
		newUser.Email = input.Email
		newUser.Phone = input.Phone
		newUser.Address = input.Address
		newUser.Gender = input.Gender

		var res = uc.model.InsertUser(newUser)
		if res == nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Cannot process data, something happend", nil))
		}

		return c.JSON(http.StatusCreated, helper.FormatResponse("success create user", res))
	}
}

func (uc *UserController) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input = model.User{}

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
		var role = "user"
		var jwtToken = helper.GenerateJWT(uc.config.Secret, uc.config.RefreshSecret, res.Id, res.Username, role)

		if jwtToken == nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Cannot process data, something happend", nil))
		}
		var info = map[string]any{}
		info["id"] = res.Id
		info["name"] = res.Name
		info["username"] = res.Username

		jwtToken["info"] = info

		return c.JSON(http.StatusOK, helper.FormatResponse("login success", jwtToken))
	}
}

func (uc *UserController) GetAllUsers() echo.HandlerFunc {
	return func(c echo.Context) error {
		pageStr := c.QueryParam("page")
		limitStr := c.QueryParam("limit")
		search := c.QueryParam("name")

		var page, limit int
		if pageStr != "" {
			page, _ = strconv.Atoi(pageStr)
		}
		if limitStr != "" {
			limit, _ = strconv.Atoi(limitStr)
		}

		if page == 0 || limit == 0 {
			var res = uc.model.SelectAll()

			if res == nil {
				return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Error get all users, ", nil))
			}

			return c.JSON(http.StatusOK, helper.FormatResponse("Success get all users, ", res))
		} else {
			res, totalCount, err := uc.model.SelectAllWithPagination(page, limit, search)

			if err != nil {
				return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Error fetching products", nil))
			}

			response := map[string]interface{}{
				"products":   res,
				"total_data": totalCount,
			}

			return c.JSON(http.StatusOK, helper.FormatResponse("Success fetching products", response))
		}
	}
}

func (uc *UserController) GetUserById() echo.HandlerFunc {
	return func(c echo.Context) error {
		var paramId = c.Param("id")

		cnv, err := strconv.Atoi(paramId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Invalid id", nil))
		}

		var res = uc.model.SelectById(cnv)
		if res == nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Error get user by id, ", nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse("Success get user", res))
	}
}

func (uc *UserController) UpdateUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		role := claims["role"].(string)
		userID := int(claims["id"].(float64))

		paramID := c.Param("id")
		id, err := strconv.Atoi(paramID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Invalid id", nil))
		}

		if role != "admin" && userID != id {
			return c.JSON(http.StatusForbidden, helper.FormatResponse("Permission denied. You don't have the required permissions.", nil))
		}

		var input = model.User{}
		if err := c.Bind(&input); err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Invalid user input", nil))
		}

		input.Id = id
		hashpwd, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			return errors.New("Gagal mengenkripsi kata sandi: " + err.Error())
		}
		input.Password = string(hashpwd)

		res, err := uc.model.Update(input)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Cannot process data, something happened: "+err.Error(), nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse("Success update data", res))
	}
}

func (uc *UserController) DeleteUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		role := claims["role"].(string)
		userID := int(claims["id"].(float64))

		paramID := c.Param("id")
		id, err := strconv.Atoi(paramID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Invalid id", nil))
		}

		if role != "admin" && userID != id {
			return c.JSON(http.StatusForbidden, helper.FormatResponse("Permission denied. You don't have the required permissions.", nil))
		}
		var paramId = c.Param("id")

		cnv, err := strconv.Atoi(paramId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Invalid id", nil))
		}

		success := uc.model.Delete(cnv)
		if !success {
			return c.JSON(http.StatusNotFound, helper.FormatResponse("User not found", nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse("Success delete user", nil))
	}
}
