package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	models "e-GoLang/models"

	"github.com/gorilla/mux"
)

// createTransactionHandler is used to create a new transaction
func createTransactionHandler(w http.ResponseWriter, r *http.Request) {
	// Parse JSON request body
	decoder := json.NewDecoder(r.Body)
	var transaction models.Transaction
	err := decoder.Decode(&transaction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Set transaction time
	transaction.TransactionTime = time.Now()

	// Save transaction to database
	err = transaction.Save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set response status code
	w.WriteHeader(http.StatusCreated)
}

// getTransactionListHandler is used to get a list of all transactions
func getTransactionListHandler(w http.ResponseWriter, r *http.Request) {
	// Get transactions from database
	transactions, err := models.GetAllTransactions()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert transactions to JSON
	transactionsJSON, err := json.Marshal(transactions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set response headers and write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(transactionsJSON)
}

// getTransactionHandler is used to get a single transaction by ID
func getTransactionHandler(w http.ResponseWriter, r *http.Request) {
	// Get transaction ID from request URL
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
		return
	}

	// Get transaction from database
	transaction, err := models.GetTransactionByID(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if transaction == nil {
		http.Error(w, "Transaction not found", http.StatusNotFound)
		return
	}

	// Convert transaction to JSON
	transactionJSON, err := json.Marshal(transaction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set response headers and write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(transactionJSON)
}

// confirmTransactionHandler is used to confirm a transaction by ID
func confirmTransactionHandler(w http.ResponseWriter, r *http.Request) {
	// Get transaction ID from request URL
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
		return
	}

	// Get transaction from database
	transaction, err := models.GetTransactionByID(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if transaction == nil {
		http.Error(w, "Transaction not found", http.StatusNotFound)
		return
	}

	// Set transaction status to confirmed
	transaction.Status = "confirmed"

	// Update transaction in database
	err = transaction.Update()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set response status code
	w.WriteHeader(http.StatusOK)
}
