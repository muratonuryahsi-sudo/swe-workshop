package ausweis

import (
	"errors"

	"gorm.io/gorm"
)

// Repository kapselt den Datenbankzugriff fuer Ausweis.
type Repository struct {
	db *gorm.DB
}

// NewRepository erstellt ein neues Repository mit der gegebenen DB-Verbindung.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// FindByMitgliedID sucht den Ausweis anhand der Mitglied-ID.
func (r *Repository) FindByMitgliedID(mitgliedID uint) (*Ausweis, error) {
	var ausweis Ausweis
	result := r.db.Where("mitglied_id = ?", mitgliedID).First(&ausweis)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &ausweis, result.Error
}

// FindByID sucht den Ausweis anhand seiner eigenen ID.
func (r *Repository) FindByID(id uint) (*Ausweis, error) {
	var ausweis Ausweis
	result := r.db.First(&ausweis, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &ausweis, result.Error
}

// Create speichert einen neuen Ausweis in der Datenbank.
func (r *Repository) Create(ausweis *Ausweis) error {
	return r.db.Create(ausweis).Error
}

// Update aktualisiert einen bestehenden Ausweis.
func (r *Repository) Update(ausweis *Ausweis) error {
	return r.db.Save(ausweis).Error
}

// Delete loescht den Ausweis mit der gegebenen ID.
func (r *Repository) Delete(id uint) error {
	return r.db.Delete(&Ausweis{}, id).Error
}
