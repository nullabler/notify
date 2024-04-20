FROM golang:1.22 AS builder

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /usr/local/bin/gateway ./cmd/gateway/gateway.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /usr/local/bin/gateway .
CMD ["/app/gateway"]
