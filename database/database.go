package database

import (
	"go-jwt/database/models"

	"github.com/jinzhu/gorm"
)

// Storager - hold logic for work with DB
type Storager interface {
	Read
	Write
}

// Read represents getting methods
type Read interface {
	GetUserByEmail(email string) (models.User, error)
	GetUsers() ([]models.User, error)
}

// Write represents inserting methods
type Write interface {
	Register(user models.User) (*gorm.DB, error)
}
