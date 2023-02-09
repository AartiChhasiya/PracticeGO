FROM golang:1.19.0

WORKDIR /Trivia

RUN go install github.com/cosmtrek/air@latest

COPY . .

RUN go mod tidy
