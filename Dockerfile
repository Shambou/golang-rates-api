FROM golang:1.18

RUN mkdir /app
COPY . /app
WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app cmd/api/main.go
CMD ["./app"]
