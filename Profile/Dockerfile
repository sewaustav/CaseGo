FROM golang:alpine AS builder
WORKDIR /app

RUN apk add --no-cache git

ARG SERVICE_NAME
ARG GITHUB_TOKEN
COPY go.mod go.sum ./
RUN git config --global url."https://${GITHUB_TOKEN}@github.com/".insteadOf "https://github.com/" && \
    go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/${SERVICE_NAME}/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations
CMD ["./main"]