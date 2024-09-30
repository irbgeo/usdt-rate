FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN  go install ./cmd/usdt-rate

FROM alpine:3.14

WORKDIR /root/

COPY --from=builder /go/bin/usdt-rate .
COPY --from=builder /app/migrations ./migrations

CMD ["./usdt-rate"]