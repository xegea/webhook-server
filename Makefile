all: build run 

build: 
	go build -o bin/webserver ./cmd 

test:
	go test

run:
	bin/webserver

up:
	docker compose up -d

cover:
	go test -cover ./...

html-cover:
	go get golang.org/x/tools/cmd/cover 
	go test -coverprofile cover.out 
	go tool cover -html=cover.out

