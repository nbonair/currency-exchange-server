package www

import (
	"fmt"

	"github.com/nbonair/currency-exchange-server/configs"
	"github.com/nbonair/currency-exchange-server/internal/wiring"
)

func Run() error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic in Run: %v", r)
		}
	}()

	// Load configuration
	cfg, err := configs.LoadConfig()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v", err)
	}

	r, cleanup, err := wiring.InitializeRouter(cfg.Database, cfg.APIs, cfg.Redis)

	if err != nil {
		return fmt.Errorf("failed to initialize router: %w", err)
	}

	defer cleanup()

	if err := r.Run(":8080"); err != nil {
		return fmt.Errorf("server encountered an error: %w", err)
	}
	return nil
}
