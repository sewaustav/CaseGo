FROM golang:1.26-alpine AS builder
WORKDIR /app

ARG GITHUB_TOKEN
ENV GOPRIVATE=github.com/sewaustav/*
RUN apk add --no-cache git && \
git config --global url."https://${GITHUB_TOKEN}@github.com/".insteadOf "https://github.com/" 

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /payment ./cmd/...

FROM alpine:latest
WORKDIR /app
COPY --from=builder /payment .
COPY --from=builder /app/migrations ./migrations
EXPOSE 8085 50052
CMD ["/app/payment"]