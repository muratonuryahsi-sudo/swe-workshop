//go:build integration

package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/muratonuryahsi-sudo/swe-workshop/internal/config"
	"github.com/muratonuryahsi-sudo/swe-workshop/internal/db"
	"github.com/muratonuryahsi-sudo/swe-workshop/internal/mitglied"
	"github.com/muratonuryahsi-sudo/swe-workshop/internal/security"
)

func uintToStr(id uint) string {
	return strconv.FormatUint(uint64(id), 10)
}

// fetchToken holt einen Access-Token von Keycloak via Password-Grant (User "user"/"p").
func fetchToken(t *testing.T) string {
	t.Helper()
	kcCfg := security.LoadKeycloakConfig()
	clientSecret := os.Getenv("KEYCLOAK_CLIENT_SECRET")
	if clientSecret == "" {
		clientSecret = "changeme-go-client-secret"
	}
	tokenURL := kcCfg.Issuer + "/protocol/openid-connect/token"

	resp, err := http.PostForm(tokenURL, url.Values{
		"grant_type":    {"password"},
		"client_id":     {kcCfg.ClientID},
		"client_secret": {clientSecret},
		"username":      {"user"},
		"password":      {"p"},
	})
	if err != nil {
		t.Fatalf("Token-Anfrage an Keycloak fehlgeschlagen: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Token-Response konnte nicht dekodiert werden: %v", err)
	}
	if result.AccessToken == "" {
		t.Fatal("Keycloak hat keinen Access-Token zurueckgegeben")
	}
	return result.AccessToken
}

// authHeader gibt den Authorization-Header-Wert fuer einen Bearer-Token zurueck.
func authHeader(token string) string {
	return "Bearer " + strings.TrimSpace(token)
}

// createTestMitglied legt ein frisches, eindeutiges Mitglied an und liefert
// dessen ID zurück. So sind die Tests unabhängig vom aktuellen Datenstand der
// DB (z.B. ob die Seed-Mitglieder existieren oder bereits einen Ausweis haben).
func createTestMitglied(t *testing.T) uint {
	t.Helper()

	cfg := config.Load()
	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		t.Fatalf("DB-Verbindung fehlgeschlagen: %v", err)
	}
	repo := mitglied.NewRepository(database)
	svc := mitglied.NewService(repo)

	email := fmt.Sprintf("integrationtest-%d@example.com", time.Now().UnixNano())
	input := &mitglied.MitgliedInput{
		Vorname:  "Integrationstest",
		Nachname: "Mitglied",
		Email:    email,
	}
	created, err := svc.Create(input)
	if err != nil {
		t.Fatalf("Test-Mitglied konnte nicht angelegt werden: %v", err)
	}
	return created.ID
}
