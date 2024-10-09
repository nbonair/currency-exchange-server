package main

import (
	"fmt"
	"log"

	"github.com/nbonair/currency-exchange-server/configs"
	"github.com/nbonair/currency-exchange-server/internal/www"
)

func main() {

	// Load configuration
	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	if err := www.Run(cfg); err != nil {
		fmt.Println(err)
	}

}
