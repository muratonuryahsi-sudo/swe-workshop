// Package config liest die Konfiguration (DB-Connection-String, Port) aus Umgebungsvariablen.
package config

import "os"

// Config enthaelt die Laufzeitkonfiguration der Anwendung.
type Config struct {
	DatabaseURL string
	Port        string
	DBPopulate  bool
	AuthEnabled bool
}

// Load liest die Konfiguration aus Umgebungsvariablen und liefert sinnvolle
// Defaults fuer die lokale Entwicklung gegen die bestehende "mitglied"-Datenbank.
// DBPopulate ist standardmaessig false, da es die Datenbank zuruecksetzt und
// neu befuellt - nur explizit ueber DB_POPULATE=true aktivieren.
// AuthEnabled ist standardmaessig false, da OIDC/Keycloak laut Aufgabenstellung
// optional ist - der Server soll auch ohne laufendes Keycloak starten.
func Load() Config {
	return Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://mitglied:p@localhost:5432/mitglied?search_path=mitglied&sslmode=disable"),
		Port:        getEnv("PORT", "8080"),
		DBPopulate:  getEnv("DB_POPULATE", "false") == "true",
		AuthEnabled: getEnv("AUTH_ENABLED", "false") == "true",
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}
