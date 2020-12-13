package services

import (
	"log"
	"time"

	"github.com/go-rod/bypass"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/jonjam/stock-checker/stores"
)

func CheckStores(storesSlice []stores.Store) []stores.StockCheckResult {
	// TODO this should be controlled via config
	devMode := true

	url, err := createControlURL(devMode)

	if err != nil {
		log.Fatalln(err)
	}

	browser, err := createBrowser(url, devMode)

	if err != nil {
		log.Fatalln(err)
	}

	defer func() {
		err := browser.Close()

		if err != nil {
			log.Println(err)
		}
	}()

	// TODO This should be controlled via config. Dev: 1, Prod:3
	pool := rod.NewPagePool(3)
	defer pool.Cleanup(func(p *rod.Page) {
		err := p.Close()

		if err != nil {
			log.Println(err)
		}
	})

	create := createCreatePageFunc(browser, devMode)

	c := make(chan stores.StockCheckResult)

	for _, s := range storesSlice {
		go func(store stores.Store) {
			c <- store.Check(pool, create)
		}(s)
	}

	numOfStores := len(storesSlice)
	results := make([]stores.StockCheckResult, 0, numOfStores)

	for i := 0; i < numOfStores; i++ {
		results = append(results, <-c)
	}

	return results
}

func createControlURL(devMode bool) (string, error) {
	launcher := launcher.New()

	if devMode {
		launcher.Devtools(true)
		launcher.Headless(false)
	}

	return launcher.Launch()
}

func createBrowser(url string, devMode bool) (*rod.Browser, error) {
	browser := rod.New().ControlURL(url)

	if devMode {
		browser.Trace(true)
		browser.SlowMotion(time.Second)
	}

	err := browser.Connect()

	if err != nil {
		return nil, err
	}

	return browser, nil
}

func createCreatePageFunc(browser *rod.Browser, devMode bool) func() *rod.Page {
	// This func will create a new configured page will be contained within a different incognito browser window.
	// It returns nil when an error occurs rather than exposing error due to https://pkg.go.dev/github.com/go-rod/rod#PagePool.Get
	return func() *rod.Page {
		browser, err := browser.Incognito()

		if err != nil {
			log.Println(err)

			return nil
		}

		page, err := bypass.Page(browser)

		if err != nil {
			log.Println(err)

			return nil
		}

		if devMode {
			err = page.SetWindow(&proto.BrowserBounds{
				WindowState: proto.BrowserWindowStateFullscreen,
			})

			if err != nil {
				log.Println(err)

				return nil
			}
		}

		return page
	}
}
