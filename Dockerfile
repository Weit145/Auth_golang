# Build stage
FROM golang:1.25.5-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
RUN go mod tidy

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /auth_service cmd/main.go

# Final stage
FROM alpine:latest

WORKDIR /

COPY --from=builder /auth_service /auth_service
COPY config/local.yaml /config/local.yaml

EXPOSE 50051

ENTRYPOINT ["/auth_service"]
