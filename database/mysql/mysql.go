package mysql

import (
	"go-jwt/logging"

	"github.com/jinzhu/gorm"

	//
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// DB holds database client
type DB struct {
	Client *gorm.DB
	Log    *logging.Handler
}

// Get database client
func Get(connStr string, dLog *logging.Handler) (*DB, error) {
	dLog.Info("Get Mysql")
	dLog.Info("Conection String: " + connStr)

	db, err := get(connStr)
	if err != nil {
		dLog.DBError("Failed to Connect DB: ", err)
		return nil, err
	}

	// Run Migration Commands
	m := &Migrate{db: db, Log: dLog}
	m.RunMigration()

	return &DB{
		Client: db,
		Log:    dLog,
	}, nil
}

func get(connStr string) (*gorm.DB, error) {
	db, err := gorm.Open("mysql", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.DB().Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
