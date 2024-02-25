FROM golang:latest as build

WORKDIR /application

COPY . .

RUN go mod download

RUN go build -o ./server ./cmd/app/main.go

FROM ubuntu:22.04

RUN apt update
# RUN apt-get install -y curl
# RUN apt install libc6 

WORKDIR /app

COPY --from=build /application/db . 
COPY --from=build /application/server ./server
#COPY ./entrypoint.sh /usr/bin/entrypoint.sh
#ENTRYPOINT [ "entrypoint.sh" ]

CMD ["/app/server"]