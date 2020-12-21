package main

import (
	"time"

	"github.com/jonjam/stock-checker/config"
	"github.com/jonjam/stock-checker/services"
	"github.com/jonjam/stock-checker/stores"
	"github.com/jonjam/stock-checker/util"

	"github.com/go-co-op/gocron"
)

func main() {
	s := gocron.NewScheduler(time.UTC)

	i := config.GetSchedulerConfig().Interval
	_, err := s.Every(i).Hour().Do(task)

	if err != nil {
		util.Logger.Fatalln(err)
	}

	s.StartBlocking()
}

func task() {
	util.Logger.Println("Starting task")

	stores := []stores.Store{
		// stores.Argos{},
		// stores.Amazon{},
		// stores.Currys{},
		// stores.Game{},
		// TODO workout why this fails in headless mode
		stores.JohnLewis{},
		// stores.ShopTo{},
		// stores.Smyths{},
	}

	results := services.CheckStores(stores)

	// TODO Only notify if one is true
	services.Notify(results)

	util.Logger.Println("Task complete")
}
