package ausweis

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/muratonuryahsi-sudo/swe-workshop/internal/shared"
)

// Router registriert alle Ausweis-Endpunkte an den gegebenen chi-Router.
// authMiddleware schuetzt das Neuanlegen (POST) und kann nil sein, wenn keine
// Authentifizierung gewuenscht ist (z.B. in Tests ohne laufendes Keycloak).
func Router(svc *Service, authMiddleware func(http.Handler) http.Handler) http.Handler {
	r := chi.NewRouter()

	r.Get("/{id}", getByID(svc))
	r.Get("/mitglied/{mitgliedID}", getByMitgliedID(svc))

	if authMiddleware != nil {
		r.With(authMiddleware).Post("/", create(svc))
	} else {
		r.Post("/", create(svc))
	}

	r.Put("/{id}", update(svc))
	r.Delete("/{id}", delete(svc))

	return r
}

// getByID behandelt GET /ausweis/{id}
func getByID(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseID(chi.URLParam(r, "id"))
		if err != nil {
			shared.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "Ungueltige ID"})
			return
		}
		ausweis, err := svc.GetByID(id)
		if err != nil {
			shared.WriteError(w, err)
			return
		}
		shared.WriteJSON(w, http.StatusOK, ausweis)
	}
}

// getByMitgliedID behandelt GET /ausweis/mitglied/{mitgliedID}
func getByMitgliedID(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mitgliedID, err := parseID(chi.URLParam(r, "mitgliedID"))
		if err != nil {
			shared.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "Ungueltige MitgliedID"})
			return
		}
		ausweis, err := svc.GetByMitgliedID(mitgliedID)
		if err != nil {
			shared.WriteError(w, err)
			return
		}
		shared.WriteJSON(w, http.StatusOK, ausweis)
	}
}

// create behandelt POST /ausweis
func create(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input AusweisInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			shared.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "Ungueltige JSON-Daten"})
			return
		}
		ausweis, err := svc.Create(&input)
		if err != nil {
			shared.WriteError(w, err)
			return
		}
		shared.WriteJSON(w, http.StatusCreated, ausweis)
	}
}

// update behandelt PUT /ausweis/{id}
func update(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseID(chi.URLParam(r, "id"))
		if err != nil {
			shared.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "Ungueltige ID"})
			return
		}
		var input AusweisInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			shared.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "Ungueltige JSON-Daten"})
			return
		}
		ausweis, err := svc.Update(id, &input)
		if err != nil {
			shared.WriteError(w, err)
			return
		}
		shared.WriteJSON(w, http.StatusOK, ausweis)
	}
}

// delete behandelt DELETE /ausweis/{id}
func delete(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseID(chi.URLParam(r, "id"))
		if err != nil {
			shared.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "Ungueltige ID"})
			return
		}
		if err := svc.Delete(id); err != nil {
			shared.WriteError(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func parseID(s string) (uint, error) {
	id, err := strconv.ParseUint(s, 10, 64)
	return uint(id), err
}
