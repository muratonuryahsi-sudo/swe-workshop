package security

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/muratonuryahsi-sudo/swe-workshop/internal/shared"
)

// NewVerifier baut einen Token-Verifizierer fuer die gegebene Keycloak-Config.
// Fuehrt einen OIDC-Discovery-Request gegen den Issuer aus (.well-known/openid-configuration).
func NewVerifier(ctx context.Context, cfg KeycloakConfig) (*oidc.IDTokenVerifier, error) {
	provider, err := oidc.NewProvider(ctx, cfg.Issuer)
	if err != nil {
		return nil, fmt.Errorf("security: Keycloak-Discovery fehlgeschlagen: %w", err)
	}
	return provider.Verifier(&oidc.Config{ClientID: cfg.Audience}), nil
}

// claims bildet die fuer die Rollenpruefung relevanten Felder des JWT ab.
// Keycloak legt Client-Rollen unter resource_access.<clientID>.roles ab.
type claims struct {
	ResourceAccess map[string]struct {
		Roles []string `json:"roles"`
	} `json:"resource_access"`
}

// RolesRequired prueft den Bearer-Token aus dem Authorization-Header gegen
// Keycloak und verlangt mindestens eine der angegebenen Rollen.
func RolesRequired(verifier *oidc.IDTokenVerifier, clientID string, roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rawToken, err := extractToken(r)
			if err != nil {
				shared.WriteError(w, err)
				return
			}

			idToken, err := verifier.Verify(r.Context(), rawToken)
			if err != nil {
				shared.WriteError(w, &shared.UnauthorizedError{Message: "Token nicht (mehr) gültig"})
				return
			}

			var c claims
			if err := idToken.Claims(&c); err != nil {
				shared.WriteError(w, &shared.UnauthorizedError{Message: "Token-Claims ungültig"})
				return
			}

			tokenRoles := c.ResourceAccess[clientID].Roles
			if !hasAnyRole(tokenRoles, roles) {
				shared.WriteError(w, &shared.ForbiddenError{Message: "Erforderliche Rolle nicht vorhanden"})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func extractToken(r *http.Request) (string, error) {
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		return "", &shared.UnauthorizedError{Message: "Authorization-Header fehlt"}
	}
	return strings.TrimPrefix(auth, "Bearer "), nil
}

func hasAnyRole(tokenRoles, required []string) bool {
	for _, requiredRole := range required {
		for _, tokenRole := range tokenRoles {
			if tokenRole == requiredRole {
				return true
			}
		}
	}
	return false
}
