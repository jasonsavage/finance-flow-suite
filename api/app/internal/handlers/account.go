package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jasonsavage/financeflow/internal/models"
)

type AccountHandler struct {
	DB *pgxpool.Pool
}

func NewAccountHandler(db *pgxpool.Pool) *AccountHandler {
	return &AccountHandler{DB: db}
}

// GET /accounts
func (h *AccountHandler) List(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query(context.Background(), `
		SELECT id, name, type, balance, currency, created_at, updated_at FROM accounts
	`)
	if err != nil {
		http.Error(w, "Failed to fetch accounts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	accounts := []models.Account{}
	for rows.Next() {
		var a models.Account
		if err := rows.Scan(&a.ID, &a.Name, &a.Type, &a.Balance, &a.Currency, &a.CreatedAt, &a.UpdatedAt); err != nil {
			http.Error(w, "Failed to scan account", http.StatusInternalServerError)
			return
		}
		accounts = append(accounts, a)
	}

	writeJSON(w, http.StatusOK, accounts)
}

// GET /accounts/{id}
func (h *AccountHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid account ID", http.StatusBadRequest)
		return
	}

	var a models.Account
	err = h.DB.QueryRow(context.Background(), `
		SELECT id, name, type, balance, currency, created_at, updated_at FROM accounts WHERE id = $1
	`, id).Scan(&a.ID, &a.Name, &a.Type, &a.Balance, &a.Currency, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, a)
}

// POST /accounts
func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	var a models.Account
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	now := time.Now()
	err := h.DB.QueryRow(context.Background(), `
		INSERT INTO accounts (name, type, balance, currency, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`, a.Name, a.Type, a.Balance, a.Currency, now, now).Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		http.Error(w, "Failed to create account", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, a)
}

// PUT /accounts/{id}
func (h *AccountHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid account ID", http.StatusBadRequest)
		return
	}

	var a models.Account
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.DB.QueryRow(context.Background(), `
		UPDATE accounts SET name=$1, type=$2, balance=$3, currency=$4, updated_at=$5
		WHERE id=$6
		RETURNING id, name, type, balance, currency, created_at, updated_at
	`, a.Name, a.Type, a.Balance, a.Currency, time.Now(), id).Scan(
		&a.ID, &a.Name, &a.Type, &a.Balance, &a.Currency, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		http.Error(w, "Failed to update account", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, a)
}

// DELETE /accounts/{id}
func (h *AccountHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid account ID", http.StatusBadRequest)
		return
	}

	result, err := h.DB.Exec(context.Background(), `DELETE FROM accounts WHERE id = $1`, id)
	if err != nil || result.RowsAffected() == 0 {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
