package db

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

func connectDB() (*gorm.DB, error) {
	db, err := gorm.Open("mysql", "root:Password@(localhost)/project-golang?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		return nil, err
	}

	return db, nil
}

func CloseDB() {
	DB.Close()
	fmt.Println("Successfully closed database connection")
}
