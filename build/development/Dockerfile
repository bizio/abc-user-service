# ---- Build Stage ----
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o bin/abc-user-service ./cmd/server/main.go

# ---- Final Stage ----
FROM alpine:3.22.2 AS final

COPY --from=builder /app/bin/abc-user-service /bin/abc-user-service
ENTRYPOINT ["bin/abc-user-service"]
