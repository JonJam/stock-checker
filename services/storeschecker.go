package main

import (
	"time"

	"github.com/go-rod/bypass"
	"github.com/go-rod/rod"
	"github.com/jonjam/stock-checker/stores"

	"sync"

	"github.com/go-rod/rod/lib/launcher"
)

func Run(storesSlice []stores.Store) []stores.StockCheckResult{
	// TODO this should be controlled via config
	devMode := true

	url := createControlUrl(devMode)

	browser := createBrowser(url, devMode)
	defer browser.MustClose()

	// TODO This should be controlled via config
	pool := rod.NewPagePool(3)
	create := createCreatePageFunc(browser, devMode)

	c := make(chan stores.StockCheckResult)

	for _, s := range storesSlice {
		go func(store stores.Store) {
			c <- store.Check(pool, create)
		}(s)
	}

	results := make([]stores.StockCheckResult, 0, len(storesSlice))

	for i := 0; i < len(storesSlice); i++ {
		results = append(results, <-c)
	}

	pool.Cleanup(func(p *rod.Page) { p.MustClose() })

	return results
}

func createControlUrl(devMode bool) string {
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