package internal

import (
	"go-jwt/database/models"
)

// Handler - represents main business logic
type Handler interface {
	Register(user models.User) (interface{}, error)
	Login(email, password string) (interface{}, error)
	GetUsers() (interface{}, error)
}
