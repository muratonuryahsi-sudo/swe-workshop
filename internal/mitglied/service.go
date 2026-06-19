package mitglied

import (
	"github.com/go-playground/validator/v10"
	"github.com/muratonuryahsi-sudo/swe-workshop/internal/shared"
)

// Service enthält die Geschäftslogik für Mitglied.
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

// GetByID sucht ein Mitglied anhand seiner ID.
func (s *Service) GetByID(id uint) (*Mitglied, error) {
	return s.repo.FindByID(id)
}

// GetAll liefert alle Mitglieder.
func (s *Service) GetAll() ([]Mitglied, error) {
	return s.repo.FindAll()
}

// Create validiert die Eingabe und legt ein neues Mitglied an.
func (s *Service) Create(input *MitgliedInput) (*Mitglied, error) {
	if err := s.validate.Struct(input); err != nil {
		return nil, &shared.ValidationError{Message: err.Error()}
	}
	mitglied := &Mitglied{
		Vorname:         input.Vorname,
		Nachname:        input.Nachname,
		Email:           input.Email,
		Geburtsdatum:    input.Geburtsdatum,
		Telefonnummer:   input.Telefonnummer,
		Geschlecht:      input.Geschlecht,
		Mitgliedsstatus: input.Mitgliedsstatus,
		Beitrittsdatum:  input.Beitrittsdatum,
		Interessen:      input.Interessen,
	}
	if err := s.repo.Create(mitglied); err != nil {
		return nil, err
	}
	return mitglied, nil
}
