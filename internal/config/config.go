// Package config liest die Konfiguration (DB-Connection-String, Port) aus Umgebungsvariablen.
package config

import "os"

// Config enthaelt die Laufzeitkonfiguration der Anwendung.
type Config struct {
	DatabaseURL string
	Port        string
}

// Load liest die Konfiguration aus Umgebungsvariablen und liefert sinnvolle
// Defaults fuer die lokale Entwicklung gegen die bestehende "mitglied"-Datenbank.
func Load() Config {
	return Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://mitglied:p@localhost:5432/mitglied?search_path=mitglied&sslmode=disable"),
		Port:        getEnv("PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}
