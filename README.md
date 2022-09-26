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


#To run the Dockerfile
docker image build --tag webhook_server-api .

#Execute docker compose
docker compose up 

#Interactive session to webhook_server-api
docker container run -i -t --rm  webhook_server-api sh

#Redis interactive session to webhook_server-redis-1
winpty docker exec -it webhook_server-redis-1 redis-cli -a <password>

#or connect to Redis using
redis-cli -h 127.0.0.1 -p 6379 -a <password>

#If connection to Redis does not work:
#from powershell:
>docker container exec -it webhook-redis sh
>ifconfig # retrieve ip
>redis-cli -h 172.22.0.2 -p 6379 -a eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81

# And then...it works
172.22.0.2:6379> keys *
1) "d27b0ad4-5c01-410d-bc66-463c752e1a5b"
2) "07ca572a-53c8-464f-9064-1fee974a7759"
