package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"golang.org/x/crypto/bcrypt"
)

var DB *gorm.DB

func connectDB() (*gorm.DB, error) {
	dbConn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DATABASE"),
	)
	db, err := gorm.Open(os.Getenv("DB_CONNECTION"), dbConn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func CloseDB(db *gorm.DB) {
	err := db.Close()
	if err != nil {
		log.Fatal("Error closing database connection")
	}
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
	log.Fatal(http.ListenAndServe(":8888", r))
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
	Status          string    `json:"status"`
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
	// Dekode request body ke dalam objek model `Address`
	var address Address
	if err := json.NewDecoder(r.Body).Decode(&address); err != nil {
		// Jika terjadi masalah saat dekode request body, kirim pesan kesalahan dengan status 400 Bad Request
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Gagal memproses request body: %v", err)
		return
	}

	// Simpan alamat baru ke dalam database
	if err := DB.Create(&address).Error; err != nil {
		// Jika terjadi masalah saat menyimpan data, kirim pesan kesalahan dengan status 500 Internal Server Error
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Gagal menyimpan alamat baru: %v", err)
		return
	}

	// Jika penyimpanan berhasil, kirim response dengan status 201 Created dan data alamat yang baru saja ditambahkan
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(address)
}

func getAddressHandler(w http.ResponseWriter, r *http.Request) {
	// Dapatkan id dari URL parameter
	vars := mux.Vars(r)
	id := vars["id"]

	// Inisialisasi objek model `Address`
	var address Address

	// Cari alamat dengan id yang diberikan dari database
	if err := DB.First(&address, id).Error; err != nil {
		// Jika alamat tidak ditemukan, kirim pesan kesalahan dengan status 404 Not Found
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Alamat dengan id %s tidak ditemukan", id)
		return
	}

	// Jika alamat ditemukan, kirim response dengan data alamat yang ditemukan
	json.NewEncoder(w).Encode(address)
}

// Get user ID from JWT
func getUserIdFromToken(w http.ResponseWriter, r *http.Request) int {
	// Parse JWT from Authorization header
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
		return 0
	}
	tokenString = strings.ReplaceAll(tokenString, "Bearer ", "")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("my-secret-key"), nil
	})
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return 0
	}

	// Get user ID from token claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return 0
	}
	userID, ok := claims["user_id"].(float64)
	if !ok {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return 0
	}

	return int(userID)
}

type contextKey string

const (
	authContextKey contextKey = "auth"
)

func updateAddressHandler(w http.ResponseWriter, r *http.Request) {
	// get user ID from JWT token
	claims, ok := r.Context().Value(authContextKey).(jwt.MapClaims)
	if !ok {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}
	userID, ok := claims["user_id"].(float64)
	if !ok {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	// get address ID from URL path
	vars := mux.Vars(r)
	addressID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "invalid address ID", http.StatusBadRequest)
		return
	}

	// parse JSON request body into Address struct
	var address Address
	err = json.NewDecoder(r.Body).Decode(&address)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// validate user ID
	if uint(userID) != address.UserID {
		http.Error(w, "you are not authorized to update this address", http.StatusUnauthorized)
		return
	}

	// update address in database
	db, err := connectDB()
	defer db.Close()

	var existingAddress Address
	err = db.Where("id = ?", addressID).First(&existingAddress).Error
	if err != nil {
		http.Error(w, "address not found", http.StatusNotFound)
		return
	}
	existingAddress.Name = address.Name
	existingAddress.Street = address.Street
	existingAddress.City = address.City
	existingAddress.Province = address.Province
	existingAddress.Zipcode = address.Zipcode
	err = db.Save(&existingAddress).Error
	if err != nil {
		http.Error(w, "failed to update address", http.StatusInternalServerError)
		return
	}

	// return updated address in response
	json.NewEncoder(w).Encode(existingAddress)
}

func deleteAddressHandler(w http.ResponseWriter, r *http.Request) {
	// Dapatkan id dari URL parameter
	vars := mux.Vars(r)
	id := vars["id"]

	// Inisialisasi objek model `Address`
	var address Address

	// Cari alamat dengan id yang diberikan dari database
	if err := DB.First(&address, id).Error; err != nil {
		// Jika alamat tidak ditemukan, kirim pesan kesalahan dengan status 404 Not Found
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Alamat dengan id %s tidak ditemukan", id)
		return
	}

	// Hapus alamat dari database
	if err := DB.Delete(&address).Error; err != nil {
		// Jika terjadi masalah saat menghapus, kirim pesan kesalahan dengan status 500 Internal Server Error
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Gagal menghapus alamat dengan id %s: %v", id, err)
		return
	}

	// Send a success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Address deleted"})
}

func createCategoryHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request body to Category struct
	var category Category
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Save Category to database using ORM
	result := DB.Create(&category)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Return JSON response with created Category object
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

func getCategoryListHandler(w http.ResponseWriter, r *http.Request) {
	// Query all Category objects from database using ORM
	var categories []Category
	result := DB.Find(&categories)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Return JSON response with list of Category objects
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

func getCategoryHandler(w http.ResponseWriter, r *http.Request) {
	// Get category ID from URL path parameter
	vars := mux.Vars(r)
	categoryID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Query Category object from database using ORM
	var category Category
	result := DB.First(&category, categoryID)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusNotFound)
		return
	}

	// Return JSON response with Category object
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

func updateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	// Get category ID from URL path parameter
	vars := mux.Vars(r)
	categoryID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Parse request body to Category struct
	var category Category
	err = json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update Category object in database using ORM
	result := DB.Model(&Category{}).Where("id = ?", categoryID).Updates(category)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Return JSON response with updated Category object
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

func deleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	// ambil id kategori dari path parameter
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// hapus kategori dari database
	category := Category{ID: uint(id)}
	if err := DB.Delete(&category).Error; err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// kirim status sukses ke client
	w.WriteHeader(http.StatusOK)
}

func createProductHandler(w http.ResponseWriter, r *http.Request) {
	// Ambil data dari request body
	var product Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Lakukan validasi data produk
	if product.Name == "" {
		http.Error(w, "Product name is required", http.StatusBadRequest)
		return
	}

	// Simpan data produk ke database
	err = DB.Create(&product).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Kirim response dengan data produk yang baru saja dibuat
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getProductListHandler(w http.ResponseWriter, r *http.Request) {
	// Ambil data produk dari database
	var products []Product
	err := DB.Find(&products).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Kirim response dengan data produk
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(products)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getProductHandler(w http.ResponseWriter, r *http.Request) {
	// Ambil ID produk dari URL parameter
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Ambil data produk dari database
	var product Product
	err = DB.First(&product, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Product not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Kirim response dengan data produk
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func updateProductHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]

	// Check if product exists
	var product Product
	result := DB.First(&product, productID)
	if result.Error != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Product not found")
		return
	}

	// Decode request body into Product struct
	var updatedProduct Product
	err := json.NewDecoder(r.Body).Decode(&updatedProduct)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request payload")
		return
	}

	// Update product fields
	product.Name = updatedProduct.Name
	product.Description = updatedProduct.Description
	product.Price = updatedProduct.Price
	product.Image = updatedProduct.Image
	product.Stock = updatedProduct.Stock
	product.UpdatedAt = time.Now()

	// Save changes to database
	DB.Save(&product)

	// Return updated product as JSON
	json.NewEncoder(w).Encode(&product)
}

func deleteProductHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]

	// Check if product exists
	var product Product
	result := DB.First(&product, productID)
	if result.Error != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Product not found")
		return
	}

	// Delete product from database
	DB.Delete(&product)

	// Return success message
	fmt.Fprintf(w, "Product deleted")
}

func createTransactionHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request body to Transaction struct
	var transaction Transaction
	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Insert transaction to database
	err = DB.Create(&transaction).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(transaction)
}

func getTransactionListHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve all transactions from database
	var transactions []Transaction
	err := DB.Find(&transactions).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return transactions as response
	json.NewEncoder(w).Encode(transactions)
}

func getTransactionHandler(w http.ResponseWriter, r *http.Request) {
	// Mendapatkan nilai id dari path parameter
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Mencari transaksi dengan id yang sesuai dari database
	var transaction Transaction
	err = DB.First(&transaction, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Mengembalikan response dengan data transaksi yang ditemukan
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transaction)
}

func confirmTransactionHandler(w http.ResponseWriter, r *http.Request) {
	// Mendapatkan nilai id dari path parameter
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Mencari transaksi dengan id yang sesuai dari database
	var transaction Transaction
	err = DB.First(&transaction, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Mengubah status transaksi menjadi "confirmed"
	transaction.Status = "confirmed"
	err = DB.Save(&transaction).Error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Mengembalikan response dengan data transaksi yang telah diubah statusnya
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transaction)
}
