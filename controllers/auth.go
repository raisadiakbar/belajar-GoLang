package controllers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct{}

var DB *gorm.DB

func connectDB() (*gorm.DB, error) {
	db, err := gorm.Open("mysql", "root:Password@(localhost)/project-golang?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		return nil, err
	}

	return db, nil
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
