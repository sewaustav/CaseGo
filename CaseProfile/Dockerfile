# syntax=docker/dockerfile:1
FROM golang:1.26.1-alpine AS builder
WORKDIR /build

RUN apk add --no-cache git

ARG GITHUB_TOKEN
COPY go.mod go.sum ./
RUN git config --global url."https://${GITHUB_TOKEN}@github.com/".insteadOf "https://github.com/" && \
    go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /app/service \
    ./cmd/case_profile

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app
COPY --from=builder /app/service /app/service
COPY --from=builder /build/migrations /app/migrations

EXPOSE 8080

ENTRYPOINT ["/app/service"]
