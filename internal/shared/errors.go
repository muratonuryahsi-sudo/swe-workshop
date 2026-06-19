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
