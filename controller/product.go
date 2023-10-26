package controller

import (
	"net/http"
	"rentcamp/config"
	"rentcamp/helper"
	"rentcamp/model"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ProductControllerInterface interface {
	CreateProduct() echo.HandlerFunc
	GetAllProduct() echo.HandlerFunc
	GetProductById() echo.HandlerFunc
	UpdateProduct() echo.HandlerFunc
	DeleteProduct() echo.HandlerFunc
}

type ProductController struct {
	config config.Config
	model  model.ProductModelInterface
}

func NewProductControllerInterface(m model.ProductModelInterface) ProductControllerInterface {
	return &ProductController{
		model: m,
	}
}

func (cpc *ProductController) CreateProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input = model.Product{}
		if err := c.Bind(&input); err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("invalid egory product input", nil))
		}

		var res = cpc.model.Insert(input)
		if res == nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Cannot process data, something happend", nil))
		}

		return c.JSON(http.StatusCreated, helper.FormatResponse("success create egory product", res))
	}
}

func (cpc *ProductController) GetAllProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		var res = cpc.model.SelectAll()

		if res == nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Error get all egory product, ", nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse("Success get all egory product, ", res))
	}
}

func (cpc *ProductController) GetProductById() echo.HandlerFunc {
	return func(c echo.Context) error {
		var paramId = c.Param("id")

		cnv, err := strconv.Atoi(paramId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Invalid id", nil))
		}

		var res = cpc.model.SelectById(cnv)
		if res == nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Error get egory product by id, ", nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse("Success get egory product", res))
	}
}

func (cpc *ProductController) UpdateProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		var paramId = c.Param("id")
		cnv, err := strconv.Atoi(paramId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Invalid id", nil))
		}

		var input = model.Product{}
		if err := c.Bind(&input); err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("invalid egory product input", nil))
		}

		input.Id = cnv

		var res = cpc.model.Update(input)
		if res == nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("cannot process data, something happend", nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse("Success update data", res))
	}
}

func (cpc *ProductController) DeleteProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		var paramId = c.Param("id")

		cnv, err := strconv.Atoi(paramId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Invalid id", nil))
		}

		success := cpc.model.Delete(cnv)
		if !success {
			return c.JSON(http.StatusNotFound, helper.FormatResponse("egory product not found", nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse("Success delete egory product", nil))
	}
}
