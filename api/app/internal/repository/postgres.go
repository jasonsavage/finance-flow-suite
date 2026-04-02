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

	accountID := fmt.Sprintf("ACC-%d", time.Now().UnixNano())

	// Insert user into DB
	var u models.User
	now := time.Now()
	err = r.DB.QueryRow(context.Background(), `
		INSERT INTO users (account_id, username, email, password_hash, first_name, last_name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, account_id, username, email, created_at, updated_at
	`, accountID, username, username, string(hash), "", "", now, now).Scan(&u.ID, &u.AccountID, &u.Username, &u.Email, &u.CreatedAt, &u.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to register user: %w", err)
	}

	return &u, nil
}

func (r *postgresDBRepo) AuthenticateUser(username, password string) (*models.User, error) {
	// Look up user by username
	var u models.User
	err := r.DB.QueryRow(context.Background(), `
		SELECT id, account_id, username, email, password_hash, created_at, updated_at FROM users WHERE username = $1
	`, username).Scan(&u.ID, &u.AccountID, &u.Username, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get user auth: %w", err)
	}

	// Compare password to hash
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid password: %w", err)
	}

	return &u, nil
}
