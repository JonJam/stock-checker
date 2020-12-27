package stores

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"go.uber.org/zap"
)

type Smyths struct {
	logger *zap.Logger
}

func NewSmyths(l *zap.Logger) Smyths {
	return Smyths{
		logger: l,
	}
}

func (s Smyths) Check(getPage func() *rod.Page, releasePage func(*rod.Page)) StockCheckResult {
	const storeName = "Smyths"

	page := getPage()
	if page == nil {
		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}
	defer releasePage(page)

	// Product details page
	if err := page.Navigate("https://www.smythstoys.com/uk/en-gb/video-games-and-tablets/xbox-gaming/xbox-series-x-%7c-s/xbox-series-x-%7c-s-consoles/xbox-series-x-1tb-console/p/192012"); err != nil {
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

	// Cookie banner
	cookieBanner, err := page.ElementR("button", "Yes, I’m happy")
	if err != nil {
		s.logger.Error("Failed to find cookie banner accept button.", zap.Error(err))

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}
	if err = cookieBanner.Click(proto.InputMouseButtonLeft); err != nil {
		s.logger.Error("Failed to click on cookie banner accept button.", zap.Error(err))

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}

	addToCartButton, err := page.Element("#addToCartButton")
	if err != nil {
		s.logger.Error("Failed to find add to cart button.", zap.Error(err))

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}

	if value, err := addToCartButton.Attribute("disabled"); err != nil {
		s.logger.Error("Failed to get disabled attribute.", zap.Error(err))

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	} else if value != nil {
		return StockCheckResult{
			StoreName: storeName,
			Status:    OutOfStock,
		}
	} else {
		return StockCheckResult{
			StoreName: storeName,
			Status:    InStock,
		}
	}
}
