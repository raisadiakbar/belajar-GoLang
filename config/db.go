package config

import (
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
