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
	product.PUT("/products/:id", cpc.UpdateProduct())
	product.DELETE("/products/:id", cpc.DeleteProduct())

	var admin = e.Group("/products")
	admin.GET("", cpc.GetAllProduct())
	admin.GET("/search", cpc.SearchProduct())
	admin.GET("/:id", cpc.GetProductById())
}

func RouteUser(e *echo.Echo, uc controller.UserControllerInterface, cfg config.Config) {
	var user = e.Group("/customer")
	user.Use(helper.Middleware())
	user.GET("/:id", uc.GetUserById())
	user.PUT("/:id", uc.UpdateUser())
	user.DELETE("/:id", uc.DeleteUser())

	var customer = e.Group("/customer")
	customer.POST("/login", uc.Login())
	customer.GET("", uc.GetAllUsers())
	customer.GET("/search", uc.SearchUsers())
	customer.POST("", uc.CreateUser())
}

func RouteCart(e *echo.Echo, cc controller.CartControllerInterface, cfg config.Config) {
	var cart = e.Group("/carts")
	cart.GET("/:cart_id", cc.GetCartByCartId())
	cart.POST("/:cart_id/items", cc.AddItemToCart())
	cart.PUT("/:cart_id/items/:item_id", cc.UpdateCartItem())
	cart.DELETE("/:cart_id/items/:item_id", cc.RemoveCartItem())
	cart.GET("/:cart_id/items", cc.GetItemsInCart())
	cart.DELETE("/:cart_id/items", cc.RemoveAllItemsFromCart())
	cart.GET("/:cart_id/total", cc.GetTotalCartPrice())
}
