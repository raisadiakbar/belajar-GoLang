package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"e-GoLang/errors"

	models "e-GoLang/models"

	"github.com/gorilla/mux"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// createAddressHandler is the handler function for the POST /api/addresses route.
func createAddressHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var address models.Address
	err := json.NewDecoder(r.Body).Decode(&address)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create address
	err = models.DB.Create(&address).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return created address
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(address)
}

// getAddressHandler is the handler function for the GET /api/addresses/{id} route.
func getAddressHandler(w http.ResponseWriter, r *http.Request) {
	// Get address ID from route parameters
	vars := mux.Vars(r)
	addressID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get address by ID
	var address models.Address
	err = models.DB.First(&address, addressID).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Return address
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(address)
}

func updateAddressHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	addressID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid address ID", http.StatusBadRequest)
		return
	}

	var address models.Address
	err = json.NewDecoder(r.Body).Decode(&address)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = models.UpdateAddress(addressID, &address)
	if err != nil {
		if errors.Is(err, models.ErrAddressNotFound) {
			http.Error(w, "Address not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to update address", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteAddressHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	addressID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid address ID", http.StatusBadRequest)
		return
	}

	err = models.DeleteAddress(addressID)
	if err != nil {
		if errors.Is(err, models.ErrAddressNotFound) {
			http.Error(w, "Address not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete address", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
