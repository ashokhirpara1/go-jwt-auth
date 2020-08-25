package mysql

import (
	"go-jwt/database/models"
	"go-jwt/logging"

	"github.com/jinzhu/gorm"
)

// Migrate Struct
type Migrate struct {
	db  *gorm.DB
	Log *logging.Handler
}

// RunMigration - Create Table in Database
func (m *Migrate) RunMigration() {
	m.Log.Info("Start Miration")

	// Check model `User`'s table exists or not
	if !m.db.HasTable(&models.User{}) {
		m.Log.Info("Creating Table User")

		// Create table for model `User`
		m.db.CreateTable(&models.User{})

		m.Log.Info("Created Table User")
	} else {
		m.Log.Info("Exists Table User")
	}

	m.Log.Info("End Miration")
}
