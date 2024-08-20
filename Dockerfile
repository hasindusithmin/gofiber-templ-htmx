FROM golang:1.22-alpine AS build-stage

WORKDIR /app

COPY . .

RUN go mod download && go mod verify

RUN go install github.com/a-h/templ/cmd/templ@latest

RUN templ generate

RUN go build -ldflags="-s -w" -o ./bin/main .

CMD ["./bin/main"]