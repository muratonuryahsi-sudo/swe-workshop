package ausleihe

import (
	"github.com/go-playground/validator/v10"
	"github.com/muratonuryahsi-sudo/swe-workshop/internal/shared"
)

// Service enthaelt die Geschaeftslogik fuer Ausleihe.
type Service struct {
	repo     *Repository
	validate *validator.Validate
}

// NewService erstellt einen neuen Service mit dem gegebenen Repository.
func NewService(repo *Repository) *Service {
	return &Service{
		repo:     repo,
		validate: validator.New(),
	}
}

// GetByMitgliedID gibt alle Ausleihen eines Mitglieds zurueck.
func (s *Service) GetByMitgliedID(mitgliedID uint) ([]Ausleihe, error) {
	ausleihen, err := s.repo.FindByMitgliedID(mitgliedID)
	if err != nil {
		return nil, err
	}
	if len(ausleihen) == 0 {
		return nil, &shared.NotFoundError{Resource: "Ausleihen fuer Mitglied", ID: mitgliedID}
	}
	return ausleihen, nil
}

// GetByID sucht eine Ausleihe anhand ihrer ID.
func (s *Service) GetByID(id uint) (*Ausleihe, error) {
	ausleihe, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if ausleihe == nil {
		return nil, &shared.NotFoundError{Resource: "Ausleihe", ID: id}
	}
	return ausleihe, nil
}

// Create legt eine neue Ausleihe an.
func (s *Service) Create(input *AusleiheInput) (*Ausleihe, error) {
	if err := s.validate.Struct(input); err != nil {
		return nil, &shared.ValidationError{Message: err.Error()}
	}
	ausleihe := &Ausleihe{
		Ausleihdatum:   input.Ausleihdatum,
		Rueckgabedatum: input.Rueckgabedatum,
		MitgliedID:     input.MitgliedID,
	}
	if err := s.repo.Create(ausleihe); err != nil {
		return nil, err
	}
	return ausleihe, nil
}

// Update aktualisiert eine bestehende Ausleihe.
func (s *Service) Update(id uint, input *AusleiheInput) (*Ausleihe, error) {
	if err := s.validate.Struct(input); err != nil {
		return nil, &shared.ValidationError{Message: err.Error()}
	}
	ausleihe, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if ausleihe == nil {
		return nil, &shared.NotFoundError{Resource: "Ausleihe", ID: id}
	}
	ausleihe.Ausleihdatum = input.Ausleihdatum
	ausleihe.Rueckgabedatum = input.Rueckgabedatum
	ausleihe.MitgliedID = input.MitgliedID
	if err := s.repo.Update(ausleihe); err != nil {
		return nil, err
	}
	return ausleihe, nil
}

// Delete loescht eine Ausleihe anhand ihrer ID.
func (s *Service) Delete(id uint) error {
	ausleihe, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if ausleihe == nil {
		return &shared.NotFoundError{Resource: "Ausleihe", ID: id}
	}
	return s.repo.Delete(id)
}
