package ausleihe

import (
	"errors"

	"gorm.io/gorm"
)

// Repository kapselt den Datenbankzugriff fuer Ausleihe.
type Repository struct {
	db *gorm.DB
}

// NewRepository erstellt ein neues Repository mit der gegebenen DB-Verbindung.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// FindByMitgliedID gibt alle Ausleihen eines Mitglieds zurueck.
func (r *Repository) FindByMitgliedID(mitgliedID uint) ([]Ausleihe, error) {
	var ausleihen []Ausleihe
	result := r.db.Where("mitglied_id = ?", mitgliedID).Find(&ausleihen)
	return ausleihen, result.Error
}

// FindByID sucht eine Ausleihe anhand ihrer ID.
func (r *Repository) FindByID(id uint) (*Ausleihe, error) {
	var ausleihe Ausleihe
	result := r.db.First(&ausleihe, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &ausleihe, result.Error
}

// Create speichert eine neue Ausleihe in der Datenbank.
func (r *Repository) Create(ausleihe *Ausleihe) error {
	return r.db.Create(ausleihe).Error
}

// Update aktualisiert eine bestehende Ausleihe.
func (r *Repository) Update(ausleihe *Ausleihe) error {
	return r.db.Save(ausleihe).Error
}

// Delete loescht die Ausleihe mit der gegebenen ID.
func (r *Repository) Delete(id uint) error {
	return r.db.Delete(&Ausleihe{}, id).Error
}
