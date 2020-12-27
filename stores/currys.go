package stores

import (
	"github.com/go-rod/rod"
	"go.uber.org/zap"
)

type Currys struct {
	logger *zap.Logger
}

func NewCurrys(l *zap.Logger) Currys {
	return Currys{
		logger: l,
	}
}

func (c Currys) Check(getPage func() *rod.Page, releasePage func(*rod.Page)) StockCheckResult {
	const storeName = "Currys"

	page := getPage()
	if page == nil {
		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}
	defer releasePage(page)

	// Product details page
	if err := page.Navigate("https://www.currys.co.uk/gbuk/gaming/console-gaming/consoles/microsoft-xbox-series-x-1-tb-10203371-pdt.html"); err != nil {
		c.logger.Error("Failed to navigate to Product details page", zap.Error(err))

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}
	if err := page.WaitLoad(); err != nil {
		c.logger.Error("Failed to wait for Product details page to load.", zap.Error(err))

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}

	// Setting Sleeper to nil to not retry
	if _, err := page.Sleeper(nil).ElementR("li", "Sorry this item is out of stock"); err == nil {
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
		c.logger.Error("Error occurred finding stock status list item.", zap.Error(err))

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}
}
