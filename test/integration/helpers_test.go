//go:build integration

package integration

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/muratonuryahsi-sudo/swe-workshop/internal/config"
	"github.com/muratonuryahsi-sudo/swe-workshop/internal/db"
	"github.com/muratonuryahsi-sudo/swe-workshop/internal/mitglied"
)

func uintToStr(id uint) string {
	return strconv.FormatUint(uint64(id), 10)
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
