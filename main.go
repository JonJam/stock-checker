package main

import (
	"time"

	"github.com/jonjam/stock-checker/services"
	"github.com/jonjam/stock-checker/stores"
	"github.com/jonjam/stock-checker/util"

	"github.com/go-co-op/gocron"
)

func main() {
	s := gocron.NewScheduler(time.UTC)

	// TODO Configure via config
	_, err := s.Every(1).Hour().Do(task)

	if err != nil {
		util.Logger.Fatalln(err)
	}

	s.StartBlocking()
}

func task() {
	util.Logger.Println("Starting task")

	stores := []stores.Store{
		stores.Argos{},
		stores.Amazon{},
		stores.Currys{},
		stores.Game{},
		stores.JohnLewis{},
		stores.ShopTo{},
		stores.Smyths{},
	}

	results := services.CheckStores(stores)

	services.Notify(results)

	util.Logger.Println("Task complete")
}
