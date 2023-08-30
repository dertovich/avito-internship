FROM golang:1.19 AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o /app/segment-service ./cmd/segment-service

FROM alpine:latest

COPY --from=builder /app/segment-service /app/segment-service

CMD ["/app/segment-service"]
