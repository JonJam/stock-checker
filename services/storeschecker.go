package services

import (
	"time"

	"github.com/go-rod/bypass"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/jonjam/stock-checker/stores"
)

func CheckStores(storesSlice []stores.Store) []stores.StockCheckResult {
	// TODO this should be controlled via config
	devMode := true

	url := createControlURL(devMode)

	browser := createBrowser(url, devMode)
	defer browser.MustClose()

	// TODO This should be controlled via config. Dev: 1, Prod:3
	pool := rod.NewPagePool(3)
	defer pool.Cleanup(func(p *rod.Page) { p.MustClose() })

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

func createControlURL(devMode bool) string {
	launcher := launcher.New()

	if devMode {
		launcher.Devtools(true)
		launcher.Headless(false)
	}

	return launcher.MustLaunch()
}

func createBrowser(url string, devMode bool) *rod.Browser {
	browser := rod.New().ControlURL(url)

	if devMode {
		browser.Trace(true)
		browser.SlowMotion(time.Second)
	}

	return browser.MustConnect()
}

func createCreatePageFunc(browser *rod.Browser, devMode bool) func() *rod.Page {
	return func() *rod.Page {
		page := bypass.MustPage(browser.MustIncognito())

		if devMode {
			page.MustWindowFullscreen()
		}

		return page
	}
}
