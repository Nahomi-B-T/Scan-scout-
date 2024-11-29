# syntax=docker/dockerfile:experimental

# Versión de Go y Alpine configuradas como variables de entorno
ARG GO_VERSION=1.20
ARG ALPINE_VERSION=3.18

#####
# BUILDER - Compila el plugin
#####
FROM golang:${GO_VERSION} AS builder
WORKDIR /app

# Copiar los archivos del proyecto
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/root/.cache/go-build \
    go mod download

COPY . .

# Compilar el binario
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o docker-scout-plugin .

#####
# RUNNER - Ejecuta el plugin
#####
FROM alpine:${ALPINE_VERSION} AS runner

# Instalar Docker CLI
RUN apk add --no-cache docker-cli

# Copiar el binario compilado
COPY --from=builder /app/docker-scout-plugin /usr/local/bin/docker-scout-plugin

# Permitir que el contenedor acceda al socket Docker
VOLUME /var/run/docker.sock

# Configurar el comando por defecto
ENTRYPOINT ["docker-scout-plugin"]