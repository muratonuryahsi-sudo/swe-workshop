# Keycloak – Lokales Setup

## Voraussetzungen

- Docker läuft lokal
- Repository `hono-appserver` liegt unter `~/vsprojects/hono-appserver`

---

## Schritt 1 – Keycloak starten

```bash
cd ~/vsprojects/hono-appserver/extras/compose/keycloak
docker compose up -d
```

Keycloak läuft danach auf **http://localhost:8880**.
Warten bis der Container healthy ist (~30 Sekunden).

Status prüfen:

```bash
docker compose ps
```

---

## Schritt 2 – Realm "go" einrichten

Einmalig ausführen (legt Realm, Client und Testuser an):

```bash
cd ~/vsprojects/swe-workshop
KEYCLOAK_URL=http://localhost:8880 \
KC_ADMIN_USERNAME=tmp \
KC_ADMIN_PASSWORD=p \
bash extras/keycloak/setup-realm.sh
```

Angelegt werden:

| Was | Wert |
|---|---|
| Realm | `go` |
| Client | `go-client` |
| Client Secret | `changeme-go-client-secret` |
| User admin | `admin` / `p` (Rolle: `admin`) |
| User user | `user` / `p` (Rolle: `user`) |

---

## Schritt 3 – Go-Server mit Auth starten

```bash
cd ~/vsprojects/swe-workshop
AUTH_ENABLED=true go run ./cmd/server
```

Ohne `AUTH_ENABLED=true` läuft der Server ohne Auth (POST offen).

---

## Schritt 4 – Token holen (curl)

```bash
curl -s -X POST http://localhost:8880/realms/go/protocol/openid-connect/token \
  -d grant_type=password \
  -d client_id=go-client \
  -d client_secret=changeme-go-client-secret \
  -d username=user \
  -d password=p \
  | jq -r '.access_token'
```

---

## Schritt 5 – Bruno

### Token-Request anlegen

- **Method:** `POST`
- **URL:** `http://localhost:8880/realms/go/protocol/openid-connect/token`
- **Body:** Form (nicht JSON)

| Key | Value |
|---|---|
| `grant_type` | `password` |
| `client_id` | `go-client` |
| `client_secret` | `changeme-go-client-secret` |
| `username` | `user` |
| `password` | `p` |

`access_token` aus der Response kopieren.

### Geschützten Request senden

Bei `POST /ausweis/`, `POST /ausleihe/` oder `POST /mitglied/`:

- **Auth Tab** → `Bearer Token`
- Token einfügen

Ohne Token → `401 Unauthorized`
Mit gültigem Token → `201 Created`

---

## Schritt 6 – Integrationstests lokal ausführen

Keycloak (Schritt 1+2) und PostgreSQL müssen laufen.

```bash
cd ~/vsprojects/swe-workshop
AUTH_ENABLED=true \
KEYCLOAK_PORT=8880 \
go test -tags integration -v ./test/integration/...
```

Auth-Tests einzeln:

```bash
go test -tags integration -v -run TestMitgliedCreate_Unauthorized ./test/integration/...
go test -tags integration -v -run TestMitgliedCreate_WithToken ./test/integration/...
go test -tags integration -v -run TestAusweisCreate_Unauthorized ./test/integration/...
go test -tags integration -v -run TestAusweisCreate_WithToken ./test/integration/...
go test -tags integration -v -run TestAusleiheCreate_Unauthorized ./test/integration/...
go test -tags integration -v -run TestAusleiheCreate_WithToken ./test/integration/...
```

---

## Keycloak Admin-UI

http://localhost:8880 → Login: `tmp` / `p`

Realm `go` → Clients → `go-client` → Roles: `admin`, `user`

---

## Keycloak stoppen

```bash
cd ~/vsprojects/hono-appserver/extras/compose/keycloak
docker compose down
```
