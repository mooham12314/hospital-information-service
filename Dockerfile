FROM golang:1.25.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/hospital-service ./cmd/server

FROM alpine:3.21
RUN adduser -D -g '' appuser
WORKDIR /app
COPY --from=builder /bin/hospital-service /app/hospital-service
USER appuser
EXPOSE 8080
ENTRYPOINT ["/app/hospital-service"]
