package main

import (
	"flag"
	"log"
	"net/http"

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

	svr := server.NewServer(*cfg)
	log.Printf("Server listening on port: %v\n", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, svr))
}
