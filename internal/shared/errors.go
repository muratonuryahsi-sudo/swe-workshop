// Package shared enthaelt von allen Entitaeten gemeinsam genutzte Fehlertypen und HTTP-Hilfsfunktionen.
package shared

import "fmt"

// NotFoundError zeigt an, dass eine angeforderte Ressource nicht existiert.
type NotFoundError struct {
	Resource string
	ID       any
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s mit ID %v nicht gefunden", e.Resource, e.ID)
}

// ValidationError zeigt an, dass Eingabedaten die fachlichen Regeln verletzen.
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// VersionOutdatedError zeigt an, dass beim optimistischen Locking eine veraltete Version verwendet wurde.
type VersionOutdatedError struct {
	Resource string
	ID       any
}

func (e *VersionOutdatedError) Error() string {
	return fmt.Sprintf("%s mit ID %v wurde zwischenzeitlich geaendert (Version veraltet)", e.Resource, e.ID)
}

// UnauthorizedError zeigt an, dass kein gueltiger Bearer-Token vorhanden ist.
type UnauthorizedError struct {
	Message string
}

func (e *UnauthorizedError) Error() string {
	return e.Message
}

// ForbiddenError zeigt an, dass der Token zwar gueltig ist, aber die
// erforderliche Rolle fehlt.
type ForbiddenError struct {
	Message string
}

func (e *ForbiddenError) Error() string {
	return e.Message
}
