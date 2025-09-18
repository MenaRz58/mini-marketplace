# Build stage
FROM golang:1.24-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux GOARCH=amd64 go build -o /app/orders-app ./orders/cmd/main/main.go

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=build /app/orders-app ./orders-app
CMD ["./orders-app"]
