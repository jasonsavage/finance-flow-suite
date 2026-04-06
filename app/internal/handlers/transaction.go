package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jasonsavage/financeflow/internal/helpers"
	"github.com/jasonsavage/financeflow/internal/middleware"
	"github.com/jasonsavage/financeflow/internal/repository"
)

type TransactionHandler struct {
	Repo repository.DatabaseRepo
}

func NewTransactionHandler(repo repository.DatabaseRepo) *TransactionHandler {
	return &TransactionHandler{Repo: repo}
}

// POST /transactions/upload
func (h *TransactionHandler) UploadTransactions(w http.ResponseWriter, r *http.Request) {
	// Parse 10MB max
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	bankAccountName := r.URL.Query().Get("bankAccount")
	if bankAccountName == "" {
		http.Error(w, "bankAccount query parameter is required", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file is required in multipart form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	transactions, err := helpers.ParseTransactionsCSV(file, userID, bankAccountName)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse csv: %v", err), http.StatusBadRequest)
		return
	}

	count, err := h.Repo.SaveTransactions(r.Context(), transactions)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to save transactions: %v", err), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Upload successful",
		"count":   count,
	})
}

// GET /transactions/list
func (h *TransactionHandler) ListTransactions(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var fromTime, toTime *time.Time

	fromStr := r.URL.Query().Get("from")
	if fromStr != "" {
		t, err := time.Parse("2006-01-02", fromStr)
		if err == nil {
			fromTime = &t
		}
	}

	toStr := r.URL.Query().Get("to")
	if toStr != "" {
		t, err := time.Parse("2006-01-02", toStr)
		if err == nil {
			toTime = &t
		}
	}

	transactions, err := h.Repo.GetTransactions(r.Context(), userID, fromTime, toTime)
	if err != nil {
		http.Error(w, "Failed to load transactions", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, transactions)
}
