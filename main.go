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

	s := []stores.Store{
		stores.Argos{},
		stores.Amazon{},
		stores.Currys{},
		// TODO Re-enable Game once add timeout
		// stores.Game{},
		// TODO John Lewis doesn't work in headless mode
		// stores.JohnLewis{},
		stores.ShopTo{},
		stores.Smyths{},
	}

	results := services.CheckStores(s)

	hasStock := false

	for _, v := range results {
		if v.Status == stores.InStock {
			hasStock = true
			break
		}
	}

	if hasStock {
		services.Notify(results)
	} else {
		util.Logger.Println(results)
	}

	util.Logger.Println("Task complete")
}
