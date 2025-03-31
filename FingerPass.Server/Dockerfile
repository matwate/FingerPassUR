FROM golang:1.24.1-alpine as base

# Development stage
FROM base as development
WORKDIR /app
RUN go install github.com/air-verse/air@latest
COPY go.mod go.sum ./
RUN go mod download
CMD ["air", "--build.cmd", "go build -o ./bin/api ./cmd/api", "--build.bin", "./bin/api"]

# Build stage
FROM base AS build
WORKDIR /build
COPY go.mod go.sum ./
RUN CGO_ENABLED=0 go build ./cmd/api -o api

#Production stage
FROM scratch AS production
WORKDIR /prod
COPY --from=build api ./
EXPOSE 8080
CMD ["/prod/api"]
