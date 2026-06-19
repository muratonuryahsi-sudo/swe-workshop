package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/muratonuryahsi-sudo/swe-workshop/internal/ausleihe"
	"github.com/muratonuryahsi-sudo/swe-workshop/internal/ausweis"
	"github.com/muratonuryahsi-sudo/swe-workshop/internal/config"
	"github.com/muratonuryahsi-sudo/swe-workshop/internal/db"
	"github.com/muratonuryahsi-sudo/swe-workshop/internal/mitglied"
	"github.com/muratonuryahsi-sudo/swe-workshop/internal/security"
)

func main() {
	// Konfiguration laden
	cfg := config.Load()

	// Datenbankverbindung herstellen
	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Datenbankverbindung fehlgeschlagen: %v", err)
	}

	// Datenbank befüllen (nur wenn explizit über DB_POPULATE=true angefordert,
	// da dies die bestehenden Daten löscht)
	if cfg.DBPopulate {
		if err := db.Populate(database); err != nil {
			log.Fatalf("Datenbank-Populate fehlgeschlagen: %v", err)
		}
	}

	// Repositories
	ausweisRepo := ausweis.NewRepository(database)
	ausleiheRepo := ausleihe.NewRepository(database)
	mitgliedRepo := mitglied.NewRepository(database)

	// Services
	ausweisSvc := ausweis.NewService(ausweisRepo)
	ausleiheSvc := ausleihe.NewService(ausleiheRepo)
	mitgliedSvc := mitglied.NewService(mitgliedRepo)

	// Auth-Middleware fuer das Neuanlegen (nur wenn AUTH_ENABLED=true, da OIDC
	// laut Aufgabenstellung optional ist)
	var authMiddleware func(http.Handler) http.Handler
	if cfg.AuthEnabled {
		kcCfg := security.LoadKeycloakConfig()
		verifier, err := security.NewVerifier(context.Background(), kcCfg)
		if err != nil {
			log.Fatalf("Keycloak-Verifizierer konnte nicht erstellt werden: %v", err)
		}
		authMiddleware = security.RolesRequired(verifier, kcCfg.ClientID, "admin", "user")
	}

	// Router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Mount("/ausweis", ausweis.Router(ausweisSvc, authMiddleware))
	r.Mount("/ausleihe", ausleihe.Router(ausleiheSvc, authMiddleware))
	r.Mount("/mitglied", mitglied.Router(mitgliedSvc, authMiddleware))

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server laeuft auf %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server-Fehler: %v", err)
	}
}
