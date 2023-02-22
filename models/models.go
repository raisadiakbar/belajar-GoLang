package model

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

func ConnectDB() {
	var err error
	DB, err = gorm.Open("mysql", "root:Password@(localhost)/project-golang?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		fmt.Println("Failed to connect to database:", err)
	} else {
		fmt.Println("Successfully connected to database")
	}
}

func CloseDB() {
	DB.Close()
	fmt.Println("Successfully closed database connection")
}

type User struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email" gorm:"unique"`
	Password  string    `json:"-"`
	Phone     string    `json:"phone" gorm:"unique"`
	Address   []Address `json:"address,omitempty" gorm:"foreignkey:UserID"`
	Store     Store     `json:"store,omitempty" gorm:"foreignkey:UserID"`
	Role      string    `json:"role" gorm:"default:'user'"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Address struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	UserID    uint      `json:"user_id"`
	Name      string    `json:"name"`
	Street    string    `json:"street"`
	City      string    `json:"city"`
	Province  string    `json:"province"`
	Zipcode   string    `json:"zipcode"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Store struct {
	ID          uint      `gorm:"primary_key" json:"id"`
	UserID      uint      `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Category struct {
	ID          uint      `gorm:"primary_key" json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsAdmin     bool      `json:"is_admin"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Product struct {
	ID          uint      `gorm:"primary_key" json:"id"`
	UserID      uint      `json:"user_id"`
	CategoryID  uint      `json:"category_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       uint      `json:"price"`
	Image       string    `json:"image"`
	Stock       uint      `json:"stock"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Transaction struct {
	ID              uint      `gorm:"primary_key" json:"id"`
	UserID          uint      `json:"user_id"`
	ProductID       uint      `json:"product_id"`
	Quantity        uint      `json:"quantity"`
	TotalPrice      uint      `json:"total_price"`
	AddressID       uint      `json:"address_id"`
	TransactionTime time.Time `json:"transaction_time"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type LogProduct struct {
	ID            uint      `gorm:"primary_key" json:"id"`
	TransactionID uint      `json:"transaction_id"`
	ProductID     uint      `json:"product_id"`
	Quantity      uint      `json:"quantity"`
	Price         uint      `json:"price"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Deklarasi fungsi GetProductQueryParams
func GetProductQueryParams(r *http.Request) (int, int, error) {
	// Mendapatkan nilai dari parameter "limit"
	limitParam := r.URL.Query().Get("limit")

	// Mengkonversi nilai "limit" menjadi integer
	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		return 0, 0, err
	}

	// Mendapatkan nilai dari parameter "page"
	pageParam := r.URL.Query().Get("page")

	// Mengkonversi nilai "page" menjadi integer
	page, err := strconv.Atoi(pageParam)
	if err != nil {
		return 0, 0, err
	}

	return limit, page, nil
}

// UpdateAddress updates the given address in the database.
func UpdateAddress(id int, address *Address) error {
	return DB.Model(&Address{}).Where("id = ?", id).Updates(address).Error
}

var ErrAddressNotFound = errors.New("address not found")

func DeleteAddress(id int) error {
	result := DB.Delete(&Address{}, id)
	if result.Error != nil {
		if result.RecordNotFound() {
			return ErrAddressNotFound
		}
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrAddressNotFound
	}
	return nil
}

func UpdateCategory(id int, category *Category) error {
	result := DB.Model(&Category{}).Where("id = ?", id).Updates(category)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("category not found")
	}
	return nil
}

// DeleteCategory deletes a category with the given ID from the database.
func DeleteCategory(id int) error {
	result := DB.Delete(&Category{}, id)
	if result.Error != nil {
		if result.RecordNotFound() {
			return ErrCategoryNotFound
		}
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrCategoryNotFound
	}
	return nil
}

var ErrCategoryNotFound = errors.New("category not found")
