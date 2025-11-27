# Build stage
FROM golang:1.24-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux GOARCH=amd64 go build -o /app/gateway-app ./gateway/cmd/main/main.go

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=build /app/gateway-app ./gateway-app
COPY gateway/static ./static
EXPOSE 9090
CMD ["./gateway-app"]
