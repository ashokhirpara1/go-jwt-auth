package mysql

import (
	"go-jwt/database/models"
)

// GetUserByEmail - create new user
func (db *DB) GetUserByEmail(email string) (models.User, error) {
	defer db.Log.Exit(db.Log.Enter())

	user := models.User{}
	if err := db.Client.Where("email = ?", email).First(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

// GetUsers -
func (db *DB) GetUsers() ([]models.User, error) {
	defer db.Log.Exit(db.Log.Enter())

	var users []models.User
	db.Client.Find(&users)

	return users, nil
}
