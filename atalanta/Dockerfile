# Build stage
FROM golang:1.23 AS builder
LABEL authors="mohammadhoseinaref"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o atalanta ./cmd

#Run stage
FROM ubuntu:latest

WORKDIR /root/

COPY --from=builder /app/atalanta .

CMD ["./atalanta"]
