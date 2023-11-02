package controller

import (
	"net/http"
	"rentcamp/helper"
	"rentcamp/model"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type CartControllerInterface interface {
	GetCartByCartId() echo.HandlerFunc
	AddItemToCart() echo.HandlerFunc
	UpdateCartItem() echo.HandlerFunc
	RemoveCartItem() echo.HandlerFunc
	GetItemsInCart() echo.HandlerFunc
	RemoveAllItemsFromCart() echo.HandlerFunc
	GetTotalCartPrice() echo.HandlerFunc
	CreateCart() echo.HandlerFunc
}

type CartController struct {
	model model.CartModelInterface
}

func NewCartControllerInterface(m model.CartModelInterface) CartControllerInterface {
	return &CartController{
		model: m,
	}
}

// revisi : update authorization
func (cc *CartController) CreateCart() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user")
		if user == nil {
			return c.JSON(http.StatusUnauthorized, helper.FormatResponse("Unauthorized 1", nil))
		}
		claims, ok := user.(*jwt.Token).Claims.(jwt.MapClaims)
		if !ok {
			return c.JSON(http.StatusUnauthorized, helper.FormatResponse("Unauthorized 2", nil))
		}
		userID, ok := claims["id"].(float64)
		if !ok {
			return c.JSON(http.StatusUnauthorized, helper.FormatResponse("Unauthorized 3", nil))
		}

		newCart, err := cc.model.CreateCart(int(userID))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Error creating cart", nil))
		}

		return c.JSON(http.StatusCreated, helper.FormatResponse("Cart created successfully", newCart))
	}
}

func (cc *CartController) GetCartByCartId() echo.HandlerFunc {
	return func(c echo.Context) error {
		var paramCartID = c.Param("cart_id")

		cartID, err := strconv.Atoi(paramCartID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Invalid cart ID", nil))
		}

		res, err := cc.model.GetCartByCartId(cartID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Error fetching cart", err))
		}

		if res == nil {
			return c.JSON(http.StatusNotFound, helper.FormatResponse("Cart not found", nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse("Cart retrieved successfully", res))
	}
}

func (cc *CartController) AddItemToCart() echo.HandlerFunc {
	return func(c echo.Context) error {
		var paramCartID = c.Param("cart_id")

		cartID, err := strconv.Atoi(paramCartID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Invalid cart ID", nil))
		}

		var input = model.CartItem{}
		if err := c.Bind(&input); err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Invalid cart item input", nil))
		}

		var res = cc.model.AddItemToCart(cartID, input)
		if res == nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Error adding item to cart", nil))
		}

		return c.JSON(http.StatusCreated, helper.FormatResponse("Item added to cart successfully", res))
	}
}

func (cc *CartController) UpdateCartItem() echo.HandlerFunc {
	return func(c echo.Context) error {
		var paramCartID = c.Param("cart_id")
		var paramItemID = c.Param("item_id")

		cartID, err := strconv.Atoi(paramCartID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Invalid cart ID", nil))
		}

		itemID, err := strconv.Atoi(paramItemID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Invalid item ID", nil))
		}

		var input = model.CartItem{}
		if err := c.Bind(&input); err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Invalid cart item input", nil))
		}

		input.ID = itemID

		var res = cc.model.UpdateCartItem(cartID, itemID, input)
		if res == nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Error updating cart item", nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse("Cart item updated successfully", res))
	}
}

func (cc *CartController) RemoveCartItem() echo.HandlerFunc {
	return func(c echo.Context) error {
		var paramCartID = c.Param("cart_id")
		var paramItemID = c.Param("item_id")

		cartID, err := strconv.Atoi(paramCartID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Invalid cart ID", nil))
		}

		itemID, err := strconv.Atoi(paramItemID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Invalid item ID", nil))
		}

		success := cc.model.RemoveCartItem(cartID, itemID)
		if !success {
			return c.JSON(http.StatusNotFound, helper.FormatResponse("Cart item not found", nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse("Cart item removed successfully", nil))
	}
}

func (cc *CartController) GetItemsInCart() echo.HandlerFunc {
	return func(c echo.Context) error {
		var paramCartID = c.Param("cart_id")

		cartID, err := strconv.Atoi(paramCartID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Invalid cart ID", nil))
		}

		var res = cc.model.GetItemsInCart(cartID)
		if res == nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Error fetching cart items", nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse("Cart items retrieved successfully", res))
	}
}

func (cc *CartController) RemoveAllItemsFromCart() echo.HandlerFunc {
	return func(c echo.Context) error {
		var paramCartID = c.Param("cart_id")

		cartID, err := strconv.Atoi(paramCartID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Invalid cart ID", nil))
		}

		success := cc.model.RemoveAllItemsFromCart(cartID)
		if !success {
			return c.JSON(http.StatusNotFound, helper.FormatResponse("No items found in the cart", nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse("All items removed from the cart", nil))
	}
}

func (cc *CartController) GetTotalCartPrice() echo.HandlerFunc {
	return func(c echo.Context) error {
		var paramCartID = c.Param("cart_id")

		cartID, err := strconv.Atoi(paramCartID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Invalid cart ID", nil))
		}

		var totalPrice = cc.model.GetTotalCartPrice(cartID)
		if totalPrice == 0 {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Error calculating total cart price", nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse("Total cart price calculated successfully", totalPrice))
	}
}
