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
	user.POST("/login", uc.Login())
	user.POST("", uc.CreateUser())
	user.GET("", uc.GetAllUsers())
	user.GET("/:id", uc.GetUserById())
	user.PUT("/:id", uc.UpdateUser())
	user.DELETE("/:id", uc.DeleteUser())
	user.GET("/search", uc.SearchUsers())
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
