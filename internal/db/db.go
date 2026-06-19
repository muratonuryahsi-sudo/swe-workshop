// Package db stellt die Verbindung zur bestehenden PostgreSQL-Datenbank "mitglied" her.
package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connect oeffnet eine GORM-Verbindung zur bestehenden Datenbank.
func Connect(databaseURL string) (*gorm.DB, error) {
	database, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("db: Verbindung fehlgeschlagen: %w", err)
	}
	return database, nil
}
