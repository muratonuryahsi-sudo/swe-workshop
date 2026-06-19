# swe-workshop

> **Mitgliederverwaltungs-API** – Ein Go-REST-Server für Vereinsmitglieder, Ausweise und Ausleihen.  
> Entwickelt im Rahmen eines Software-Engineering-Workshops an der **Hochschule Karlsruhe**.

---

## Inhaltsverzeichnis

- [Überblick](#überblick)
- [Architektur](#architektur)
- [Tech-Stack](#tech-stack)
- [API-Endpunkte](#api-endpunkte)
- [Datenmodell](#datenmodell)
- [Erste Schritte](#erste-schritte)
  - [Voraussetzungen](#voraussetzungen)
  - [Konfiguration](#konfiguration)
  - [Lokaler Start (ohne Docker)](#lokaler-start-ohne-docker)
  - [Mit Docker](#mit-docker)
- [Authentifizierung (Keycloak / OIDC)](#authentifizierung-keycloak--oidc)
- [Datenbank-Seeding](#datenbank-seeding)
- [Tests](#tests)
- [CI/CD](#cicd)
- [Projektstruktur](#projektstruktur)
- [Lizenz](#lizenz)

---

## Überblick

`swe-workshop` ist ein leichtgewichtiger HTTP-REST-Server in **Go**, der die drei Kernentitäten einer Vereinsverwaltung abbildet:

| Entität | Beschreibung |
|---|---|
| **Mitglied** | Vereinsmitglied mit Stammdaten (Name, E-Mail, Geburtsdatum, etc.) |
| **Ausweis** | 1:1-Mitgliedsausweis mit Ausstellungs- und Ablaufdatum |
| **Ausleihe** | N:1-Ausleihvorgang eines Mitglieds mit Ausleihdatum und Rückgabedatum |

Die API ist vollständig **JSON-basiert**, optional **OIDC-gesichert** (Keycloak) und liefert aussagekräftige Fehlermeldungen mit passenden HTTP-Statuscodes.

---

## Architektur

Das Projekt folgt einer **layered architecture** mit klarer Trennung zwischen Router, Service und Repository:

```
cmd/server/
└── main.go             ← Einstiegspunkt: Konfiguration, DI, Server-Start

internal/
├── config/             ← Konfiguration aus Umgebungsvariablen
├── db/                 ← Datenbankverbindung, Migrate & Seeding
├── mitglied/           ← Router · Service · Repository · Model
├── ausweis/            ← Router · Service · Repository · Model
├── ausleihe/           ← Router · Service · Repository · Model
├── dev/                ← Entwicklungs-Endpunkt: POST /dev/db-populate
├── security/           ← OIDC-Middleware (Keycloak)
└── shared/             ← Gemeinsame Fehlertypen & HTTP-Hilfsfunktionen

test/integration/       ← Integrationstests (Build-Tag: integration)
extras/
├── bruno/              ← Bruno API-Collection
└── keycloak/           ← Keycloak-Setup-Skript & Dokumentation
```

---

## Tech-Stack

| Kategorie | Technologie |
|---|---|
| Sprache | [Go 1.25](https://go.dev/) |
| HTTP-Router | [go-chi/chi v5](https://github.com/go-chi/chi) |
| ORM | [GORM](https://gorm.io/) mit PostgreSQL-Treiber |
| Datenbank | [PostgreSQL](https://www.postgresql.org/) |
| Authentifizierung | [Keycloak](https://www.keycloak.org/) via [coreos/go-oidc v3](https://github.com/coreos/go-oidc) |
| Validierung | [go-playground/validator v10](https://github.com/go-playground/validator) |
| Container | Docker (Distroless `gcr.io/distroless/static-debian12` oder Chainguard `cgr.io/chainguard/static`) |
| CI | GitHub Actions |
| API-Client | [Bruno](https://www.usebruno.com/) |

---

## API-Endpunkte

### Mitglied `/mitglied`

| Methode | Pfad | Beschreibung | Auth erforderlich |
|---|---|---|---|
| `GET` | `/mitglied` | Alle Mitglieder abrufen | Nein |
| `GET` | `/mitglied/{id}` | Einzelnes Mitglied abrufen | Nein |
| `POST` | `/mitglied` | Neues Mitglied anlegen | Ja _(wenn `AUTH_ENABLED=true`)_ |

### Ausweis `/ausweis`

| Methode | Pfad | Beschreibung | Auth erforderlich |
|---|---|---|---|
| `GET` | `/ausweis/{id}` | Ausweis nach ID abrufen | Nein |
| `GET` | `/ausweis/mitglied/{mitgliedID}` | Ausweis eines Mitglieds abrufen | Nein |
| `POST` | `/ausweis` | Neuen Ausweis anlegen | Ja _(wenn `AUTH_ENABLED=true`)_ |
| `PUT` | `/ausweis/{id}` | Ausweis aktualisieren | Nein |
| `DELETE` | `/ausweis/{id}` | Ausweis löschen | Nein |

### Ausleihe `/ausleihe`

| Methode | Pfad | Beschreibung | Auth erforderlich |
|---|---|---|---|
| `GET` | `/ausleihe/{id}` | Ausleihe nach ID abrufen | Nein |
| `GET` | `/ausleihe/mitglied/{mitgliedID}` | Alle Ausleihen eines Mitglieds | Nein |
| `POST` | `/ausleihe` | Neue Ausleihe anlegen | Ja _(wenn `AUTH_ENABLED=true`)_ |
| `PUT` | `/ausleihe/{id}` | Ausleihe aktualisieren | Nein |
| `DELETE` | `/ausleihe/{id}` | Ausleihe löschen | Nein |

### Entwicklung `/dev`

| Methode | Pfad | Beschreibung | Auth erforderlich |
|---|---|---|---|
| `POST` | `/dev/db-populate` | Datenbank zurücksetzen & neu befüllen | Ja _(Admin-Rolle, wenn `AUTH_ENABLED=true`)_ |

---

## Datenmodell

### Mitglied

```sql
CREATE TABLE mitglied (
    id              INTEGER PRIMARY KEY,
    version         INTEGER NOT NULL DEFAULT 0,   -- Optimistisches Locking
    vorname         TEXT NOT NULL,
    nachname        TEXT NOT NULL,
    email           TEXT NOT NULL UNIQUE,
    geburtsdatum    DATE,
    telefonnummer   TEXT UNIQUE,
    geschlecht      ENUM('MAENNLICH', 'WEIBLICH', 'DIVERS'),
    mitgliedsstatus ENUM('AKTIV', 'INAKTIV'),
    beitrittsdatum  DATE,
    interessen      JSONB,
    erzeugt         TIMESTAMP NOT NULL,
    aktualisiert    TIMESTAMP NOT NULL
);
```

### Ausweis _(1:1 zu Mitglied)_

```sql
CREATE TABLE ausweis (
    id                INTEGER PRIMARY KEY,
    ausstellungsdatum DATE NOT NULL,
    ablaufdatum       DATE NOT NULL,
    mitglied_id       INTEGER NOT NULL UNIQUE REFERENCES mitglied ON DELETE CASCADE
);
```

### Ausleihe _(N:1 zu Mitglied)_

```sql
CREATE TABLE ausleihe (
    id              INTEGER PRIMARY KEY,
    ausleihdatum    DATE NOT NULL,
    rueckgabedatum  DATE NOT NULL,
    mitglied_id     INTEGER NOT NULL REFERENCES mitglied ON DELETE CASCADE
);
```

---

## Erste Schritte

### Voraussetzungen

- [Go 1.25+](https://go.dev/dl/)
- [PostgreSQL](https://www.postgresql.org/) (lokal oder via Docker)
- _(optional)_ [Docker](https://www.docker.com/)
- _(optional)_ [Keycloak](https://www.keycloak.org/) für Authentifizierung

### Konfiguration

Der Server liest seine Konfiguration **ausschließlich aus Umgebungsvariablen**:

| Variable | Standard | Beschreibung |
|---|---|---|
| `DATABASE_URL` | `postgres://mitglied:p@localhost:5432/mitglied?search_path=mitglied&sslmode=disable` | PostgreSQL-Connection-String |
| `PORT` | `8080` | HTTP-Port des Servers |
| `DB_POPULATE` | `false` | Beim Start Datenbank zurücksetzen & mit Seed-Daten befüllen |
| `AUTH_ENABLED` | `false` | OIDC-Authentifizierung via Keycloak aktivieren |

Bei aktivierter Authentifizierung werden zusätzlich folgende Variablen benötigt:

| Variable | Beschreibung |
|---|---|
| `KEYCLOAK_ISSUER` | Issuer-URL des Keycloak-Realms (z.B. `http://localhost:8880/realms/go`) |
| `KEYCLOAK_CLIENT_ID` | Client-ID (z.B. `go-client`) |
| `KEYCLOAK_AUDIENCE` | Audience (normalerweise identisch mit Client-ID) |

### Lokaler Start (ohne Docker)

```bash
# Repository klonen
git clone https://github.com/muratonuryahsi-sudo/swe-workshop.git
cd swe-workshop

# Abhängigkeiten herunterladen
go mod download

# Server starten (Standardkonfiguration)
go run ./cmd/server

# Mit expliziter Konfiguration
DATABASE_URL="postgres://mitglied:p@localhost:5432/mitglied?search_path=mitglied&sslmode=disable" \
PORT=8080 \
go run ./cmd/server
```

Der Server ist danach unter `http://localhost:8080` erreichbar.

### Mit Docker

**Standard-Image (Google Distroless):**

```bash
# Image bauen
docker build --tag murat_jeton/mitglied:latest .

# Container starten
docker run --rm \
  -e DATABASE_URL="postgres://mitglied:p@host.docker.internal:5432/mitglied?search_path=mitglied&sslmode=disable" \
  -e PORT=8080 \
  -e DB_POPULATE=false \
  -p 8080:8080 \
  murat_jeton/mitglied:latest
```

**Chainguard-Image (alternatives gehärtetes Basis-Image):**

```bash
docker build --tag murat_jeton/mitglied:chainguard -f Dockerfile.chainguard .
```

Beide Images sind **statische, shell-lose Container** ohne Package Manager für eine minimale Angriffsfläche.

---

## Authentifizierung (Keycloak / OIDC)

Standardmäßig ist die Authentifizierung **deaktiviert** – der Server startet ohne laufendes Keycloak.  
Mit `AUTH_ENABLED=true` werden `POST`-Endpunkte durch einen **Bearer-Token-Check** gesichert.

### Rollen

| Rolle | Zugriff |
|---|---|
| `user` | `POST /mitglied`, `POST /ausweis`, `POST /ausleihe` |
| `admin` | Zusätzlich `POST /dev/db-populate` |

### Keycloak-Setup

Eine vollständige Anleitung zur lokalen Keycloak-Einrichtung (Realm, Client, Testbenutzer) befindet sich in [`extras/keycloak/README.md`](extras/keycloak/README.md).

**Kurzanleitung – Token holen:**

```bash
curl -s -X POST http://localhost:8880/realms/go/protocol/openid-connect/token \
  -d grant_type=password \
  -d client_id=go-client \
  -d client_secret=changeme-go-client-secret \
  -d username=user \
  -d password=p \
  | jq -r '.access_token'
```

**Geschützten Endpunkt aufrufen:**

```bash
TOKEN=$(curl -s -X POST http://localhost:8880/realms/go/protocol/openid-connect/token \
  -d grant_type=password -d client_id=go-client \
  -d client_secret=changeme-go-client-secret \
  -d username=user -d password=p | jq -r '.access_token')

curl -X POST http://localhost:8080/mitglied \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"vorname":"Max","nachname":"Mustermann","email":"max@example.com"}'
```

---

## Datenbank-Seeding

Beim ersten Start oder zur Entwicklung kann die Datenbank automatisch mit Testdaten befüllt werden:

```bash
# Einmalig beim Start
DB_POPULATE=true go run ./cmd/server

# Oder zur Laufzeit via API (setzt DB zurück und befüllt neu)
curl -X POST http://localhost:8080/dev/db-populate
```

> ⚠️ **Achtung:** `DB_POPULATE=true` löscht alle bestehenden Daten und setzt das Schema zurück. Nur für Entwicklung/Tests verwenden.

---

## Tests

### Unit- und Build-Tests

```bash
# Formatierung prüfen
gofmt -l .

# Statische Analyse
go vet ./...

# Build prüfen
go build ./...

# Alle Tests ausführen
go test ./...
```

### Integrationstests

Integrationstests erfordern eine laufende PostgreSQL-Datenbank und optonal Keycloak.

```bash
# Ohne Auth
go test -tags integration -v ./test/integration/...

# Mit Keycloak-Authentifizierung
AUTH_ENABLED=true \
KEYCLOAK_PORT=8880 \
go test -tags integration -v ./test/integration/...

# Einzelne Testfälle
go test -tags integration -v -run TestMitgliedCreate_Unauthorized ./test/integration/...
go test -tags integration -v -run TestMitgliedCreate_WithToken ./test/integration/...
go test -tags integration -v -run TestAusweisCreate_Unauthorized ./test/integration/...
go test -tags integration -v -run TestAusleiheCreate_Unauthorized ./test/integration/...
```

### API-Tests mit Bruno

In `extras/bruno/go-appserver/` befindet sich eine vollständige **Bruno-Collection** mit vorbereiteten Requests für alle Endpunkte, einschließlich Token-Anfrage und geschützte Endpunkte.

1. [Bruno](https://www.usebruno.com/) installieren
2. Collection `extras/bruno/go-appserver/` öffnen
3. Environment konfigurieren (Host, Port)
4. Requests ausführen

---

## CI/CD

GitHub Actions führt bei jedem Push automatisch folgende Schritte aus:

| Schritt | Beschreibung |
|---|---|
| **Formatierung** | `gofmt -l .` – schlägt fehl bei unformatierten Dateien |
| **Statische Analyse** | `go vet ./...` |
| **Build** | `go build ./...` |
| **Tests** | `go test ./...` |

Die Pipeline ist in [`.github/workflows/ci.yml`](.github/workflows/ci.yml) definiert.

---

## Projektstruktur

```
swe-workshop/
├── .github/
│   └── workflows/
│       └── ci.yml                  ← GitHub Actions CI-Pipeline
├── cmd/
│   └── server/
│       └── main.go                 ← Einstiegspunkt
├── internal/
│   ├── ausleihe/
│   │   ├── model.go                ← Datenmodell & Validierung
│   │   ├── repository.go           ← Datenbankzugriff (GORM)
│   │   ├── service.go              ← Geschäftslogik
│   │   └── router.go               ← HTTP-Handler (chi)
│   ├── ausweis/
│   │   ├── model.go
│   │   ├── repository.go
│   │   ├── service.go
│   │   └── router.go
│   ├── config/
│   │   └── config.go               ← Konfiguration aus Env-Variablen
│   ├── db/
│   │   ├── db.go                   ← Datenbankverbindung
│   │   ├── populate.go             ← Schema-Reset & Seeding
│   │   ├── seed.go                 ← Seed-Logik
│   │   └── resources/
│   │       ├── create-table.sql
│   │       └── drop-table.sql
│   ├── dev/
│   │   └── router.go               ← POST /dev/db-populate
│   ├── mitglied/
│   │   ├── model.go
│   │   ├── repository.go
│   │   ├── service.go
│   │   └── router.go
│   ├── security/
│   │   ├── config.go               ← Keycloak-Konfiguration
│   │   └── middleware.go           ← OIDC Bearer-Token-Prüfung
│   └── shared/
│       ├── errors.go               ← Gemeinsame Fehlertypen
│       └── http.go                 ← JSON-Response-Hilfsfunktionen
├── test/
│   └── integration/                ← Integrationstests (Build-Tag: integration)
├── extras/
│   ├── bruno/                      ← Bruno API-Collection
│   └── keycloak/
│       ├── README.md               ← Keycloak-Setup-Anleitung
│       └── setup-realm.sh          ← Realm-Initialisierungsskript
├── Dockerfile                      ← Distroless-Image (Produktion)
├── Dockerfile.chainguard           ← Chainguard-Image (Alternative)
├── go.mod
└── go.sum
```

---

## Lizenz

Dieses Projekt steht unter der **GNU General Public License v3.0 or later**.  
Siehe [GNU GPL v3](https://www.gnu.org/licenses/gpl-3.0.html) für Details.

Copyright © 2026 – Murat Yahsi, Jeton Rama, Hochschule Karlsruhe
