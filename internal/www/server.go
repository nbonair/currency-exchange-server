package www

import (
	"fmt"
	"log"

	"github.com/nbonair/currency-exchange-server/configs"
	"github.com/nbonair/currency-exchange-server/internal/dataaccess/db"
	"github.com/nbonair/currency-exchange-server/internal/handler"
	"github.com/nbonair/currency-exchange-server/internal/lib/openexchangerates"
	"github.com/nbonair/currency-exchange-server/internal/repo"
	"github.com/nbonair/currency-exchange-server/internal/router"
	"github.com/nbonair/currency-exchange-server/internal/service"
)

func Run(cfg *configs.Config) error {
	pool, cleanup, err := db.InitializeDB(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer cleanup()

	exchangeRateRepo := repo.NewExchangeRateRepository(pool.Pool)
	appID := cfg.APIs.APIKeys["openexchangerates"]

	apiClient, err := openexchangerates.NewOpenExchangeRateClient(appID)
	if err != nil {
		fmt.Println("Get External API error")
	}

	exchangeRateService := service.NewExchangeRateService(exchangeRateRepo, apiClient)
	exchangeRateHandler := handler.NewExchangeRateHandler(exchangeRateService)

	r := router.SetupRouter(exchangeRateHandler)

	return r.Run(":8080")
}
