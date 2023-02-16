package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/myusername/myapp/models"
)

// CreateAddressRequest represents the request body for creating an address.
type CreateAddressRequest struct {
	RecipientName  string `json:"recipient_name"`
	RecipientPhone string `json:"recipient_phone"`
	StreetAddress  string `json:"street_address"`
	City           string `json:"city"`
	Province       string `json:"province"`
	ZipCode        string `json:"zip_code"`
}

// AddressResponse represents the response body for an address.
type AddressResponse struct {
	ID             int64  `json:"id"`
	RecipientName  string `json:"recipient_name"`
	RecipientPhone string `json:"recipient_phone"`
	StreetAddress  string `json:"street_address"`
	City           string `json:"city"`
	Province       string `json:"province"`
	ZipCode        string `json:"zip_code"`
	UserID         int64  `json:"user_id"`
}

// createAddressHandler creates a new address.
func createAddressHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := extractUserIDFromToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var req CreateAddressRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	address := Address{
		RecipientName:  req.RecipientName,
		RecipientPhone: req.RecipientPhone,
		StreetAddress:  req.StreetAddress,
		City:           req.City,
		Province:       req.Province,
		ZipCode:        req.ZipCode,
		UserID:         userID,
	}
	err = db.CreateAddress(&address)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := AddressResponse{
		ID:             address.ID,
		RecipientName:  address.RecipientName,
		RecipientPhone: address.RecipientPhone,
		StreetAddress:  address.StreetAddress,
		City:           address.City,
		Province:       address.Province,
		ZipCode:        address.ZipCode,
		UserID:         address.UserID,
	}
	json.NewEncoder(w).Encode(res)
}

// getAddressHandler gets an address by ID.
func getAddressHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := extractUserIDFromToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	params := mux.Vars(r)
	id, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid address ID", http.StatusBadRequest)
		return
	}

	address, err := db.GetAddress(id)
	if err != nil {
		if errors.Is(err, ErrAddressNotFound) {
			http.Error(w, "Address not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if address.UserID != userID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	res := AddressResponse{
		ID:             address.ID,
		RecipientName:  address.RecipientName,
		RecipientPhone: address.RecipientPhone,
		StreetAddress:  address.StreetAddress,
		City:           address.City,
		Province:       address.Province,
		ZipCode:        address.ZipCode,
		UserID:         address.UserID,
	}
	json.NewEncoder(w).Encode(res)
}

// Update an address by id
func updateAddressHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid address ID", http.StatusBadRequest)
		return
	}

	// Get address from database
	address, err := models.GetAddressById(id)
	if err != nil {
		http.Error(w, "Address not found", http.StatusNotFound)
		return
	}

	// Parse request body to update address
	err = json.NewDecoder(r.Body).Decode(&address)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate and update address in database
	err = models.ValidateAddress(&address)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = models.UpdateAddress(&address)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return updated address as JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(address)
}

// Delete an address by id
func deleteAddressHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid address ID", http.StatusBadRequest)
		return
	}

	// Delete address from database
	err = models.DeleteAddress(id)
	if err != nil {
		http.Error(w, "Address not found", http.StatusNotFound)
		return
	}

	// Return success message as JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Address deleted successfully"})
}
