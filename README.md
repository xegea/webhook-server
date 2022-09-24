webhook server

https://hub.docker.com/_/redis/tags?page=1&name=alpine

Download docker alpine image - 28.47MB uncompressed
docker pull redis:7.0.5-alpine

Create container and remove (-rm) when close
docker run --rm --name test-redis redis:7.0.5-alpine redis-server --loglevel warning

the same as above but execute as detached (--detach or -d)
docker run -d --rm --name test-redis redis:7.0.5-alpine redis-server --loglevel warning

run in Powershell (https://stackoverflow.com/a/45869400/2147883):
docker exec -it test-redis redis-cli

or run if I'm in Windows (https://stackoverflow.com/a/50483923/2147883):
winpty docker exec -it test-redis redis-cli


To run the Dockerfile
docker build --tag webhook_server-api .

