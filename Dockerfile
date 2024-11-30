# syntax=docker/dockerfile:experimental

# Versión de Go y Alpine configuradas como variables de entorno
ARG GO_VERSION=1.21.12        # Actualizada a una versión más reciente
ARG ALPINE_VERSION=3.20      # Versión de Alpine actualizada

#####
# BUILDER - Compila el plugin
#####
FROM golang:${GO_VERSION}-alpine AS builder

# Establecer directorio de trabajo
WORKDIR /app

# Instalar dependencias necesarias para la compilación
RUN apk add --no-cache git

# Copiar los archivos de configuración del proyecto
COPY go.mod go.sum ./

# Descargar las dependencias (esto puede aprovechar la caché de Go)
RUN --mount=type=cache,target=/root/.cache/go-build \
    go mod download

# Copiar todo el código fuente
COPY . .

# Compilar el binario con optimización (sin dependencias de C)
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o docker-scout-plugin .

#####
# RUNNER - Ejecuta el plugin
#####
FROM alpine:${ALPINE_VERSION} AS runner

# Instalar Docker CLI en la imagen final
RUN apk add --no-cache docker-cli

# Copiar el binario compilado desde la fase builder
COPY --from=builder /app/docker-scout-plugin /usr/local/bin/docker-scout-plugin

# Permitir que el contenedor acceda al socket Docker
VOLUME /var/run/docker.sock

# Configurar el comando por defecto
ENTRYPOINT ["docker-scout-plugin"]
