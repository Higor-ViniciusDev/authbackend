FROM golang:1.24-alpine

WORKDIR /appauth

# Install dependencies needed for CGO if necessary (though go-pg usually doesn't need it)
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /usr/local/bin/auth-backend ./cmd/api

EXPOSE 8080

CMD ["auth-backend"]