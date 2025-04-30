# Etapa 1: Build do binário
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copiar arquivos de dependência e baixar módulos
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copiar o restante do código
COPY . .

# Build do binário para Linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app

# Etapa 2: Imagem final mínima
FROM alpine:latest

# ✅ inclui timezones para evitar panic no .In()
RUN apk add --no-cache tzdata

WORKDIR /root/

# Copiar apenas o binário da etapa de build
COPY --from=builder /app/app .

# Expor a porta que seu app usa
EXPOSE 8080

# Comando para rodar o app
CMD ["./app"]
