FROM golang:1.24 AS builder

WORKDIR /appauth

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o auth-backend ./cmd/api

FROM alpine:3.21

RUN apk --no-cache add ca-certificates

WORKDIR /appauth

COPY --from=builder /appauth/auth-backend .

EXPOSE 8080

CMD ["./auth-backend"]
