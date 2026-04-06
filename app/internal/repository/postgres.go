package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jasonsavage/financeflow/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type postgresDBRepo struct {
	DB *pgxpool.Pool
}

func NewPostgresRepo(conn *pgxpool.Pool) DatabaseRepo {
	return &postgresDBRepo{
		DB: conn,
	}
}

func (r *postgresDBRepo) RegisterUser(username, password string) (*models.User, error) {
	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	userID := fmt.Sprintf("USR-%d", time.Now().UnixNano())

	// Insert user into DB
	var u models.User
	now := time.Now()
	err = r.DB.QueryRow(context.Background(), `
		INSERT INTO users (id, username, email, password_hash, first_name, last_name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, username, email, created_at, updated_at
	`, userID, username, username, string(hash), "", "", now, now).Scan(&u.ID, &u.Username, &u.Email, &u.CreatedAt, &u.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to register user: %w", err)
	}

	return &u, nil
}

func (r *postgresDBRepo) AuthenticateUser(username, password string) (*models.User, error) {
	// Look up user by username
	var u models.User
	err := r.DB.QueryRow(context.Background(), `
		SELECT id, username, email, password_hash, created_at, updated_at FROM users WHERE username = $1
	`, username).Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get user auth: %w", err)
	}

	// Compare password to hash
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid password: %w", err)
	}

	return &u, nil
}

func (r *postgresDBRepo) SaveTransactions(ctx context.Context, transactions []models.Transaction) (int, error) {
	if len(transactions) == 0 {
		return 0, nil
	}

	var inserted int
	for _, t := range transactions {
		res, err := r.DB.Exec(ctx, `
			INSERT INTO transactions (
				id, user_id, date, description, category, 
				deposit, withdrawal, bank_account_name, created_at, updated_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			ON CONFLICT (id) DO NOTHING
		`, t.ID, t.UserID, t.Date, t.Description, t.Category,
			t.Deposit, t.Withdrawal, t.BankAccountName, t.CreatedAt, t.UpdatedAt)

		if err != nil {
			return inserted, fmt.Errorf("failed to insert transaction: %w", err)
		}
		if res.RowsAffected() > 0 {
			inserted++
		}
	}

	return inserted, nil
}

func (r *postgresDBRepo) GetTransactions(ctx context.Context, userID string, from, to *time.Time) ([]models.Transaction, error) {
	query := `SELECT id, user_id, date, description, category, deposit, withdrawal, bank_account_name, created_at, updated_at FROM transactions WHERE user_id = $1`
	args := []interface{}{userID}
	
	if from != nil {
		args = append(args, *from)
		query += fmt.Sprintf(" AND date >= $%d", len(args))
	}
	if to != nil {
		args = append(args, *to)
		query += fmt.Sprintf(" AND date <= $%d", len(args))
	}

	query += " ORDER BY date DESC"

	rows, err := r.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query transactions: %w", err)
	}
	defer rows.Close()

	// Initialize as empty slice rather than nil so it serializes as [] in JSON when empty
	transactions := make([]models.Transaction, 0)
	for rows.Next() {
		var t models.Transaction
		err := rows.Scan(
			&t.ID, &t.UserID, &t.Date, &t.Description, 
			&t.Category, &t.Deposit, &t.Withdrawal, &t.BankAccountName, 
			&t.CreatedAt, &t.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, t)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return transactions, nil
}

func (r *postgresDBRepo) Ping(ctx context.Context) error {
	return r.DB.Ping(ctx)
}

func (r *postgresDBRepo) GetUserDetails(ctx context.Context, userID string) (*models.User, error) {
	var u models.User
	err := r.DB.QueryRow(ctx, `
		SELECT id, username, email, first_name, last_name, password_hash, created_at, updated_at 
		FROM users WHERE id = $1
	`, userID).Scan(&u.ID, &u.Username, &u.Email, &u.FirstName, &u.LastName, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get user details: %w", err)
	}
	return &u, nil
}

func (r *postgresDBRepo) UpdateUserDetails(ctx context.Context, userID string, firstName, lastName, email string) error {
	now := time.Now()
	_, err := r.DB.Exec(ctx, `
		UPDATE users SET first_name = $1, last_name = $2, email = $3, updated_at = $4
		WHERE id = $5
	`, firstName, lastName, email, now, userID)
	if err != nil {
		return fmt.Errorf("failed to update user details: %w", err)
	}
	return nil
}
