FROM golang:latest as build

WORKDIR /app

COPY . .

RUN go build -o /server /app/internal/basescript/saveInBase.go

FROM ubuntu:22.04 as prod

WORKDIR /app

COPY --from=build /server .

CMD  ["/app/server"]