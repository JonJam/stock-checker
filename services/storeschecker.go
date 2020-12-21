package services

import (
	"time"

	"github.com/go-rod/bypass"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/jonjam/stock-checker/config"
	"github.com/jonjam/stock-checker/stores"
	"github.com/jonjam/stock-checker/util"
)

func CheckStores(storesSlice []stores.Store) []stores.StockCheckResult {
	url, err := createControlURL()

	if err != nil {
		util.Logger.Fatalln(err)
	}

	browser, err := createBrowser(url)

	if err != nil {
		util.Logger.Fatalln(err)
	}

	defer func() {
		err := browser.Close()

		if err != nil {
			util.Logger.Println(err)
		}
	}()

	pool := rod.NewPagePool(config.GetRodConfig().PagePoolSize)
	defer pool.Cleanup(func(p *rod.Page) {
		err := p.Close()

		if err != nil {
			util.Logger.Println(err)
		}
	})

	c := make(chan stores.StockCheckResult)

	get := createGetPageFunc(browser, pool)
	release := createReleasePageFunc(pool)

	for _, s := range storesSlice {
		go func(store stores.Store) {
			c <- store.Check(get, release)
		}(s)
	}

	numOfStores := len(storesSlice)
	results := make([]stores.StockCheckResult, 0, numOfStores)

	for i := 0; i < numOfStores; i++ {
		results = append(results, <-c)
	}

	return results
}

func createControlURL() (string, error) {
	launcher := launcher.New()

	launcher.Devtools(config.GetRodConfig().DevTools)
	launcher.Headless(config.GetRodConfig().Headless)

	return launcher.Launch()
}

func createBrowser(url string) (*rod.Browser, error) {
	browser := rod.New().ControlURL(url)

	browser.Trace(config.GetRodConfig().Trace)

	if config.GetRodConfig().SlowMotion {
		browser.SlowMotion(time.Second * 10)
	}

	err := browser.Connect()

	if err != nil {
		return nil, err
	}

	return browser, nil
}

func createGetPageFunc(browser *rod.Browser, pool rod.PagePool) func() *rod.Page {
	create := createCreatePageFunc(browser)

	// Gets a page from the pool and configures a timeout for store to perform all operations with it
	return func() *rod.Page {
		// TODO Implement timeout
		return pool.Get(create)
	}
}

func createReleasePageFunc(pool rod.PagePool) func(*rod.Page) {
	return func(page *rod.Page) {
		// TODO Implement cancel timeout
		pool.Put(page)
	}
}

func createCreatePageFunc(browser *rod.Browser) func() *rod.Page {
	// This func will create a new configured page will be contained within a different incognito browser window.
	// It returns nil when an error occurs rather than exposing error due to https://pkg.go.dev/github.com/go-rod/rod#PagePool.Get
	return func() *rod.Page {
		browser, err := browser.Incognito()

		if err != nil {
			util.Logger.Println(err)

			return nil
		}

		page, err := bypass.Page(browser)

		if err != nil {
			util.Logger.Println(err)

			return nil
		}

		return page
	}
}
