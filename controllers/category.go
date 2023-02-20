package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	models "e-GoLang/models"
)

// CreateCategoryRequest struct
type CreateCategoryRequest struct {
	Name string `json:"name"`
}

// UpdateCategoryRequest struct
type UpdateCategoryRequest struct {
	Name string `json:"name"`
}

// createCategoryHandler is used to create a new category
func createCategoryHandler(w http.ResponseWriter, r *http.Request) {
	// Get the body from the request
	var req CreateCategoryRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a new category
	category := models.Category{
		Name: req.Name,
	}
	result := models.DB.Create(&category)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Return the created category
	json.NewEncoder(w).Encode(category)
}

// getCategoryListHandler is used to get a list of categories
func getCategoryListHandler(w http.ResponseWriter, r *http.Request) {
	// Get the page and limit from the query string
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 10
	}

	// Get a list of categories with pagination
	var categories []models.Category
	result := models.DB.Offset((page - 1) * limit).Limit(limit).Find(&categories)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Get the total count of categories
	var count int64
	result = models.DB.Model(&models.Category{}).Count(&count)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Return the list of categories
	response := map[string]interface{}{
		"data":         categories,
		"total_data":   count,
		"current_page": page,
		"per_page":     limit,
	}
	json.NewEncoder(w).Encode(response)
}

// getCategoryHandler is used to get a single category
func getCategoryHandler(w http.ResponseWriter, r *http.Request) {
	// Get the ID parameter from the URL
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the category with the specified ID
	var category models.Category
	result := models.DB.First(&category, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			http.Error(w, "Category not found", http.StatusNotFound)
			return
		}
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Return the category
	json.NewEncoder(w).Encode(category)
}

func updateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// get the id from the request params
	params := mux.Vars(r)
	idStr := params["ID"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	// decode the request body into a Category model
	var category models.Category
	err = json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// update the category with the given id
	err = models.UpdateCategory(id, &category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// encode the updated category into a response and return it
	json.NewEncoder(w).Encode(category)
}

func deleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// get the id from the request params
	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	// delete the category with the given id
	err = models.DeleteCategory(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// return a success message
	json.NewEncoder(w).Encode(map[string]string{"message": "Category deleted"})
}

// func route() {
// 	r := mux.NewRouter()

// 	// Category routes
// 	r.HandleFunc("/api/categories", createCategoryHandler).Methods("POST")
// 	r.HandleFunc("/api/categories", getCategoryListHandler).Methods("GET")
// 	r.HandleFunc("/api/categories/{id}", getCategoryHandler).Methods("GET")
// 	r.HandleFunc("/api/categories/{id}", updateCategoryHandler).Methods("PUT")
// 	r.HandleFunc("/api/categories/{id}", deleteCategoryHandler).Methods("DELETE")

// 	// Serve the API
// 	http.ListenAndServe(":8000", r)
// }
