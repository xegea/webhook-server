# syntax=docker/dockerfile:1

FROM golang:1.19-alpine as build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /webhook_server

# FROM gcr.io/distroless/base-debian10

# WORKDIR /

# COPY --from=build /webhook_server /webhook_server

EXPOSE 8080

CMD [ "/webhook_server" ]