FROM golang:latest as build

WORKDIR /app

COPY * .

RUN go mod download

RUN go build -o /server ./cmd/app/main.go

FROM ubuntu:20.04

RUN apt-get update
RUN apt-get install -y curl

COPY --from=build /app/server .
#COPY ./entrypoint.sh /usr/bin/entrypoint.sh
#ENTRYPOINT [ "entrypoint.sh" ]

CMD ["/server"]