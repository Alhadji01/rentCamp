package controller

import (
	"context"
	"net/http"
	"rentcamp/config"
	"rentcamp/helper"
	"rentcamp/model"
	"strconv"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/golang-jwt/jwt/v5"
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

func NewProductControllerInterface(m model.ProductModelInterface, cfg config.Config) ProductControllerInterface {
	return &ProductController{
		model:  m,
		config: cfg,
	}
}

// Revisi: Authorization admin only
func (cpc *ProductController) CreateProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		userToken := c.Get("user").(*jwt.Token)
		tokenData, ok := helper.ExtractToken(userToken).(map[string]interface{})
		if !ok {
			return c.JSON(http.StatusUnauthorized, helper.FormatResponse("Invalid Token", "Invalid or missing token"))
		}
		role, roleOk := tokenData["role"].(string)
		if !roleOk {
			return c.JSON(http.StatusUnauthorized, helper.FormatResponse("Role information missing in the token", nil))
		}
		if role != "admin" {
			return c.JSON(http.StatusUnauthorized, helper.FormatResponse("You don't have permission", nil))
		}

		var input = model.Product{}
		if err := c.Bind(&input); err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Invalid product input", nil))
		}

		image, err := c.FormFile("image")
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Error uploading image", nil))
		}

		cld, err := cloudinary.NewFromParams(
			cpc.config.CDN_Cloud_Name,
			cpc.config.CDN_API_Key,
			cpc.config.CDN_API_Secret,
		)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Failed to initialize Cloudinary", nil))
		}

		uploadParams := uploader.UploadParams{
			Folder: cpc.config.CDN_Folder_Name,
		}

		fileReader, err := image.Open()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Failed to read image file", nil))
		}
		defer fileReader.Close()

		result, err := cld.Upload.Upload(context.TODO(), fileReader, uploadParams)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Failed to upload image to Cloudinary: "+err.Error(), nil))
		}

		input.Image = result.SecureURL

		createdProduct := cpc.model.InsertProduct(input)
		if createdProduct == nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Cannot process data, something happened", nil))
		}

		return c.JSON(http.StatusCreated, helper.FormatResponse("Success create product", createdProduct))
	}
}

// Revisi: Penggabungan getall dengan search pagination
func (cpc *ProductController) GetAllProduct() echo.HandlerFunc {
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
			var res = cpc.model.SelectAll()

			if res == nil {
				return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Error get all users, ", nil))
			}

			return c.JSON(http.StatusOK, helper.FormatResponse("Success get all users, ", res))
		} else {
			res, totalCount, err := cpc.model.SelectAllWithPagination(page, limit, search)

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

func (cpc *ProductController) GetProductById() echo.HandlerFunc {
	return func(c echo.Context) error {
		var paramId = c.Param("id")

		cnv, err := strconv.Atoi(paramId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Invalid id", nil))
		}

		var res = cpc.model.SelectById(cnv)
		if res == nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Error get product by id, ", nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse("Success get product", res))
	}
}

// Revisi: Authorization admin only
func (cpc *ProductController) UpdateProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		userToken := c.Get("user").(*jwt.Token)
		tokenData, ok := helper.ExtractToken(userToken).(map[string]interface{})
		if !ok {
			return c.JSON(http.StatusUnauthorized, helper.FormatResponse("Invalid Token", "Invalid or missing token"))
		}
		role, roleOk := tokenData["role"].(string)
		if !roleOk {
			return c.JSON(http.StatusUnauthorized, helper.FormatResponse("Role information missing in the token", nil))
		}
		if role != "admin" {
			return c.JSON(http.StatusUnauthorized, helper.FormatResponse("You don't have permission", nil))
		}

		var paramId = c.Param("id")
		cnv, err := strconv.Atoi(paramId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Invalid id", nil))
		}
		var input = model.Product{}
		if err := c.Bind(&input); err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("invalid product input", nil))
		}
		image, err := c.FormFile("image")
		if err != nil {
			existingProduct := cpc.model.SelectById(cnv)
			if existingProduct == nil {
				return c.JSON(http.StatusNotFound, helper.FormatResponse("Product not found", nil))
			}
			input.Image = existingProduct.Image
		} else {
			cld, err := cloudinary.NewFromParams(
				cpc.config.CDN_Cloud_Name,
				cpc.config.CDN_API_Key,
				cpc.config.CDN_API_Secret,
			)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Failed to initialize Cloudinary", nil))
			}
			uploadParams := uploader.UploadParams{
				Folder: cpc.config.CDN_Folder_Name,
			}

			fileReader, err := image.Open()
			if err != nil {
				return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Failed to read image file", nil))
			}
			defer fileReader.Close()

			result, err := cld.Upload.Upload(context.TODO(), fileReader, uploadParams)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Failed to upload image to Cloudinary", nil))
			}
			input.Image = result.SecureURL
		}

		input.Id = cnv

		var res = cpc.model.Update(input)
		if res == nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("cannot process data, something happened", nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse("Success update data", res))
	}
}

// Revisi: Authorization admin only
func (cpc *ProductController) DeleteProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		userToken := c.Get("user").(*jwt.Token)
		tokenData, ok := helper.ExtractToken(userToken).(map[string]interface{})
		if !ok {
			return c.JSON(http.StatusUnauthorized, helper.FormatResponse("Invalid Token", "Invalid or missing token"))
		}
		role, roleOk := tokenData["role"].(string)
		if !roleOk {
			return c.JSON(http.StatusUnauthorized, helper.FormatResponse("Role information missing in the token", nil))
		}
		if role != "admin" {
			return c.JSON(http.StatusUnauthorized, helper.FormatResponse("You don't have permission", nil))
		}
		var paramId = c.Param("id")

		cnv, err := strconv.Atoi(paramId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Invalid id", nil))
		}

		success := cpc.model.Delete(cnv)
		if !success {
			return c.JSON(http.StatusNotFound, helper.FormatResponse("category product not found", nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse("Success delete product", nil))
	}
}
