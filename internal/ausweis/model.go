package ausweis

import "time"

// Ausweis repraesentiert den Mitgliedsausweis (1:1 zu Mitglied).
type Ausweis struct {
	ID                uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Ausstellungsdatum time.Time `gorm:"column:ausstellungsdatum;not null" json:"ausstellungsdatum"`
	Ablaufdatum       time.Time `gorm:"column:ablaufdatum;not null"       json:"ablaufdatum"`
	MitgliedID        uint      `gorm:"column:mitglied_id;not null;unique" json:"mitglied_id"`
}

// TableName gibt den Tabellennamen im Schema "mitglied" an.
func (Ausweis) TableName() string {
	return "mitglied.ausweis"
}

// AusweisInput enthaelt die validierten Felder fuer Create/Update.
type AusweisInput struct {
	Ausstellungsdatum time.Time `json:"ausstellungsdatum" validate:"required"`
	Ablaufdatum       time.Time `json:"ablaufdatum"       validate:"required,gtfield=Ausstellungsdatum"`
	MitgliedID        uint      `json:"mitglied_id"       validate:"required,min=1"`
}
