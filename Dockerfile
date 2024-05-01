FROM golang:latest AS builder

RUN apt-get update && apt-get install -y git

WORKDIR /app/build

RUN git clone https://github.com/DannyAlas/my-site .

WORKDIR /app/build/cmd/build

RUN go mod download

RUN go build -o /app/build/build ./build.go

WORKDIR /app/build/cmd/serve

RUN go mod download

RUN go build -o /app/build/serve ./webserver.go

FROM golang:latest

COPY --from=builder /app/build /app/build

EXPOSE 8080

CMD ["/app/build/webserver"]

