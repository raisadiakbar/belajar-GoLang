package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// Login and register routes
	r.HandleFunc("/api/auth/register", registerHandler).Methods("POST")
	r.HandleFunc("/api/auth/login", loginHandler).Methods("POST")

	// Account routes
	r.HandleFunc("/api/accounts/me", getAccountHandler).Methods("GET")
	r.HandleFunc("/api/accounts/me", updateAccountHandler).Methods("PUT")

	// Address routes
	r.HandleFunc("/api/addresses", createAddressHandler).Methods("POST")
	r.HandleFunc("/api/addresses/{id}", getAddressHandler).Methods("GET")
	r.HandleFunc("/api/addresses/{id}", updateAddressHandler).Methods("PUT")
	r.HandleFunc("/api/addresses/{id}", deleteAddressHandler).Methods("DELETE")

	// Category routes
	r.HandleFunc("/api/categories", createCategoryHandler).Methods("POST")
	r.HandleFunc("/api/categories", getCategoryListHandler).Methods("GET")
	r.HandleFunc("/api/categories/{id}", getCategoryHandler).Methods("GET")
	r.HandleFunc("/api/categories/{id}", updateCategoryHandler).Methods("PUT")
	r.HandleFunc("/api/categories/{id}", deleteCategoryHandler).Methods("DELETE")

	// Product routes
	r.HandleFunc("/api/products", createProductHandler).Methods("POST")
	r.HandleFunc("/api/products", getProductListHandler).Methods("GET")
	r.HandleFunc("/api/products/{id}", getProductHandler).Methods("GET")
	r.HandleFunc("/api/products/{id}", updateProductHandler).Methods("PUT")
	r.HandleFunc("/api/products/{id}", deleteProductHandler).Methods("DELETE")

	// Transaction routes
	r.HandleFunc("/api/transactions", createTransactionHandler).Methods("POST")
	r.HandleFunc("/api/transactions", getTransactionListHandler).Methods("GET")
	r.HandleFunc("/api/transactions/{id}", getTransactionHandler).Methods("GET")
	r.HandleFunc("/api/transactions/{id}/confirm", confirmTransactionHandler).Methods("POST")

	// Serve the API
	http.ListenAndServe(":8000", r)
}
