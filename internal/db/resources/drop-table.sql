DROP INDEX IF EXISTS
    ausweis_mitglied_id_idx,
    ausleihe_mitglied_id_idx,
    mitglied_nachname_idx;

DROP TABLE IF EXISTS
    ausweis,
    ausleihe,
    mitglied CASCADE;

DROP TYPE IF EXISTS
    geschlecht,
    mitgliedsstatus;
