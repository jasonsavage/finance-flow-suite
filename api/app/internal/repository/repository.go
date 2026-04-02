package repository

import "github.com/jasonsavage/financeflow/internal/models"

type DatabaseRepo interface {
	RegisterUser(username, password string) (*models.User, error)
	AuthenticateUser(username, password string) (*models.User, error)
}
