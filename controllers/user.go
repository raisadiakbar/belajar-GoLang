package controllers

import (
	model "e-GoLang/models"
	"encoding/json"
	"net/http"

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

func getAccountHandler(w http.ResponseWriter, r *http.Request) {
	// Mendapatkan user dari token JWT
	user := r.Context().Value("user").(model.User)

	// Menampilkan data user
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func updateAccountHandler(w http.ResponseWriter, r *http.Request) {
	// Mendapatkan user dari token JWT
	user := r.Context().Value("user").(model.User)

	// Mengambil data yang diberikan oleh user pada body request
	var updatedUser model.User
	err := json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Memperbarui data user
	user.Name = updatedUser.Name
	user.Email = updatedUser.Email
	user.Phone = updatedUser.Phone

	DB.Save(user)

	// Menampilkan data user yang telah diperbarui
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
