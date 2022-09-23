FROM golang:latest

COPY . ./src
WORKDIR ./src

RUN go mod download

EXPOSE 8080

CMD [ "go", "run", "main.go" ]
