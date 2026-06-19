//go:build integration

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/muratonuryahsi-sudo/swe-workshop/internal/config"
	"github.com/muratonuryahsi-sudo/swe-workshop/internal/db"
	"github.com/muratonuryahsi-sudo/swe-workshop/internal/mitglied"
	"github.com/muratonuryahsi-sudo/swe-workshop/internal/security"
)

func setupMitgliedServerWithAuth(t *testing.T) *httptest.Server {
	t.Helper()
	cfg := config.Load()
	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		t.Fatalf("DB-Verbindung fehlgeschlagen: %v", err)
	}
	kcCfg := security.LoadKeycloakConfig()
	verifier, err := security.NewVerifier(context.Background(), kcCfg)
	if err != nil {
		t.Fatalf("Keycloak-Verifizierer konnte nicht erstellt werden: %v", err)
	}
	authMiddleware := security.RolesRequired(verifier, kcCfg.ClientID, "admin", "user")
	repo := mitglied.NewRepository(database)
	svc := mitglied.NewService(repo)
	return httptest.NewServer(mitglied.Router(svc, authMiddleware))
}

func TestMitgliedCreate_Unauthorized(t *testing.T) {
	srv := setupMitgliedServerWithAuth(t)
	defer srv.Close()

	input := mitglied.MitgliedInput{
		Vorname:  "Auth",
		Nachname: "Test",
		Email:    fmt.Sprintf("auth-unauth-%d@example.com", time.Now().UnixNano()),
	}
	body, _ := json.Marshal(input)

	resp, err := http.Post(srv.URL+"/", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("POST fehlgeschlagen: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("erwartet 401, bekommen %d", resp.StatusCode)
	}
}

func TestMitgliedCreate_WithToken(t *testing.T) {
	srv := setupMitgliedServerWithAuth(t)
	defer srv.Close()

	token := fetchToken(t)
	input := mitglied.MitgliedInput{
		Vorname:  "Auth",
		Nachname: "Test",
		Email:    fmt.Sprintf("auth-token-%d@example.com", time.Now().UnixNano()),
	}
	body, _ := json.Marshal(input)

	req, err := http.NewRequest(http.MethodPost, srv.URL+"/", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("Request konnte nicht erstellt werden: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authHeader(token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("POST fehlgeschlagen: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("erwartet 201, bekommen %d", resp.StatusCode)
	}
}
