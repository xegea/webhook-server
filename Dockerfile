# syntax=docker/dockerfile:1

FROM golang:1.19 as build

WORKDIR /usr/src/webhook-server

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/webhook-server/ 


FROM alpine:latest as prod
RUN apk --no-cache add ca-certificates
WORKDIR /
COPY --from=build /usr/local/bin /app
COPY --from=build /usr/src/webhook-server/.env /app/

EXPOSE 8080

CMD [ "/app/webhook_server", "-env=/app/.env" ]