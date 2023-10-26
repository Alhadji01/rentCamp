package model

import (
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Product struct {
	Id          int            `gorm:"primaryKey;type:smallint" json:"id" form:"id"`
	Name        string         `gorm:"type:varchar(100);not null" json:"name" form:"name"`
	Image       string         `gorm:"type:text"`
	Description string         `gorm:"type:text;not null" json:"description" form:"description"`
	Price       int            `gorm:"type:smallint;not null" json:"unit_price" form:"unit_price"`
	Stock       int            `gorm:"type:smallint;not null" json:"stock" form:"stock"`
	CreatedAt   time.Time      `gorm:"type:timestamp DEFAULT CURRENT_TIMESTAMP" json:"created_at" form:"created_at"`
	UpdatedAt   time.Time      `gorm:"type:timestamp DEFAULT CURRENT_TIMESTAMP" json:"updated_at" form:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at" form:"deleted_at"`
	AdminId     int            `json:"admin_id" form:"admin_id"`
}

type ProductModelInterface interface {
	Insert(newProduct Product) *Product
	SelectAll() []Product
	SelectById(ProductId int) *Product
	Update(updatedData Product) *Product
	Delete(ProductId int) bool
}

type ProductsModel struct {
	db *gorm.DB
}

func NewProductsModel(db *gorm.DB) ProductModelInterface {
	return &ProductsModel{
		db: db,
	}
}

func (cpm *ProductsModel) Insert(newProduct Product) *Product {
	if err := cpm.db.Create(&newProduct).Error; err != nil {
		logrus.Error("Model : Insert data error, ", err.Error())
		return nil
	}

	return &newProduct
}

func (cpm *ProductsModel) SelectAll() []Product {
	var data = []Product{}
	if err := cpm.db.Find(&data).Error; err != nil {
		logrus.Error("Model : Cannot get all category product, ", err.Error())
		return nil
	}

	return data
}

func (cpm *ProductsModel) SelectById(ProductId int) *Product {
	var data = Product{}
	if err := cpm.db.Where("id = ?", ProductId).First(&data).Error; err != nil {
		logrus.Error("Model : Data with that ID was not found, ", err.Error())
		return nil
	}

	return &data
}

func (cpm *ProductsModel) Update(updatedData Product) *Product {
	var qry = cpm.db.Table("products").Where("id = ?", updatedData.Id).Update("name", updatedData.Name)
	if err := qry.Error; err != nil {
		logrus.Error("Model : update error, ", err.Error())
		return nil
	}

	if dataCount := qry.RowsAffected; dataCount < 1 {
		logrus.Error("Model : Update error, ", "no data effected")
		return nil
	}

	var updatedProduct = Product{}
	if err := cpm.db.Where("id = ?", updatedData.Id).First(&updatedProduct).Error; err != nil {
		logrus.Error("Model : Error get updated data, ", err.Error())
		return nil
	}

	return &updatedProduct
}

func (cpm *ProductsModel) Delete(ProductId int) bool {
	var data = Product{}
	data.Id = ProductId

	if err := cpm.db.Where("id = ?", ProductId).First(&data).Error; err != nil {
		logrus.Error("Model: Error finding data to delete, ", err.Error())
		return false
	}

	if err := cpm.db.Delete(&data).Error; err != nil {
		logrus.Error("Model : Error delete data, ", err.Error())
		return false
	}

	return true
}
