FROM golang:latest as bd

WORKDIR /app

COPY . .

RUN go build -o ./server ./internal/basescript/saveInBase.go


FROM ubuntu:22.04

WORKDIR /app

COPY --from=bd /app/server ./server

# RUN chmod +x ./server

CMD  ["/app/server"]

