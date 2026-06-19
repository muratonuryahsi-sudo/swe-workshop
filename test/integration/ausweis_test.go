//go:build integration

package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/muratonuryahsi-sudo/swe-workshop/internal/ausweis"
	"github.com/muratonuryahsi-sudo/swe-workshop/internal/config"
	"github.com/muratonuryahsi-sudo/swe-workshop/internal/db"
)

func setupAusweisServer(t *testing.T) *httptest.Server {
	t.Helper()
	cfg := config.Load()
	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		t.Fatalf("DB-Verbindung fehlgeschlagen: %v", err)
	}
	repo := ausweis.NewRepository(database)
	svc := ausweis.NewService(repo)
	return httptest.NewServer(ausweis.Router(svc, nil))
}

func TestAusweisCreate(t *testing.T) {
	srv := setupAusweisServer(t)
	defer srv.Close()

	input := ausweis.AusweisInput{
		Ausstellungsdatum: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Ablaufdatum:       time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		MitgliedID:        createTestMitglied(t),
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

	var result ausweis.Ausweis
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Response-Body konnte nicht dekodiert werden: %v", err)
	}
	if result.ID == 0 {
		t.Error("erwartet ID > 0 in der Antwort")
	}
	if result.MitgliedID != input.MitgliedID {
		t.Errorf("erwartet MitgliedID %d, bekommen %d", input.MitgliedID, result.MitgliedID)
	}
}

func TestAusweisGetByID(t *testing.T) {
	srv := setupAusweisServer(t)
	defer srv.Close()

	// Erst anlegen
	input := ausweis.AusweisInput{
		Ausstellungsdatum: time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
		Ablaufdatum:       time.Date(2027, 2, 1, 0, 0, 0, 0, time.UTC),
		MitgliedID:        createTestMitglied(t),
	}
	body, _ := json.Marshal(input)
	postResp, err := http.Post(srv.URL+"/", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("POST fehlgeschlagen: %v", err)
	}
	defer postResp.Body.Close()

	var created ausweis.Ausweis
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

	var result ausweis.Ausweis
	if err := json.NewDecoder(getResp.Body).Decode(&result); err != nil {
		t.Fatalf("GET-Response konnte nicht dekodiert werden: %v", err)
	}
	if result.ID != created.ID {
		t.Errorf("erwartet ID %d, bekommen %d", created.ID, result.ID)
	}
}

func TestAusweisGetByID_NotFound(t *testing.T) {
	srv := setupAusweisServer(t)
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
