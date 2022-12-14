all: build run 

build: 
	go build -o bin/webserver

test:
	go test

run:
	go build -o bin/webserver && ./bin/webserver 

up:
	docker compose up --build -d

down:
	docker compose down

image:
	docker image build -t xavieregea/webhook-server:1.0.2 .

container:
	docker container run --name webhook-server xavieregea/webhook-server:1.0.2

cover:
	go test -cover ./...

html-cover:
	go get golang.org/x/tools/cmd/cover 
	go test -coverprofile cover.out 
	go tool cover -html=cover.out

