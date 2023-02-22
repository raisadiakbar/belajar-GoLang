package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"e-GoLang/models"
)

func RespondWithJSON(w http.ResponseWriter, payload interface{}) {
	response, err := json.Marshal(payload)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error while encoding response to JSON"))
		return
	}

	w.Write(response)
}

// Create a new product
func createProductHandler(w http.ResponseWriter, r *http.Request) {
	// Validate user input
	var newProduct models.Product
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newProduct)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err = models.ValidateProduct(&newProduct); err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Create the new product
	if err := models.CreateProduct(&newProduct); err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Set response status code
	w.WriteHeader(http.StatusOK)
}

// Get a list of all products
func getProductListHandler(w http.ResponseWriter, r *http.Request) {
	limit, page, err := models.GetProductQueryParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	products, err := models.GetProductList(limit, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Mengkonversi hasil query menjadi format JSON
	jsonBytes, err := json.Marshal(products)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set header Content-Type sebagai application/json
	w.Header().Set("Content-Type", "application/json")

	// Mengembalikan hasil query sebagai response
	w.Write(jsonBytes)
}

// Get a product by ID
func getProductHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the product ID from the URL parameter
	vars := mux.Vars(r)
	productID := vars["id"]

	// Retrieve the product
	product, err := models.GetProduct(productID)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "Product not found")
		return
	}

	// Set response status code
	w.WriteHeader(http.StatusOK)
}

// Update a product by ID
func updateProductHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the product ID from the URL parameter
	vars := mux.Vars(r)
	productID := vars["id"]

	// Validate user input
	var updatedProduct models.Product
	err := json.NewDecoder(r.Body).Decode(&updatedProduct)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err = models.ValidateProduct(&updatedProduct); err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Update the product
	if err := models.UpdateProduct(productID, &updatedProduct); err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}


	// Set response status code
	w.WriteHeader(http.StatusOK)

