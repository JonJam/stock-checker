package stores

import (
	"github.com/go-rod/rod"
	"go.uber.org/zap"
)

type ShopTo struct {
	logger *zap.Logger
}

func NewShopTo(l *zap.Logger) ShopTo {
	return ShopTo{
		logger: l,
	}
}

func (s ShopTo) Check(getPage func() *rod.Page, releasePage func(*rod.Page)) StockCheckResult {
	const storeName = "ShopTo"

	page := getPage()
	if page == nil {
		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}
	defer releasePage(page)

	// Product details page
	if err := page.Navigate("https://www.shopto.net/en/xbxhw01-xbox-series-x-p191471/"); err != nil {
		s.logger.Error("Failed to navigate to product details page.", zap.Error(err))

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}
	if err := page.WaitLoad(); err != nil {
		s.logger.Error("Failed to wait for product details page to load.", zap.Error(err))

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}

	// Setting Sleeper to nil to not retry
	if _, err := page.Sleeper(nil).ElementR("button", "REGISTER NOW"); err == nil {
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
		s.logger.Error("Error occurred finding register now button.", zap.Error(err))

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}
}
