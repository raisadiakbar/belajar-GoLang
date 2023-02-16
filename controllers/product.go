package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/example/app/models"
	"github.com/example/app/utils"
)

// Create a new product
func createProductHandler(w http.ResponseWriter, r *http.Request) {
	// Validate user input
	var newProduct models.Product
	err := json.NewDecoder(r.Body).Decode(&newProduct)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err = newProduct.Validate(); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Create the new product
	if err := models.CreateProduct(&newProduct); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Send the response
	utils.RespondWithJSON(w, http.StatusCreated, newProduct)
}

// Get a list of all products
func getProductListHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve query parameters for pagination and filtering
	limit, offset, nameFilter := utils.GetProductQueryParams(r)

	// Retrieve the products
	products, err := models.GetProductList(limit, offset, nameFilter)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Send the response
	utils.RespondWithJSON(w, http.StatusOK, products)
}

// Get a product by ID
func getProductHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the product ID from the URL parameter
	vars := mux.Vars(r)
	productID := vars["id"]

	// Retrieve the product
	product, err := models.GetProduct(productID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Product not found")
		return
	}

	// Send the response
	utils.RespondWithJSON(w, http.StatusOK, product)
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
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err = updatedProduct.Validate(); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Update the product
	if err := models.UpdateProduct(productID, &updatedProduct); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Send the response
	utils.RespondWithJSON(w, http.StatusOK, updatedProduct)
}

// Delete a product by ID
func deleteProductHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the product ID from the URL parameter
	vars := mux.Vars(r)
	productID := vars["id"]

	// Delete the product
	if err := models.DeleteProduct(productID); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Send the response
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
