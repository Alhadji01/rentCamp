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
	"github.com/labstack/echo/v4"
)

type ProductControllerInterface interface {
	CreateProduct() echo.HandlerFunc
	GetAllProduct() echo.HandlerFunc
	GetProductById() echo.HandlerFunc
	UpdateProduct() echo.HandlerFunc
	DeleteProduct() echo.HandlerFunc
	SearchProduct() echo.HandlerFunc
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

func (cpc *ProductController) CreateProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input = model.Product{}
		if err := c.Bind(&input); err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Invalid category product input", nil))
		}
		image, err := c.FormFile("image")
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Image upload is required", nil))
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
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Failed to upload image to Cloudinary", nil))
		}
		input.Image = result.SecureURL

		createdProduct := cpc.model.InsertProduct(input)
		if createdProduct == nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Cannot process data, something happened", nil))
		}

		return c.JSON(http.StatusCreated, helper.FormatResponse("Success create category product", createdProduct))
	}
}

func (cpc *ProductController) GetAllProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		var res = cpc.model.SelectAll()

		if res == nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Error get all users, ", nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse("Success get all users, ", res))
	}
}

func (cpc *ProductController) SearchProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		page := c.QueryParam("page")
		limit := c.QueryParam("limit")
		search := c.QueryParam("name")

		pageNumber, err := strconv.Atoi(page)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Invalid page parameter", nil))
		}

		limitNumber, err := strconv.Atoi(limit)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Invalid limit parameter", nil))
		}

		res, totalCount, err := cpc.model.SelectWithPagination(pageNumber, limitNumber, search)

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
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("invalid category product input", nil))
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
