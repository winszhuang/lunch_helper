FROM golang:alpine
USER root
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o app .
CMD go run main.go --APP_ENV=prod