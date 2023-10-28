package model

import (
	"time"

	"gorm.io/gorm"
)

type Cart struct {
	ID        int `gorm:"primaryKey" json:"id" form:"id"`
	UserId    int `json:"user_id" form:"user_id"`
	ProductId int `json:"product_id" form:"product_id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at" form:"deleted_at"`
}
