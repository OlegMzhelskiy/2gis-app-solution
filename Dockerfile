FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server ./cmd/server/main.go

FROM debian:bookworm-slim

COPY --from=builder /app/server /server
#COPY config.yaml /config.yaml

EXPOSE 8080

CMD ["/server"]