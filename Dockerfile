# syntax=docker.io/docker/dockerfile-upstream:1.24.0
# check=error=true

# Copyright (C) 2026 - present, Murat Yahsi, Jeton Rama, Hochschule Karlsruhe
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program. If not, see <https://www.gnu.org/licenses/>.

# Aufruf:   docker build --tag murat_jeton/mitglied:2026.4.1-hardened .
#               ggf. --progress=plain
#               ggf. --no-cache
#
#           Windows:   Get-Content Dockerfile | docker run --rm --interactive hadolint/hadolint:v2.14.0-debian
#           macOS:     cat Dockerfile | docker run --rm --interactive hadolint/hadolint:v2.14.0-debian
#
#           docker run --rm \
#             -e DATABASE_URL=postgres://mitglied:p@host.docker.internal:5432/mitglied?search_path=mitglied&sslmode=disable \
#             -e PORT=8080 \
#             -e DB_POPULATE=false \
#             -p 8080:8080 \
#             murat_jeton/mitglied:2026.4.1-hardened
#
#           docker debug murat_jeton/mitglied:2026.4.1-hardened
#           docker network ls
#           docker save murat_jeton/mitglied:2026.4.1-hardened > mitglied.tar

# https://docs.docker.com/engine/reference/builder/#syntax
# https://github.com/moby/buildkit/blob/master/frontend/dockerfile/docs/reference.md
# https://hub.docker.com/r/docker/dockerfile

ARG GO_VERSION=1.25

# Stage 1: Go-Binary kompilieren
FROM golang:${GO_VERSION}-bookworm AS build

WORKDIR /app

# Abhaengigkeiten zuerst herunterladen (Cache-Layer)
RUN --mount=type=bind,source=go.mod,target=go.mod \
  --mount=type=bind,source=go.sum,target=go.sum \
  --mount=type=cache,target=/root/.cache/go-build \
  go mod download

# Statisches Binary bauen (kein CGO, keine glibc-Abhaengigkeit im finalen Image)
RUN --mount=type=bind,source=.,target=. \
  --mount=type=cache,target=/root/.cache/go-build \
  CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server

# Stage 2: Finales gehaertetes Image
# Distroless static: kein Shell, kein Package Manager, minimale Angriffsfläche
FROM gcr.io/distroless/static-debian12:nonroot AS final

WORKDIR /opt/app

COPY --from=build /app/server ./server

# Anzeige bei "docker inspect ..."
# https://specs.opencontainers.org/image-spec/annotations
# https://spdx.org/licenses
# MAINTAINER ist deprecated https://docs.docker.com/engine/reference/builder/#maintainer-deprecated
LABEL org.opencontainers.image.title="mitglied" \
  org.opencontainers.image.description="Appserver mit 'hardened' Basis-Image Go und Debian 12" \
  org.opencontainers.image.version="2026.4.1-bookworm" \
  org.opencontainers.image.licenses="GPL-3.0-or-later" \
  org.opencontainers.image.authors="Murat Yahsi, Jeton Rama"

# Laufzeit-Konfiguration via Umgebungsvariablen (siehe internal/config/config.go)
# Keine Defaults hier - Werte werden beim "docker run -e ..." uebergeben
ENV DATABASE_URL="" \
  PORT="" \
  DB_POPULATE=""

EXPOSE 8080

# Bei CMD statt ENTRYPOINT kann das Kommando bei "docker run ..." ueberschrieben werden
# "Array Syntax" damit auch <Strg>C funktioniert
ENTRYPOINT ["/opt/app/server"]
