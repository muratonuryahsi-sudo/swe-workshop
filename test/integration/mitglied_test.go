//go:build integration

package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/muratonuryahsi-sudo/swe-workshop/internal/config"
	"github.com/muratonuryahsi-sudo/swe-workshop/internal/db"
	"github.com/muratonuryahsi-sudo/swe-workshop/internal/mitglied"
)

func setupMitgliedServer(t *testing.T) *httptest.Server {
	t.Helper()
	cfg := config.Load()
	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		t.Fatalf("DB-Verbindung fehlgeschlagen: %v", err)
	}
	repo := mitglied.NewRepository(database)
	svc := mitglied.NewService(repo)
	return httptest.NewServer(mitglied.Router(svc))
}

func TestMitgliedCreate(t *testing.T) {
	srv := setupMitgliedServer(t)
	defer srv.Close()

	input := mitglied.MitgliedInput{
		Vorname:  "Integrationstest",
		Nachname: "Create",
		Email:    fmt.Sprintf("create-%d@example.com", time.Now().UnixNano()),
	}
	body, _ := json.Marshal(input)

	resp, err := http.Post(srv.URL+"/", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("POST fehlgeschlagen: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("erwartet 201, bekommen %d", resp.StatusCode)
	}

	var result mitglied.Mitglied
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Response-Body konnte nicht dekodiert werden: %v", err)
	}
	if result.ID == 0 {
		t.Error("erwartet ID > 0 in der Antwort")
	}
	if result.Email != input.Email {
		t.Errorf("erwartet Email %q, bekommen %q", input.Email, result.Email)
	}
}

func TestMitgliedCreate_ValidationError(t *testing.T) {
	srv := setupMitgliedServer(t)
	defer srv.Close()

	input := mitglied.MitgliedInput{
		Vorname:  "Ohne",
		Nachname: "Email",
		Email:    "keine-gueltige-email",
	}
	body, _ := json.Marshal(input)

	resp, err := http.Post(srv.URL+"/", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("POST fehlgeschlagen: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("erwartet 422, bekommen %d", resp.StatusCode)
	}
}

func TestMitgliedGetByID(t *testing.T) {
	srv := setupMitgliedServer(t)
	defer srv.Close()

	// Erst anlegen
	input := mitglied.MitgliedInput{
		Vorname:  "Integrationstest",
		Nachname: "GetByID",
		Email:    fmt.Sprintf("getbyid-%d@example.com", time.Now().UnixNano()),
	}
	body, _ := json.Marshal(input)
	postResp, err := http.Post(srv.URL+"/", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("POST fehlgeschlagen: %v", err)
	}
	defer postResp.Body.Close()

	var created mitglied.Mitglied
	if err := json.NewDecoder(postResp.Body).Decode(&created); err != nil {
		t.Fatalf("POST-Response konnte nicht dekodiert werden: %v", err)
	}

	// Dann per ID lesen
	getResp, err := http.Get(srv.URL + "/" + uintToStr(created.ID))
	if err != nil {
		t.Fatalf("GET fehlgeschlagen: %v", err)
	}
	defer getResp.Body.Close()

	if getResp.StatusCode != http.StatusOK {
		t.Errorf("erwartet 200, bekommen %d", getResp.StatusCode)
	}

	var result mitglied.Mitglied
	if err := json.NewDecoder(getResp.Body).Decode(&result); err != nil {
		t.Fatalf("GET-Response konnte nicht dekodiert werden: %v", err)
	}
	if result.ID != created.ID {
		t.Errorf("erwartet ID %d, bekommen %d", created.ID, result.ID)
	}
}

func TestMitgliedGetByID_NotFound(t *testing.T) {
	srv := setupMitgliedServer(t)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/999999")
	if err != nil {
		t.Fatalf("GET fehlgeschlagen: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("erwartet 404, bekommen %d", resp.StatusCode)
	}
}

func TestMitgliedGetAll(t *testing.T) {
	srv := setupMitgliedServer(t)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/")
	if err != nil {
		t.Fatalf("GET fehlgeschlagen: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("erwartet 200, bekommen %d", resp.StatusCode)
	}

	var result []mitglied.Mitglied
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Response-Body konnte nicht dekodiert werden: %v", err)
	}
	if len(result) == 0 {
		t.Error("erwartet mindestens ein Mitglied in der Liste")
	}
}
