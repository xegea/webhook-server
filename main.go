package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/xegea/webhook_server/pkg/config"
	"github.com/xegea/webhook_server/pkg/server"
)

func main() {

	env := flag.String("env", ".env", ".env path")
	flag.Parse()

	cfg, err := config.LoadConfig(env)
	if err != nil {
		log.Fatalf("unable to load config: %+v", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisHost,
		Password: cfg.RedisPassword,
		DB:       0,
	})

	pong, err := client.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Redis PING: %v\n", pong)

	svr := server.NewServer(
		*cfg,
		*client,
	)

	fmt.Printf("Server listening on port: %v\n", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, svr))
}
