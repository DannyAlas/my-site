# Build stage
FROM golang:1.22.2-alpine3.18 AS build
WORKDIR /app
COPY . .
RUN go build -o ./cmd/build/build.go ./cmd/build
RUN go build -o ./cmd/serve/serve.go ./cmd/serve

# Serve
FROM alpine:3.14
WORKDIR /app
COPY --from=build /app/cmd/build/build.go /app/cmd/build/build.go
COPY --from=build /app/cmd/serve/serve.go /app/cmd/serve/serve.go
COPY --from=build /app/cmd/serve/serve /app/cmd/serve/serve
COPY --from=build /app/cmd/build/build /app/cmd/build/build
CMD ["/app/cmd/serve/serve"]
