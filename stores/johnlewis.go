package stores

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"go.uber.org/zap"
)

type JohnLewis struct {
	logger *zap.Logger
}

func NewJohnLewis(l *zap.Logger) JohnLewis {
	return JohnLewis{
		logger: l,
	}
}

func (j JohnLewis) Check(getPage func() *rod.Page, releasePage func(*rod.Page)) StockCheckResult {
	const storeName = "John Lewis"

	page := getPage()
	if page == nil {
		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}
	defer releasePage(page)

	// Home page
	if err := page.Navigate("https://www.johnlewis.com/"); err != nil {
		j.logger.Error("Failed to load to home page.", zap.Error(err))

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}
	if err := page.WaitLoad(); err != nil {
		j.logger.Error("Failed to wait for home page to load.", zap.Error(err))

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}

	// Cookie banner
	cookieBannerSubmit, err := page.Element(`button[data-test="allow-all"]`)
	if err != nil {
		j.logger.Error("Failed to find cookie banner accept button.", zap.Error(err))

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}

	if err = cookieBannerSubmit.Click(proto.InputMouseButtonLeft); err != nil {
		j.logger.Error("Failed to click cookie banner accept button.", zap.Error(err))

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}

	// Search
	searchInput, err := page.Element("#desktopSearch")
	if err != nil {
		j.logger.Error("Failed to find search input.", zap.Error(err))

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}

	if err = searchInput.Input("xbox series x console"); err != nil {
		j.logger.Error("Failed to enter search term.", zap.Error(err))

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}

	searchButton, err := searchInput.Next()
	if err != nil {
		j.logger.Error("Failed to find search button.", zap.Error(err))

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}

	if err = searchButton.Click(proto.InputMouseButtonLeft); err != nil {
		j.logger.Error("Failed to click on search button.", zap.Error(err))

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}

	// Search results page (https://www.johnlewis.com/search)
	if err = page.WaitLoad(); err != nil {
		j.logger.Error("Failed to wait for search results page to load.", zap.Error(err))

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}

	productCardTitle, err := page.ElementR("h2", "Microsoft Xbox Series X Console")
	if err != nil {
		j.logger.Error("Failed to find console title.", zap.Error(err))

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}

	productCardLinkDiv, err := productCardTitle.Parent()
	if err != nil {
		j.logger.Error("Failed to find console card.", zap.Error(err))

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}

	productCardLink, err := productCardLinkDiv.Parent()
	if err != nil {
		j.logger.Error("Failed to find console link.", zap.Error(err))

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}

	if err := productCardLink.Click(proto.InputMouseButtonLeft); err != nil {
		j.logger.Error("Failed to click on console link.", zap.Error(err))

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}

	// Product details page
	if err = page.WaitLoad(); err != nil {
		j.logger.Error("Failed to wait for Product details page to load.", zap.Error(err))

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}

	addToBasketButton, err := page.Element("div.add-to-basket-summary-and-cta > button")
	if err != nil {
		j.logger.Error("Failed to find add to basket button.", zap.Error(err))

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}

	addToBasketButtonText, err := addToBasketButton.Text()

	if err != nil {
		j.logger.Error("Failed to get add to basket button text.", zap.Error(err))

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}

	if addToBasketButtonText == "Out of stock" {
		return StockCheckResult{
			StoreName: storeName,
			Status:    OutOfStock,
		}
	}

	return StockCheckResult{
		StoreName: storeName,
		Status:    InStock,
	}
}
