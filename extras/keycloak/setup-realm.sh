#!/usr/bin/env bash
# Konfiguriert den Keycloak-Realm "go" via REST-API.
# Voraussetzung: Keycloak laeuft im start-dev-Modus auf http://localhost:8080
# Aufruf: extras/keycloak/setup-realm.sh

set -euo pipefail

KC_URL="${KEYCLOAK_URL:-http://localhost:8080}"
CLIENT_SECRET="${CLIENT_SECRET:-changeme-go-client-secret}"
KC_ADMIN_USERNAME="${KC_ADMIN_USERNAME:-admin}"
KC_ADMIN_PASSWORD="${KC_ADMIN_PASSWORD:-admin}"

echo "Warte auf Keycloak unter $KC_URL ..."
until curl -fs "$KC_URL/realms/master" > /dev/null 2>&1; do
  sleep 3
done
echo "Keycloak ist bereit."

# Admin-Token vom Master-Realm holen
TOKEN=$(curl -s -X POST "$KC_URL/realms/master/protocol/openid-connect/token" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "client_id=admin-cli&username=$KC_ADMIN_USERNAME&password=$KC_ADMIN_PASSWORD&grant_type=password" | jq -r '.access_token')

if [ -z "$TOKEN" ] || [ "$TOKEN" = "null" ]; then
  echo "Fehler: Konnte keinen Admin-Token holen." >&2
  exit 1
fi

# Realm "go" anlegen
echo "Lege Realm 'go' an ..."
curl -s -o /dev/null -w "%{http_code}" -X POST "$KC_URL/admin/realms" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"realm":"go","enabled":true,"accessTokenLifespan":1800,"ssoSessionIdleTimeout":3600}'

# Client anlegen
echo "Lege Client 'go-client' an ..."
curl -s -o /dev/null -X POST "$KC_URL/admin/realms/go/clients" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"clientId\": \"go-client\",
    \"name\": \"Go Client\",
    \"enabled\": true,
    \"publicClient\": false,
    \"clientAuthenticatorType\": \"client-secret\",
    \"secret\": \"$CLIENT_SECRET\",
    \"directAccessGrantsEnabled\": true,
    \"serviceAccountsEnabled\": true,
    \"standardFlowEnabled\": true,
    \"rootUrl\": \"http://localhost:8080\",
    \"redirectUris\": [\"*\"],
    \"webOrigins\": [\"+\"]
  }"

# Client-UUID ermitteln
CLIENT_UUID=$(curl -s "$KC_URL/admin/realms/go/clients?clientId=go-client" \
  -H "Authorization: Bearer $TOKEN" | jq -r '.[0].id')
echo "Client-UUID: $CLIENT_UUID"

# Rollen anlegen
for ROLE in admin user; do
  echo "Lege Rolle '$ROLE' an ..."
  curl -s -o /dev/null -X POST "$KC_URL/admin/realms/go/clients/$CLIENT_UUID/roles" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "{\"name\":\"$ROLE\"}"
done

# Hilfsfunktion: User anlegen und Rolle zuweisen
create_user() {
  local USERNAME="$1"
  local EMAIL="$2"
  local FIRSTNAME="$3"
  local LASTNAME="$4"
  local ROLE="$5"

  echo "Lege User '$USERNAME' an ..."
  curl -s -o /dev/null -X POST "$KC_URL/admin/realms/go/users" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "{
      \"username\": \"$USERNAME\",
      \"email\": \"$EMAIL\",
      \"firstName\": \"$FIRSTNAME\",
      \"lastName\": \"$LASTNAME\",
      \"enabled\": true,
      \"emailVerified\": true,
      \"credentials\": [{\"type\":\"password\",\"value\":\"p\",\"temporary\":false}]
    }"

  local USER_ID
  USER_ID=$(curl -s "$KC_URL/admin/realms/go/users?username=$USERNAME" \
    -H "Authorization: Bearer $TOKEN" | jq -r '.[0].id')

  local ROLE_JSON
  ROLE_JSON=$(curl -s "$KC_URL/admin/realms/go/clients/$CLIENT_UUID/roles/$ROLE" \
    -H "Authorization: Bearer $TOKEN")

  echo "Weise Rolle '$ROLE' dem User '$USERNAME' zu ..."
  curl -s -o /dev/null -X POST \
    "$KC_URL/admin/realms/go/users/$USER_ID/role-mappings/clients/$CLIENT_UUID" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "[$ROLE_JSON]"
}

create_user "admin" "admin@acme.com" "Go" "Admin" "admin"
create_user "user"  "user@acme.com"  "Go" "User"  "user"

echo "Keycloak-Realm 'go' vollstaendig konfiguriert."
