package db

import (
	"embed"
	"encoding/csv"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

//go:embed resources/drop-table.sql resources/create-table.sql resources/seed/*.csv
var resources embed.FS

// Populate setzt das Datenbankschema im Schema "mitglied" zurück und füllt es
// mit den Seed-Daten aus resources/seed/*.csv neu. Gedacht für lokale
// Entwicklung und Tests gegen einen bekannten, reproduzierbaren Datenstand.
func Populate(database *gorm.DB) error {
	if err := execSQLFile(database, "resources/drop-table.sql"); err != nil {
		return fmt.Errorf("db: Drop-Tabellen fehlgeschlagen: %w", err)
	}
	if err := execSQLFile(database, "resources/create-table.sql"); err != nil {
		return fmt.Errorf("db: Create-Tabellen fehlgeschlagen: %w", err)
	}
	if err := seedMitglieder(database); err != nil {
		return fmt.Errorf("db: Seed Mitglieder fehlgeschlagen: %w", err)
	}
	if err := seedAusweise(database); err != nil {
		return fmt.Errorf("db: Seed Ausweise fehlgeschlagen: %w", err)
	}
	if err := seedAusleihen(database); err != nil {
		return fmt.Errorf("db: Seed Ausleihen fehlgeschlagen: %w", err)
	}
	return nil
}

func execSQLFile(database *gorm.DB, path string) error {
	content, err := resources.ReadFile(path)
	if err != nil {
		return err
	}
	return database.Exec(string(content)).Error
}

// readSeedCSV liest eine Seed-CSV-Datei (Semikolon-getrennt, mit Header) ein
// und gibt die Datenzeilen ohne Header zurück.
func readSeedCSV(path string) ([][]string, error) {
	content, err := resources.ReadFile(path)
	if err != nil {
		return nil, err
	}
	reader := csv.NewReader(strings.NewReader(string(content)))
	reader.Comma = ';'
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, nil
	}
	return rows[1:], nil
}

// nullableValue wandelt das CSV-Literal "null" in einen SQL-NULL-Wert (nil) um.
func nullableValue(value string) any {
	if value == "null" {
		return nil
	}
	return value
}
