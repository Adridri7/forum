FROM golang:1.23-alpine

# Installer les dépendances nécessaires pour CGO et sqlite3
RUN apk update && apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

COPY go.mod ./
COPY . .

# Construire l'application avec CGO activé
RUN CGO_ENABLED=1 GOOS=linux go build -o main

# Nettoyer les paquets pour réduire la taille de l'image
RUN apk del gcc musl-dev

EXPOSE 8080

CMD ["./main"]