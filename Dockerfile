FROM golang:1.21-bullseye AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	go build -ldflags='-s -w' -o /app/transcription-service ./cmd

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
WORKDIR /app

COPY --from=builder /app/transcription-service /usr/local/bin/transcription-service

RUN addgroup -S app && adduser -S -G app app
USER app

ENTRYPOINT ["/usr/local/bin/transcription-service"]