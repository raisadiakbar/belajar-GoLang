package controllers

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"

	"github.com/raisadiakbar/belajar-GoLang/models"
	"github.com/raisadiakbar/belajar-GoLang/repository"
)

type AuthController struct{}

var authRepository = repository.AuthRepository{}

// Register a new user
func (ac AuthController) Register(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err = user.Validate(); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// check if email already exists
	if authRepository.EmailExists(user.Email) {
		respondWithError(w, http.StatusBadRequest, "Email already exists")
		return
	}

	// check if phone number already exists
	if authRepository.PhoneExists(user.Phone) {
		respondWithError(w, http.StatusBadRequest, "Phone number already exists")
		return
	}

	// hash the password before storing in the database
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	if err = authRepository.Create(&user); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// create a new store for the user
	store := models.Store{UserID: user.ID}
	if err = repository.StoreRepository{}.Create(&store); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]string{"message": "User created successfully"})
}

// Login user
func (ac AuthController) Login(w http.ResponseWriter, r *http.Request) {
	loginRequest := models.LoginRequest{}
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// get user with given email
	user, err := authRepository.FindByEmail(loginRequest.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// check if password is correct
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
	})
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error generating token")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"token": tokenString})
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	ac := AuthController{}
	ac.Register(w, r)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	ac := AuthController{}
	ac.Login(w, r)
}

// Login and register routes
func AuthRoutes(r *mux.Router) {
	r.HandleFunc("/api/auth/register", registerHandler).Methods("POST")
	r.HandleFunc("/api/auth/login", loginHandler).Methods("POST")
}
