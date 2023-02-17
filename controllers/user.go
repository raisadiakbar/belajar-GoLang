package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type User struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email" gorm:"unique"`
	Password  string    `json:"-"`
	Phone     string    `json:"phone" gorm:"unique"`
	Address   string    `json:"address,omitempty" gorm:"foreignkey:UserID"`
	Store     string    `json:"store,omitempty" gorm:"foreignkey:UserID"`
	Role      string    `json:"role" gorm:"default:'user'"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type contextKey string

const userIDContextKey contextKey = "userID"

// getAccountHandler is the handler for GET /api/accounts/me.
func getAccountHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve user ID from context.
	userID, err := getUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Retrieve user from database.
	user, err := getUserByID(userID)
	if err != nil {
		http.Error(w, "failed to retrieve user", http.StatusInternalServerError)
		return
	}

	// Return user data as JSON response.
	json.NewEncoder(w).Encode(user)
}

// updateAccountHandler is the handler for PUT /api/accounts/me.
func updateAccountHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve user ID from context.
	userID, err := getUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Retrieve updated user data from request body.
	var updatedUser User
	err = json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		http.Error(w, "failed to decode request body", http.StatusBadRequest)
		return
	}

	// // Ensure that user ID in request body matches user ID in context.
	// if updatedUser.ID != userID {
	// 	http.Error(w, "invalid user ID", http.StatusBadRequest)
	// 	return
	// }

	// Retrieve existing user from database.
	existingUser, err := getUserByID(userID)
	if err != nil {
		http.Error(w, "failed to retrieve user", http.StatusInternalServerError)
		return
	}

	// Merge updated user data with existing user data.
	existingUser.Name = updatedUser.Name
	existingUser.Email = updatedUser.Email
	existingUser.Phone = updatedUser.Phone

	// Update user in database.
	err = updateUser(existingUser)
	if err != nil {
		http.Error(w, "failed to update user", http.StatusInternalServerError)
		return
	}

	// Return updated user data as JSON response.
	json.NewEncoder(w).Encode(existingUser)
}

// getUserIDFromContext retrieves the user ID from the context.
func getUserIDFromContext(ctx context.Context) (int, error) {
	// Retrieve user ID from context.
	userID, ok := ctx.Value(userIDContextKey).(int)
	if !ok {
		return 0, errors.New("user ID not found in context")
	}

	return userID, nil
}

// getUserByID retrieves a user from the database by ID.
func getUserByID(userID int) (*User, error) {
	// TODO: implement retrieval of user from database.
	return nil, errors.New("not implemented")
}

// updateUser updates a user in the database.
func updateUser(user *User) error {
	// TODO: implement update of user in database.
	return errors.New("not implemented")
}
