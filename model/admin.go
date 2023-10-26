package model

import (
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Admin struct {
	Id        int    `gorm:"primaryKey;type:smallint" json:"id" form:"id"`
	Username  string `gorm:"type:varchar(25);not null" json:"username" form:"username"`
	Password  string `gorm:"type:varchar(25);not null" json:"password" form:"password"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Products  []Product      `json:"products"`
}

type Login struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

type AdminModelInterface interface {
	Login(username string, password string) *Admin
	Insert(newItem Admin) *Admin
}

type AdminsModel struct {
	db *gorm.DB
}

func NewUsersModel(db *gorm.DB) AdminModelInterface {
	return &AdminsModel{
		db: db,
	}
}

func (um *AdminsModel) Insert(newUser Admin) *Admin {
	if err := um.db.Create(&newUser).Error; err != nil {
		logrus.Error("Model : Insert data error, ", err.Error())
		return nil
	}

	return &newUser
}

func (um *AdminsModel) Login(username string, password string) *Admin {
	var data = Admin{}
	if err := um.db.Where("username = ?", username).First(&data).Error; err != nil {
		logrus.Error("Model : Login data error, ", err.Error())
		return nil
	}
	if data.Id == 0 {
		logrus.Error("Model : Login data error, ", nil)
		return nil
	}

	return &data
}
