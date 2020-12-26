package stores

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"go.uber.org/zap"
)

type Argos struct {
	logger *zap.Logger
}

func NewArgos(l *zap.Logger) Argos {
	return Argos{
		logger: l,
	}
}

func (a Argos) Check(getPage func() *rod.Page, releasePage func(*rod.Page)) StockCheckResult {
	a.logger.Info("Checking store.")
	const storeName = "Argos"

	page := getPage()
	if page == nil {
		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}
	defer releasePage(page)

	// Home (https://www.argos.co.uk/)
	if err := page.Navigate("https://www.argos.co.uk/"); err != nil {
		a.logger.Error("Failed to navigate to Argos home page.", zap.Error(err))

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}
	if err := page.WaitLoad(); err != nil {
		a.logger.Error("Failed to wait for Argos home page to load.", zap.Error(err))

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}

	// Cookie banner
	cookieBannerSubmit, err := page.Element("#consent_prompt_submit")
	if err != nil {
		a.logger.Error("Failed to find accept cookies button.", zap.Error(err))

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}
	if err = cookieBannerSubmit.Click(proto.InputMouseButtonLeft); err != nil {
		a.logger.Error("Failed to click accept cookies button.", zap.Error(err))

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}

	// Search
	searchInput, err := page.Element("#searchTerm")
	if err != nil {
		a.logger.Error("Failed to find search term input.", zap.Error(err))

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}
	if err = searchInput.Input("xbox series x console"); err != nil {
		a.logger.Error("Failed to input search term.", zap.Error(err))

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}

	searchButton, err := page.Element(`button[data-test="search-button"]`)
	if err != nil {
		a.logger.Error("Failed to find search button.", zap.Error(err))

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}
	if err = searchButton.Click(proto.InputMouseButtonLeft); err != nil {
		a.logger.Error("Failed to click on search button.", zap.Error(err))

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}

	// Search page (https://www.argos.co.uk/search/)
	if err := page.WaitLoad(); err != nil {
		a.logger.Error("Failed to wait for Search page to load.", zap.Error(err))

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}

	consoleLink, err := page.ElementR("a", "Xbox Series X 1TB Console")
	if err != nil {
		a.logger.Error("Failed to find console link.", zap.Error(err))

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}
	if err = consoleLink.Click(proto.InputMouseButtonLeft); err != nil {
		a.logger.Error("Failed to click on console link.", zap.Error(err))

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}

	// Out of stock page (https://www.argos.co.uk/vp/oos/xbox.html)
	if err := page.WaitLoad(); err != nil {
		a.logger.Error("Failed to wait for console page to load.", zap.Error(err))

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}

	// Setting Sleeper to nil to not retry
	if _, err = page.Sleeper(nil).ElementR("h1", "Sorry, Xbox is currently unavailable."); err == nil {
		return StockCheckResult{
			StoreName: storeName,
			Status:    OutOfStock,
		}
	} else if err.Error() == "cannot find element" {
		return StockCheckResult{
			StoreName: storeName,
			Status:    InStock,
		}
	} else {
		a.logger.Error("Error occurred finding title.", zap.Error(err))

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}
}
