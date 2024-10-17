package www

import (
	"fmt"
	"log"

	"github.com/nbonair/currency-exchange-server/configs"
	"github.com/nbonair/currency-exchange-server/internal/wiring"
)

func Run() error {
	// Load configuration
	cfg, err := configs.LoadConfig()
	if err != nil {
		fmt.Println("Failed to load configuration: %v", err)
	}

	r, cleanup, err := wiring.InitializeRouter(cfg.Database, cfg.APIs, cfg.Redis)

	if err != nil {
		log.Fatalf("Failed to initialize router: %v", err)
	}

	defer cleanup()
	return r.Run(":8080")
}
