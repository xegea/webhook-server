package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/xegea/webhook_server/pkg/config"
	"github.com/xegea/webhook_server/pkg/server"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("unable to load config: %+v", err)
	}

	svr := server.NewServer(*cfg)
	fmt.Printf("Server listening on port: %v\n", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, svr))
}
