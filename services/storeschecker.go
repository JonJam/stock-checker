package services

import (
	"fmt"
	"time"

	"github.com/go-rod/bypass"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/jonjam/stock-checker/config"
	"github.com/jonjam/stock-checker/stores"

	"go.uber.org/zap"
)

type StoresChecker struct {
	config config.Config
	logger *zap.Logger
}

func NewStoresChecker(c config.Config, l *zap.Logger) StoresChecker {
	return StoresChecker{
		config: c,
		logger: l,
	}
}

func (s StoresChecker) CheckStores(storesSlice []stores.Store) []stores.StockCheckResult {
	url, err := s.createControlURL()

	if err != nil {
		s.logger.Error("Failed to create control URL.", zap.Error(err))

		return []stores.StockCheckResult{}
	}

	browser, err := s.createBrowser(url)

	if err != nil {
		s.logger.Error("Failed to create browser.", zap.Error(err))

		return []stores.StockCheckResult{}
	}

	defer func() {
		err := browser.Close()

		if err != nil {
			s.logger.Error("Failed to close browser.", zap.Error(err))
		}
	}()

	pool := rod.NewPagePool(s.config.GetRodConfig().PagePoolSize)
	defer pool.Cleanup(func(p *rod.Page) {
		err := p.Close()

		if err != nil {
			s.logger.Error("Failed to close page.", zap.Error(err))
		}
	})

	c := make(chan stores.StockCheckResult)

	get := s.createGetPageFunc(browser, pool)
	release := s.createReleasePageFunc(pool)

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

func (s StoresChecker) createControlURL() (string, error) {
	launcher := launcher.New().Set("--no-sandbox")

	launcher.Devtools(s.config.GetRodConfig().DevTools)
	launcher.Headless(s.config.GetRodConfig().Headless)

	return launcher.Launch()
}

func (s StoresChecker) createBrowser(url string) (*rod.Browser, error) {
	browser := rod.New().ControlURL(url)

	browser.Logger(newCustomRodLogger(s.logger))

	browser.Trace(s.config.GetRodConfig().Trace)

	if s.config.GetRodConfig().SlowMotion {
		browser.SlowMotion(time.Second)
	}

	err := browser.Connect()

	if err != nil {
		return nil, err
	}

	return browser, nil
}

func (s StoresChecker) createGetPageFunc(browser *rod.Browser, pool rod.PagePool) func() *rod.Page {
	create := s.createCreatePageFunc(browser)

	// Gets a page from the pool and configures a timeout for store to perform all operations with it
	return func() *rod.Page {
		// TODO Implement timeout
		return pool.Get(create)
	}
}

func (s StoresChecker) createReleasePageFunc(pool rod.PagePool) func(*rod.Page) {
	return func(page *rod.Page) {
		// TODO Implement cancel timeout
		pool.Put(page)
	}
}

func (s StoresChecker) createCreatePageFunc(browser *rod.Browser) func() *rod.Page {
	// This func will create a new configured page will be contained within a different incognito browser window.
	// It returns nil when an error occurs rather than exposing error due to https://pkg.go.dev/github.com/go-rod/rod#PagePool.Get
	return func() *rod.Page {
		browser, err := browser.Incognito()

		if err != nil {
			s.logger.Error("Failed to create incognito browser.", zap.Error(err))

			return nil
		}

		page, err := bypass.Page(browser)

		if err != nil {
			s.logger.Error("Failed to create page.", zap.Error(err))

			return nil
		}

		return page
	}
}

type customRodLogger struct {
	logger *zap.Logger
}

func newCustomRodLogger(l *zap.Logger) customRodLogger {
	return customRodLogger{
		logger: l,
	}
}

func (c customRodLogger) Println(msg ...interface{}) {
	s := fmt.Sprintln(msg...)

	c.logger.Debug(s)
}
