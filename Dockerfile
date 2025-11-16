FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

RUN go build -o service ./cmd/main.go


FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/service /app/service
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY migrations ./migrations

EXPOSE 8080

CMD sh -c '\
  echo "waiting for postgres..."; \
  until goose -dir ./migrations status >/dev/null 2>&1; do \
    echo "postgres is not ready yet, retrying..."; \
    sleep 2; \
  done; \
  echo "postgres is up, running migrations..."; \
  goose -dir ./migrations up && \
  ./service \
'
