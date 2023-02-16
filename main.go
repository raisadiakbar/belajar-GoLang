// package main

// import (
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"time"
// )

// func main() {

// 	// Start the server
// 	fmt.Println("Server is listening on port 8080...")
// 	srv := &http.Server{
// 		Handler:      r,
// 		Addr:         "127.0.0.1:8080",
// 		WriteTimeout: 15 * time.Second,
// 		ReadTimeout:  15 * time.Second,
// 	}
// 	log.Fatal(srv.ListenAndServe())
// }

// // Register handler
// func registerHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintln(w, "Register")
// }

// // Login handler
// func loginHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintln(w, "Login")
// }
