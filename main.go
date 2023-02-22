package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"golang.org/x/crypto/bcrypt"
)

var DB *gorm.DB

func connectDB() (*gorm.DB, error) {
	db, err := gorm.Open("mysql", "root:Password@(localhost)/project-golang?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		return nil, err
	}

	return db, nil
}

func CloseDB() {
	DB.Close()
	fmt.Println("Successfully closed database connection")
}

func main() {
	// Membuat koneksi ke database
	db, err := connectDB()

	if err != nil {
		fmt.Println("Failed to connect to database:", err)
	} else {
		fmt.Println("Successfully connected to database")
	}

	defer db.Close()

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
	log.Fatal(http.ListenAndServe(":8080", r))
}

type User struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email" gorm:"unique"`
	Password  string    `json:"-"`
	Phone     string    `json:"phone" gorm:"unique"`
	Address   []Address `json:"address,omitempty" gorm:"foreignkey:UserID"`
	Store     Store     `json:"store,omitempty" gorm:"foreignkey:UserID"`
	Role      string    `json:"role" gorm:"default:'user'"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Address struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	UserID    uint      `json:"user_id"`
	Name      string    `json:"name"`
	Street    string    `json:"street"`
	City      string    `json:"city"`
	Province  string    `json:"province"`
	Zipcode   string    `json:"zipcode"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Store struct {
	ID          uint      `gorm:"primary_key" json:"id"`
	UserID      uint      `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Category struct {
	ID          uint      `gorm:"primary_key" json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsAdmin     bool      `json:"is_admin"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Product struct {
	ID          uint      `gorm:"primary_key" json:"id"`
	UserID      uint      `json:"user_id"`
	CategoryID  uint      `json:"category_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       uint      `json:"price"`
	Image       string    `json:"image"`
	Stock       uint      `json:"stock"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Transaction struct {
	ID              uint      `gorm:"primary_key" json:"id"`
	UserID          uint      `json:"user_id"`
	ProductID       uint      `json:"product_id"`
	Quantity        uint      `json:"quantity"`
	TotalPrice      uint      `json:"total_price"`
	AddressID       uint      `json:"address_id"`
	TransactionTime time.Time `json:"transaction_time"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type LogProduct struct {
	ID            uint      `gorm:"primary_key" json:"id"`
	TransactionID uint      `json:"transaction_id"`
	ProductID     uint      `json:"product_id"`
	Quantity      uint      `json:"quantity"`
	Price         uint      `json:"price"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func isEmailExist(email string) bool {
	db, err := connectDB()
	defer db.Close()

	var count int
	err = db.Raw("SELECT count(*) FROM users WHERE email=?", email).Scan(&count).Error
	if err != nil {
		log.Fatal(err)
	}
	return count > 0
}

func isPhoneExist(phone string) bool {
	db, err := connectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var count int
	row := db.Raw("SELECT count(*) FROM users WHERE phone=?", phone).Row()
	row.Scan(&count)

	return count > 0
}

func createUser(user User) error {
	db, err := connectDB()
	defer db.Close()

	if err != nil {
		return err
	}

	if isEmailExist(user.Email) {
		return errors.New("Email already exists")
	}

	if isPhoneExist(user.Phone) {
		return errors.New("Phone already exists")
	}

	if result := db.Exec("INSERT INTO users(name, email, phone, password) VALUES (?, ?, ?, ?)", user.Name, user.Email, user.Phone, user.Password); result.Error != nil {
		return result.Error
	}

	return nil
}

func getUserByEmail(email string) (*User, error) {
	db, err := connectDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var user User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func generateToken(userID int64) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	// ambil data yang diberikan oleh user pada body request
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// cek apakah email atau no telepon sudah terdaftar
	if isEmailExist(user.Email) {
		http.Error(w, "Email already exist", http.StatusBadRequest)
		return
	}
	if isPhoneExist(user.Phone) {
		http.Error(w, "Phone already exist", http.StatusBadRequest)
		return
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	// simpan data user ke database
	err = createUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// kirim response ke user
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "User created"})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// ambil data yang diberikan oleh user pada body request
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// cek apakah email atau no telepon terdaftar
	userData, err := getUserByEmail(user.Email)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusBadRequest)
		return
	}

	// bandingkan password yang diberikan oleh user dengan password yang tersimpan di database
	err = bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(user.Password))
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusBadRequest)
		return
	}

	// generate token JWT
	token, err := generateToken(int64(userData.ID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// kirim response ke user
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func getAccountHandler(w http.ResponseWriter, r *http.Request) {
	// Mendapatkan user dari token JWT
	user := r.Context().Value("user").(*User)

	// Menampilkan data user
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func updateAccountHandler(w http.ResponseWriter, r *http.Request) {
	// Mendapatkan user dari token JWT
	user := r.Context().Value("user").(*User)

	// Mengambil data yang diberikan oleh user pada body request
	var updatedUser User
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

func createAddressHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var address models.Address
	if err := json.NewDecoder(r.Body).Decode(&address); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate request body
	if err := validate.Struct(address); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Get user ID from token
	userID, err := getUserIDFromToken(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	// Check if address already exists for user
	if _, err := services.GetAddressByUserIDAndLabel(userID, address.Label); err == nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Address with label '%s' already exists", address.Label))
		return
	}

	// Create address
	if err := services.CreateAddress(userID, &address); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create address")
		return
	}

	// Respond with created address
	respondWithJSON(w, http.StatusCreated, address)
}

func getAddressHandler(w http.ResponseWriter, r *http.Request) {
	// Get user ID from token
	userID, err := getUserIDFromToken(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	// Get address ID from URL parameter
	vars := mux.Vars(r)
	addressID := vars["id"]

	// Get address by ID
	address, err := services.GetAddressByUserIDAndID(userID, addressID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			respondWithError(w, http.StatusNotFound, "Address not found")
		} else {
			respondWithError(w, http.StatusInternalServerError, "Failed to get address")
		}
		return
	}

	// Respond with address
	respondWithJSON(w, http.StatusOK, address)
}
