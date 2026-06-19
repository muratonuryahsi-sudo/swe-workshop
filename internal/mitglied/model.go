// Package mitglied implementiert Router, Service und Repository für die Entität Mitglied.
package mitglied

import (
	"encoding/json"
	"time"
)

// Geschlecht entspricht dem Postgres-Enum "geschlecht".
type Geschlecht string

const (
	GeschlechtMaennlich Geschlecht = "MAENNLICH"
	GeschlechtWeiblich  Geschlecht = "WEIBLICH"
	GeschlechtDivers    Geschlecht = "DIVERS"
)

// Mitgliedsstatus entspricht dem Postgres-Enum "mitgliedsstatus".
type Mitgliedsstatus string

const (
	MitgliedsstatusAktiv   Mitgliedsstatus = "AKTIV"
	MitgliedsstatusInaktiv Mitgliedsstatus = "INAKTIV"
)

// Mitglied repräsentiert ein Vereinsmitglied.
type Mitglied struct {
	ID              uint             `gorm:"primaryKey;autoIncrement"        json:"id"`
	Version         int              `gorm:"column:version;not null"         json:"version"`
	Vorname         string           `gorm:"column:vorname;not null"         json:"vorname"`
	Nachname        string           `gorm:"column:nachname;not null"        json:"nachname"`
	Email           string           `gorm:"column:email;not null;unique"    json:"email"`
	Geburtsdatum    *time.Time       `gorm:"column:geburtsdatum"             json:"geburtsdatum,omitempty"`
	Telefonnummer   *string          `gorm:"column:telefonnummer;unique"     json:"telefonnummer,omitempty"`
	Geschlecht      *Geschlecht      `gorm:"column:geschlecht;type:geschlecht"             json:"geschlecht,omitempty"`
	Mitgliedsstatus *Mitgliedsstatus `gorm:"column:mitgliedsstatus;type:mitgliedsstatus" json:"mitgliedsstatus,omitempty"`
	Beitrittsdatum  *time.Time       `gorm:"column:beitrittsdatum"           json:"beitrittsdatum,omitempty"`
	Interessen      json.RawMessage  `gorm:"column:interessen;type:jsonb"    json:"interessen,omitempty"`
	Erzeugt         time.Time        `gorm:"column:erzeugt;autoCreateTime"   json:"erzeugt"`
	Aktualisiert    time.Time        `gorm:"column:aktualisiert;autoUpdateTime" json:"aktualisiert"`
}

// TableName gibt den Tabellennamen im Schema "mitglied" an.
func (Mitglied) TableName() string {
	return "mitglied.mitglied"
}

// MitgliedInput enthält die validierten Felder für Create/Update.
type MitgliedInput struct {
	Vorname         string           `json:"vorname"         validate:"required,max=64"`
	Nachname        string           `json:"nachname"        validate:"required,max=64"`
	Email           string           `json:"email"           validate:"required,email"`
	Geburtsdatum    *time.Time       `json:"geburtsdatum"    validate:"omitempty"`
	Telefonnummer   *string          `json:"telefonnummer"   validate:"omitempty"`
	Geschlecht      *Geschlecht      `json:"geschlecht"      validate:"omitempty,oneof=MAENNLICH WEIBLICH DIVERS"`
	Mitgliedsstatus *Mitgliedsstatus `json:"mitgliedsstatus" validate:"omitempty,oneof=AKTIV INAKTIV"`
	Beitrittsdatum  *time.Time       `json:"beitrittsdatum"  validate:"omitempty"`
	Interessen      json.RawMessage  `json:"interessen"      validate:"omitempty"`
}
