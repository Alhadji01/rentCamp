package model

import (
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Id        int            `gorm:"primaryKey;type:smallint" json:"id" form:"id"`
	Name      string         `gorm:"type:varchar(100);not null" json:"name" form:"name"`
	Username  string         `gorm:"type:varchar(25);not null" json:"username" form:"username"`
	Password  string         `gorm:"type:varchar(255);not null" json:"password" form:"password"`
	Email     string         `gorm:"type:varchar(50);not null" json:"email" form:"email"`
	Phone     string         `gorm:"type:varchar(15);not null" json:"phone" form:"phone"`
	Address   string         `gorm:"type:varchar(255);not null" json:"address" form:"address"`
	Gender    string         `gorm:"type:ENUM('m','f');not null" json:"gender" form:"gender"`
	CreatedAt time.Time      `gorm:"type:timestamp DEFAULT CURRENT_TIMESTAMP" json:"created_at" form:"created_at"`
	UpdatedAt time.Time      `gorm:"type:timestamp DEFAULT CURRENT_TIMESTAMP" json:"updated_at" form:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at" form:"deleted_at"`
	Carts     []Cart         `json:"cart"`
}

type LoginUser struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

type UserModelInterface interface {
	Login(username string, password string) *User
	InsertUser(newItem User) *User
	SelectAll() []User
	SelectAllWithPagination(page, limit int, search string) ([]User, int64, error)
	SelectById(userId int) *User
	Update(updatedData User) (*User, error)
	Delete(userId int) bool
}

type UsersModel struct {
	db *gorm.DB
}

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

func (um *UsersModel) InsertUser(newUser User) *User {
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

func (um *UsersModel) SelectAllWithPagination(page, limit int, search string) ([]User, int64, error) {
	var Users []User
	var totalCount int64

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	if search != "" {
		if err := um.db.Model(&User{}).
			Where("name LIKE ?", "%"+search+"%").
			Count(&totalCount).Error; err != nil {
			return nil, 0, err
		}
	} else {
		if err := um.db.Model(&User{}).
			Count(&totalCount).Error; err != nil {
			return nil, 0, err
		}
	}

	offset := (page - 1) * limit

	if offset < 0 {
		offset = 0
	}

	if search != "" {
		if err := um.db.Where("name LIKE ?", "%"+search+"%").
			Offset(offset).
			Limit(limit).
			Find(&Users).Error; err != nil {
			return nil, 0, err
		}
	} else {
		if err := um.db.Offset(offset).
			Limit(limit).
			Find(&Users).Error; err != nil {
			return nil, 0, err
		}
	}

	return Users, totalCount, nil
}

func (um *UsersModel) SelectById(userId int) *User {
	var data = User{}
	if err := um.db.Preload("cartz").Where("id = ?", userId).First(&data).Error; err != nil {
		logrus.Error("Model : Data with that ID was not found, ", err.Error())
		return nil
	}

	return &data
}

func (um *UsersModel) Update(updatedData User) (*User, error) {
	var data map[string]interface{} = make(map[string]interface{})

	if updatedData.Name != "" {
		data["name"] = updatedData.Name
	}
	if updatedData.Username != "" {
		data["username"] = updatedData.Username
	}
	if updatedData.Password != "" {
		hashpwd, err := bcrypt.GenerateFromPassword([]byte(updatedData.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, errors.New("Gagal mengenkripsi kata sandi: " + err.Error())
		}
		data["password"] = string(hashpwd)
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
		return nil, errors.New("Gagal memperbarui data pengguna: " + err.Error())
	}

	if dataCount := qry.RowsAffected; dataCount < 1 {
		return nil, errors.New("Gagal memperbarui data pengguna: tidak ada data yang terpengaruh")
	}

	var updatedUser = User{}
	if err := um.db.Where("id = ?", updatedData.Id).First(&updatedUser).Error; err != nil {
		return nil, errors.New("Gagal mengambil data pengguna yang diperbarui: " + err.Error())
	}

	return &updatedUser, nil
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
