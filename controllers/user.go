package controllers

import (
	model "e-GoLang/models"
	"encoding/json"
	"net/http"
)

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
