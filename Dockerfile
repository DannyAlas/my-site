FROM golang:latest

WORKDIR /app

RUN git clone https://github.com/DannyAlas/my-site .

RUN go build -o webserver ./cmd/server/webserver

RUN go build -o build ./cmd/build/build.go

EXPOSE 8080

CMD ["./webserver"]
