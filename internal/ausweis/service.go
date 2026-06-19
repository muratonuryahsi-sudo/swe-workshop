package ausweis

import (
	"github.com/go-playground/validator/v10"
	"github.com/muratonuryahsi-sudo/swe-workshop/internal/shared"
)

// Service enthaelt die Geschaeftslogik fuer Ausweis.
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

// GetByMitgliedID sucht den Ausweis eines Mitglieds.
func (s *Service) GetByMitgliedID(mitgliedID uint) (*Ausweis, error) {
	ausweis, err := s.repo.FindByMitgliedID(mitgliedID)
	if err != nil {
		return nil, err
	}
	if ausweis == nil {
		return nil, &shared.NotFoundError{Resource: "Ausweis", ID: mitgliedID}
	}
	return ausweis, nil
}

// GetByID sucht einen Ausweis anhand seiner ID.
func (s *Service) GetByID(id uint) (*Ausweis, error) {
	ausweis, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if ausweis == nil {
		return nil, &shared.NotFoundError{Resource: "Ausweis", ID: id}
	}
	return ausweis, nil
}

// Create legt einen neuen Ausweis an.
func (s *Service) Create(input *AusweisInput) (*Ausweis, error) {
	if err := s.validate.Struct(input); err != nil {
		return nil, &shared.ValidationError{Message: err.Error()}
	}
	ausweis := &Ausweis{
		Ausstellungsdatum: input.Ausstellungsdatum,
		Ablaufdatum:       input.Ablaufdatum,
		MitgliedID:        input.MitgliedID,
	}
	if err := s.repo.Create(ausweis); err != nil {
		return nil, err
	}
	return ausweis, nil
}

// Update aktualisiert einen bestehenden Ausweis.
func (s *Service) Update(id uint, input *AusweisInput) (*Ausweis, error) {
	if err := s.validate.Struct(input); err != nil {
		return nil, &shared.ValidationError{Message: err.Error()}
	}
	ausweis, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if ausweis == nil {
		return nil, &shared.NotFoundError{Resource: "Ausweis", ID: id}
	}
	ausweis.Ausstellungsdatum = input.Ausstellungsdatum
	ausweis.Ablaufdatum = input.Ablaufdatum
	ausweis.MitgliedID = input.MitgliedID
	if err := s.repo.Update(ausweis); err != nil {
		return nil, err
	}
	return ausweis, nil
}

// Delete loescht einen Ausweis anhand seiner ID.
func (s *Service) Delete(id uint) error {
	ausweis, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if ausweis == nil {
		return &shared.NotFoundError{Resource: "Ausweis", ID: id}
	}
	return s.repo.Delete(id)
}
