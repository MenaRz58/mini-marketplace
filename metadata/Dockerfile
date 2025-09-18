# Build stage
FROM golang:1.24-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux GOARCH=amd64 go build -o /app/metadata-app ./metadata/cmd/main/main.go

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=build /app/metadata-app ./metadata-app
CMD ["./metadata-app"]
