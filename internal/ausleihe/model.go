package ausleihe

import "time"

// Ausleihe repraesentiert einen Ausleihvorgang (1:N zu Mitglied).
type Ausleihe struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Ausleihdatum   time.Time `gorm:"column:ausleihdatum;not null"   json:"ausleihdatum"`
	Rueckgabedatum time.Time `gorm:"column:rueckgabedatum;not null" json:"rueckgabedatum"`
	MitgliedID     uint      `gorm:"column:mitglied_id;not null"    json:"mitglied_id"`
}

// TableName gibt den Tabellennamen im Schema "mitglied" an.
func (Ausleihe) TableName() string {
	return "mitglied.ausleihe"
}

// AusleiheInput enthaelt die validierten Felder fuer Create.
type AusleiheInput struct {
	Ausleihdatum   time.Time `json:"ausleihdatum"   validate:"required"`
	Rueckgabedatum time.Time `json:"rueckgabedatum" validate:"required,gtfield=Ausleihdatum"`
	MitgliedID     uint      `json:"mitglied_id"    validate:"required,min=1"`
}
