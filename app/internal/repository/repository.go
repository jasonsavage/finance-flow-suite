package repository

import (
	"context"
	"time"

	"github.com/jasonsavage/financeflow/internal/models"
)

type DatabaseRepo interface {
	RegisterUser(username, password string) (*models.User, error)
	AuthenticateUser(username, password string) (*models.User, error)
	SaveTransactions(ctx context.Context, transactions []models.Transaction) (int, error)
	GetTransactions(ctx context.Context, userID string, from, to *time.Time) ([]models.Transaction, error)
	GetUserDetails(ctx context.Context, userID string) (*models.User, error)
	UpdateUserDetails(ctx context.Context, userID string, firstName, lastName, email string) error
	Ping(ctx context.Context) error
}
