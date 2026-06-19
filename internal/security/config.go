// Package security enthält die JWT/OIDC-Verifizierung gegen Keycloak.
package security

import (
	"fmt"
	"os"
)

// KeycloakConfig enthält die fuer die Token-Verifizierung notwendigen Angaben.
type KeycloakConfig struct {
	Issuer   string
	ClientID string
	Audience string
}

// LoadKeycloakConfig liest die Keycloak-Konfiguration aus Umgebungsvariablen.
func LoadKeycloakConfig() KeycloakConfig {
	schema := getEnv("KEYCLOAK_SCHEMA", "http")
	host := getEnv("KEYCLOAK_HOST", "localhost")
	port := getEnv("KEYCLOAK_PORT", "8880")
	realm := getEnv("KEYCLOAK_REALM", "go")

	return KeycloakConfig{
		Issuer:   fmt.Sprintf("%s://%s:%s/realms/%s", schema, host, port, realm),
		ClientID: getEnv("KEYCLOAK_CLIENT_ID", "go-client"),
		Audience: getEnv("KEYCLOAK_AUDIENCE", "account"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}
