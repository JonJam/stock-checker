package main

import (
	"time"

	"github.com/jonjam/stock-checker/config"
	"github.com/jonjam/stock-checker/services"
	"github.com/jonjam/stock-checker/stores"

	"github.com/go-co-op/gocron"

	"go.uber.org/zap"
)

func main() {
	var logger *zap.Logger

	if config.GetLogConfig().Development {
		logger, _ = zap.NewDevelopment()
	} else {
		logger, _ = zap.NewProduction()
	}

	defer logger.Sync()

	s := gocron.NewScheduler(time.UTC)

	i := config.GetSchedulerConfig().Interval
	_, err := s.Every(i).Hour().Do(func() {
		checkStores(logger)
	})

	if err != nil {
		logger.Fatal("Failed to create job.", zap.Error(err))
	}

	s.StartBlocking()
}

func checkStores(logger *zap.Logger) {
	logger.Info("Starting task.")

	s := []stores.Store{
		stores.NewArgos(logger),
		stores.NewAmazon(logger),
		stores.NewCurrys(logger),
		// TODO Re-enable Game once add timeout
		// stores.NewGame(logger),
		// TODO John Lewis doesn't work in headless mode
		// stores.NewJohnLewis(logger),
		stores.NewShopTo(logger),
		stores.NewSmyths(logger),
	}

	c := services.NewStoresChecker(logger)
	results := c.CheckStores(s)

	hasStock := false

	for _, v := range results {
		if v.Status == stores.InStock {
			hasStock = true
			break
		}
	}

	if hasStock {
		n := services.NewNotifier(logger)
		n.Notify(results)
	}

	logger.Info("Task complete.", zap.Any("results", results))
}
