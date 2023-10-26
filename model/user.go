package model

import (
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type User struct {
	Id        int            `gorm:"primaryKey;type:smallint" json:"id" form:"id"`
	Name      string         `gorm:"type:varchar(100);not null" json:"name" form:"name"`
	Username  string         `gorm:"type:varchar(25);not null" json:"username" form:"username"`
	Password  string         `gorm:"type:varchar(25);not null" json:"password" form:"password"`
	Email     string         `gorm:"type:varchar(50);not null" json:"email" form:"email"`
	Phone     string         `gorm:"type:varchar(15);not null" json:"phone" form:"phone"`
	Address   string         `gorm:"type:varchar(255);not null" json:"address" form:"address"`
	Gender    string         `gorm:"type:ENUM('m','f');not null" json:"gender" form:"gender"`
	CreatedAt time.Time      `gorm:"type:timestamp DEFAULT CURRENT_TIMESTAMP" json:"created_at" form:"created_at"`
	UpdatedAt time.Time      `gorm:"type:timestamp DEFAULT CURRENT_TIMESTAMP" json:"updated_at" form:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at" form:"deleted_at"`
}

// Declaration Contract
type UserModelInterface interface {
	Login(username string, password string) *User
	Insert(newItem User) *User
	SelectAll() []User
	SelectById(userId int) *User
	Update(updatedData User) *User
	Delete(userId int) bool
}

// Interaction with db
type UsersModel struct {
	db *gorm.DB
}

// New Instance from UsersModel
func NewUsersModel(db *gorm.DB) UserModelInterface {
	return &UsersModel{
		db: db,
	}
}

func (um *UsersModel) Login(username string, password string) *User {
	var data = User{}
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

func (um *UsersModel) Insert(newUser User) *User {
	if err := um.db.Create(&newUser).Error; err != nil {
		logrus.Error("Model : Insert data error, ", err.Error())
		return nil
	}

	return &newUser
}

func (um *UsersModel) SelectAll() []User {
	var data = []User{}
	if err := um.db.Find(&data).Error; err != nil {
		logrus.Error("Model : Cannot get all users, ", err.Error())
		return nil
	}

	return data
}

func (um *UsersModel) SelectById(userId int) *User {
	var data = User{}
	if err := um.db.Where("id = ?", userId).First(&data).Error; err != nil {
		logrus.Error("Model : Data with that ID was not found, ", err.Error())
		return nil
	}

	return &data
}

func (um *UsersModel) Update(updatedData User) *User {
	var data map[string]interface{} = make(map[string]interface{})

	if updatedData.Name != "" {
		data["name"] = updatedData.Name
	}
	if updatedData.Username != "" {
		data["username"] = updatedData.Username
	}
	if updatedData.Password != "" {
		data["password"] = updatedData.Password
	}
	if updatedData.Email != "" {
		data["email"] = updatedData.Email
	}
	if updatedData.Phone != "" {
		data["phone"] = updatedData.Phone
	}
	if updatedData.Address != "" {
		data["address"] = updatedData.Address
	}
	if updatedData.Gender != "" {
		data["gender"] = updatedData.Gender
	}

	var qry = um.db.Table("users").Where("id = ?", updatedData.Id).Updates(data)
	if err := qry.Error; err != nil {
		logrus.Error("Model : update error, ", err.Error())
		return nil
	}

	if dataCount := qry.RowsAffected; dataCount < 1 {
		logrus.Error("Model : Update error, ", "no data effected")
		return nil
	}

	var updatedUser = User{}
	if err := um.db.Where("id = ?", updatedData.Id).First(&updatedUser).Error; err != nil {
		logrus.Error("Model : Error get updated data, ", err.Error())
		return nil
	}

	return &updatedUser
}

func (um *UsersModel) Delete(userId int) bool {
	var data = User{}
	data.Id = userId

	if err := um.db.Where("id = ?", userId).First(&data).Error; err != nil {
		logrus.Error("Model: Error finding data to delete, ", err.Error())
		return false
	}

	if err := um.db.Delete(&data).Error; err != nil {
		logrus.Error("Model : Error delete data, ", err.Error())
		return false
	}

	return true
}
