package main

import (
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
)

func main() {
	// Konfiguration laden
	cfg := config.Load()

	// Datenbankverbindung herstellen
	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Datenbankverbindung fehlgeschlagen: %v", err)
	}

	// Datenbank befüllen
	if err := db.Populate(database); err != nil {
		log.Fatalf("Datenbank-Populate fehlgeschlagen: %v", err)
	}

	// Repositories
	ausweisRepo := ausweis.NewRepository(database)
	ausleiheRepo := ausleihe.NewRepository(database)
	mitgliedRepo := mitglied.NewRepository(database)

	// Services
	ausweisSvc := ausweis.NewService(ausweisRepo)
	ausleiheSvc := ausleihe.NewService(ausleiheRepo)
	mitgliedSvc := mitglied.NewService(mitgliedRepo)

	// Router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Mount("/ausweis", ausweis.Router(ausweisSvc))
	r.Mount("/ausleihe", ausleihe.Router(ausleiheSvc))
	r.Mount("/mitglied", mitglied.Router(mitgliedSvc))

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server laeuft auf %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server-Fehler: %v", err)
	}
}
