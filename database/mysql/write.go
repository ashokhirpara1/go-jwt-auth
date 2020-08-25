package mysql

import (
	"go-jwt/database/models"

	"github.com/jinzhu/gorm"
)

// Register - create new user
func (db *DB) Register(user models.User) (*gorm.DB, error) {
	defer db.Log.Exit(db.Log.Enter())

	createdUser := db.Client.Create(&user)

	if createdUser.Error != nil {
		return nil, createdUser.Error
	}

	return createdUser, nil
}
