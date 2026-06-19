//go:build integration

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/muratonuryahsi-sudo/swe-workshop/internal/ausleihe"
	"github.com/muratonuryahsi-sudo/swe-workshop/internal/config"
	"github.com/muratonuryahsi-sudo/swe-workshop/internal/db"
	"github.com/muratonuryahsi-sudo/swe-workshop/internal/security"
)

func setupAusleiheServerWithAuth(t *testing.T) *httptest.Server {
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
	repo := ausleihe.NewRepository(database)
	svc := ausleihe.NewService(repo)
	return httptest.NewServer(ausleihe.Router(svc, authMiddleware))
}

func TestAusleiheCreate_Unauthorized(t *testing.T) {
	srv := setupAusleiheServerWithAuth(t)
	defer srv.Close()

	input := ausleihe.AusleiheInput{
		Ausleihdatum:   time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
		Rueckgabedatum: time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC),
		MitgliedID:     createTestMitglied(t),
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

func TestAusleiheCreate_WithToken(t *testing.T) {
	srv := setupAusleiheServerWithAuth(t)
	defer srv.Close()

	token := fetchToken(t)
	input := ausleihe.AusleiheInput{
		Ausleihdatum:   time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
		Rueckgabedatum: time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC),
		MitgliedID:     createTestMitglied(t),
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
