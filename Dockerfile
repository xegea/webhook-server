# syntax=docker/dockerfile:1

FROM golang:1.19 as build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /webhook_server

FROM alpine as prod

WORKDIR /

COPY --from=build webhook_server /app/webhook_server
COPY --from=build /app/.env /app/

EXPOSE 8080

CMD [ "/app/webhook_server", "-env=/app/.env" ]