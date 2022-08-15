FROM golang:1.18-stretch

WORKDIR /src

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY app/go.* ./

RUN go mod download

COPY migrations ./migrations
COPY app/templates .templates

COPY app ./app

WORKDIR ./app

RUN go build -v -o server cmd/server/main.go

CMD [ "./server" ]