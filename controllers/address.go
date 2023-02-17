package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/raisadiakbar/belajar-GoLang/models"
	"github.com/raisadiakbar/belajar-GoLang/repository"
)

type AddressController struct{}

var addressRepository = repository.AddressRepository{}

// RespondWithError writes an error response to the http.ResponseWriter
func respondWithError(w http.ResponseWriter, status int, message string) {
	response := map[string]string{"error": message}
	responseJSON, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(responseJSON)
}

// RespondWithJSON writes a JSON response to the http.ResponseWriter
func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

// Create a new address
func (ac AddressController) CreateAddress(w http.ResponseWriter, r *http.Request) {
	address := models.Address{}
	err := json.NewDecoder(r.Body).Decode(&address)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err = address.Validate(); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err = addressRepository.Create(&address); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, address)
}

// Get an address by ID
func (ac AddressController) GetAddress(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid address ID")
		return
	}

	address, err := addressRepository.FindByID(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Address not found")
		return
	}

	respondWithJSON(w, http.StatusOK, address)
}

func createAddressHandler(w http.ResponseWriter, r *http.Request) {
	ac := AddressController{}
	ac.CreateAddress(w, r)
}

func getAddressHandler(w http.ResponseWriter, r *http.Request) {
	ac := AddressController{}
	ac.GetAddress(w, r)
}
