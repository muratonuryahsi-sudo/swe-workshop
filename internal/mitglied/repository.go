package mitglied

import (
	"errors"

	"github.com/muratonuryahsi-sudo/swe-workshop/internal/shared"
	"gorm.io/gorm"
)

// Repository kapselt den Datenbankzugriff für Mitglied.
type Repository struct {
	db *gorm.DB
}

// NewRepository erstellt ein neues Repository für Mitglied.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// FindByID lädt ein Mitglied anhand seiner ID.
// Existiert kein Mitglied mit dieser ID, wird ein *shared.NotFoundError zurückgegeben.
func (r *Repository) FindByID(id uint) (*Mitglied, error) {
	var mitglied Mitglied
	err := r.db.First(&mitglied, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, &shared.NotFoundError{Resource: "Mitglied", ID: id}
	}
	if err != nil {
		return nil, err
	}
	return &mitglied, nil
}

// FindAll lädt alle Mitglieder.
func (r *Repository) FindAll() ([]Mitglied, error) {
	var mitglieder []Mitglied
	if err := r.db.Find(&mitglieder).Error; err != nil {
		return nil, err
	}
	return mitglieder, nil
}

// Create speichert ein neues Mitglied in der Datenbank.
func (r *Repository) Create(mitglied *Mitglied) error {
	return r.db.Create(mitglied).Error
}
