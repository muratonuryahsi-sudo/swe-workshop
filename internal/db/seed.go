package db

import "gorm.io/gorm"

// seedMitglieder liest resources/seed/mitglied.csv und fügt die Zeilen ein.
// Spalten: id;version;vorname;nachname;email;geburtsdatum;telefonnummer;
//
//	geschlecht;mitgliedsstatus;beitrittsdatum;interessen;erzeugt;aktualisiert
func seedMitglieder(database *gorm.DB) error {
	rows, err := readSeedCSV("resources/seed/mitglied.csv")
	if err != nil {
		return err
	}
	const insert = `
		INSERT INTO mitglied
			(id, version, vorname, nachname, email, geburtsdatum, telefonnummer,
			 geschlecht, mitgliedsstatus, beitrittsdatum, interessen, erzeugt, aktualisiert)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?::jsonb, ?, ?)`
	for _, row := range rows {
		err := database.Exec(insert,
			row[0], row[1], row[2], row[3], row[4], row[5], row[6],
			row[7], row[8], row[9], nullableValue(row[10]), row[11], row[12],
		).Error
		if err != nil {
			return err
		}
	}
	return nil
}

// seedAusweise liest resources/seed/ausweis.csv und fügt die Zeilen ein.
// Spalten: id;ausstellungsdatum;ablaufdatum;mitglied_id
func seedAusweise(database *gorm.DB) error {
	rows, err := readSeedCSV("resources/seed/ausweis.csv")
	if err != nil {
		return err
	}
	const insert = `
		INSERT INTO ausweis (id, ausstellungsdatum, ablaufdatum, mitglied_id)
		VALUES (?, ?, ?, ?)`
	for _, row := range rows {
		if err := database.Exec(insert, row[0], row[1], row[2], row[3]).Error; err != nil {
			return err
		}
	}
	return nil
}

// seedAusleihen liest resources/seed/ausleihe.csv und fügt die Zeilen ein.
// Spalten: id;ausleihdatum;rueckgabedatum;mitglied_id
func seedAusleihen(database *gorm.DB) error {
	rows, err := readSeedCSV("resources/seed/ausleihe.csv")
	if err != nil {
		return err
	}
	const insert = `
		INSERT INTO ausleihe (id, ausleihdatum, rueckgabedatum, mitglied_id)
		VALUES (?, ?, ?, ?)`
	for _, row := range rows {
		if err := database.Exec(insert, row[0], row[1], row[2], row[3]).Error; err != nil {
			return err
		}
	}
	return nil
}
