package stores

import (
	"github.com/go-rod/rod"
	"go.uber.org/zap"
)

type Amazon struct {
	logger *zap.Logger
}

func NewAmazon(l *zap.Logger) Amazon {
	return Amazon{
		logger: l,
	}
}

func (a Amazon) Check(getPage func() *rod.Page, releasePage func(*rod.Page)) StockCheckResult {
	a.logger.Info("Checking store.")

	const storeName = "Amazon"

	page := getPage()
	if page == nil {
		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}
	defer releasePage(page)

	// Product details page
	if err := page.Navigate("https://www.amazon.co.uk/Xbox-RRT-00007-Series-X/dp/B08H93GKNJ/ref=sr_1_1"); err != nil {
		a.logger.Error("Failed to navigate to Product details page.", zap.Error(err))

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}
	if err := page.WaitLoad(); err != nil {
		a.logger.Error("Failed to wait for Product details page to load.", zap.Error(err))

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}

	// Setting Sleeper to nil to not retry
	if _, err := page.Sleeper(nil).Element("#add-to-cart-button"); err == nil {
		return StockCheckResult{
			StoreName: storeName,
			Status:    InStock,
		}
	} else if err.Error() == "cannot find element" {
		return StockCheckResult{
			StoreName: storeName,
			Status:    OutOfStock,
		}
	} else {
		a.logger.Error("Error occurred finding add to card button.", zap.Error(err))

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}
}
