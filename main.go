package main

import (
	"log"

	"github.com/ultimateboy/discapi/api"
	"github.com/ultimateboy/discapi/config"
)

func main() {
	log.Println("Starting Disc API...")

	cfg, err := config.ParseConfig()
	if err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}

	err = cfg.Log()
	if err != nil {
		log.Fatalf("Failed to log config: %v", err)
	}

	api, err := api.NewAPI(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize the API: %v", err)
	}

	err = api.Start()
	if err != nil {
		log.Fatalf("Failed to start the API: %v", err)
	}
}
