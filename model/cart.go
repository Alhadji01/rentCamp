package model

import (
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Cart struct {
	ID        int            `gorm:"primaryKey" json:"id" form:"id"`
	UserID    int            `json:"user_id" form:"user_id"`
	CreatedAt time.Time      `gorm:"type:timestamp DEFAULT CURRENT_TIMESTAMP" json:"created_at" form:"created_at"`
	UpdatedAt time.Time      `gorm:"type:timestamp DEFAULT CURRENT_TIMESTAMP" json:"updated_at" form:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at" form:"deleted_at"`
	CartItem  []CartItem
}

type CartItem struct {
	ID        int            `gorm:"primaryKey" json:"id" form:"id"`
	CartID    int            `json:"cart_id" form:"cart_id"`
	ProductID int            `json:"product_id" form:"product_id"`
	Quantity  int            `json:"quantity" form:"quantity"`
	CreatedAt time.Time      `gorm:"type:timestamp DEFAULT CURRENT_TIMESTAMP" json:"created_at" form:"created_at"`
	UpdatedAt time.Time      `gorm:"type:timestamp DEFAULT CURRENT_TIMESTAMP" json:"updated_at" form:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at" form:"deleted_at"`
	Product   ProductResponse
}

type ProductResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	Stock       int    `json:"stock"`
	Image       string `json:"image"`
	AdminId     int    `json:"admin_id"`
}

func (ProductResponse) TableName() string {
	return "Products"
}

type CartModelInterface interface {
	GetCartByCartId(cartID int) (*Cart, error)
	AddItemToCart(cartID int, newItem CartItem) *CartItem
	UpdateCartItem(cartID, itemID int, updatedItem CartItem) *CartItem
	RemoveCartItem(cartID, itemID int) bool
	GetItemsInCart(cartID int) []CartItem
	RemoveAllItemsFromCart(cartID int) bool
	GetTotalCartPrice(cartID int) int
	CreateCart(userID int) (*Cart, error)
}

type CartModel struct {
	db *gorm.DB
}

func NewCartModel(db *gorm.DB) CartModelInterface {
	return &CartModel{
		db: db,
	}
}

func (cm *CartModel) CreateCart(userID int) (*Cart, error) {
	newCart := Cart{
		UserID: userID,
	}

	if err := cm.db.Create(&newCart).Error; err != nil {
		logrus.Error("Cart Model: Error creating cart, ", err.Error())
		return nil, err
	}

	return &newCart, nil
}

func (cm *CartModel) GetCartByCartId(cartID int) (*Cart, error) {
	var cart = Cart{}
	if err := cm.db.Preload("CartItem").Where("id = ?", cartID).First(&cart).Error; err != nil {
		logrus.Error("Cart Model: Error fetching cart, ", err.Error())
		return nil, err
	}
	return &cart, nil
}

func (cm *CartModel) AddItemToCart(cartID int, newItem CartItem) *CartItem {
	newItem.CartID = cartID
	if err := cm.db.Create(&newItem).Error; err != nil {
		logrus.Error("Cart Model: Error adding item to cart, ", err.Error())
		return nil
	}
	return &newItem
}

func (cm *CartModel) UpdateCartItem(cartID, itemID int, updatedItem CartItem) *CartItem {
	if err := cm.db.Where("cart_id = ? AND id = ?", cartID, itemID).Updates(&updatedItem).Error; err != nil {
		logrus.Error("Cart Model: Error updating cart item, ", err.Error())
		return nil
	}
	return &updatedItem
}

func (cm *CartModel) RemoveCartItem(cartID, itemID int) bool {
	if err := cm.db.Where("cart_id = ? AND id = ?", cartID, itemID).Delete(&CartItem{}).Error; err != nil {
		logrus.Error("Cart Model: Error removing cart item, ", err.Error())
		return false
	}
	return true
}

func (cm *CartModel) GetItemsInCart(cartID int) []CartItem {
	var items = []CartItem{}
	if err := cm.db.Preload("Product").Where("cart_id = ?", cartID).Find(&items).Error; err != nil {
		logrus.Error("Cart Model: Error fetching cart items, ", err.Error())
		return nil
	}
	return items
}

func (cm *CartModel) RemoveAllItemsFromCart(cartID int) bool {
	if err := cm.db.Where("cart_id = ?", cartID).Delete(&CartItem{}).Error; err != nil {
		logrus.Error("Cart Model: Error removing all cart items, ", err.Error())
		return false
	}
	return true
}

func (cm *CartModel) GetTotalCartPrice(cartID int) int {
	var totalPrice int
	if err := cm.db.Table("cart_items").Joins("JOIN products ON cart_items.product_id = products.id").Where("cart_id = ?", cartID).Select("SUM(cart_items.quantity * products.price)").Row().Scan(&totalPrice); err != nil {
		logrus.Error("Cart Model: Error calculating total cart price, ", err.Error())
		return 0
	}
	return totalPrice
}
