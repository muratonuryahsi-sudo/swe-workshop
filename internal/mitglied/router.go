package mitglied

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/muratonuryahsi-sudo/swe-workshop/internal/shared"
)

// Router registriert alle Mitglied-Endpunkte an den gegebenen chi-Router.
// authMiddleware schuetzt das Neuanlegen (POST) und kann nil sein, wenn keine
// Authentifizierung gewuenscht ist (z.B. in Tests ohne laufendes Keycloak).
func Router(svc *Service, authMiddleware func(http.Handler) http.Handler) http.Handler {
	r := chi.NewRouter()

	r.Get("/", getAll(svc))
	r.Get("/{id}", getByID(svc))

	if authMiddleware != nil {
		r.With(authMiddleware).Post("/", create(svc))
	} else {
		r.Post("/", create(svc))
	}

	return r
}

// getAll behandelt GET /mitglied
func getAll(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mitglieder, err := svc.GetAll()
		if err != nil {
			shared.WriteError(w, err)
			return
		}
		shared.WriteJSON(w, http.StatusOK, mitglieder)
	}
}

// getByID behandelt GET /mitglied/{id}
func getByID(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseID(chi.URLParam(r, "id"))
		if err != nil {
			shared.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "Ungültige ID"})
			return
		}
		mitglied, err := svc.GetByID(id)
		if err != nil {
			shared.WriteError(w, err)
			return
		}
		shared.WriteJSON(w, http.StatusOK, mitglied)
	}
}

// create behandelt POST /mitglied
func create(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input MitgliedInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			shared.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "Ungültiger Request-Body"})
			return
		}
		mitglied, err := svc.Create(&input)
		if err != nil {
			shared.WriteError(w, err)
			return
		}
		shared.WriteJSON(w, http.StatusCreated, mitglied)
	}
}

// parseID wandelt den String-URL-Parameter in eine uint-ID um.
func parseID(raw string) (uint, error) {
	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}
