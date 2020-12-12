package stores

import (
	"log"

	"github.com/go-rod/rod"
)

type JohnLewis struct {
}

func (j JohnLewis) Check(pool rod.PagePool, create func() *rod.Page) StockCheckResult {
	page := pool.Get(create)
	defer pool.Put(page)

	// Home page
	page.MustNavigate("https://www.johnlewis.com/")
	page.MustWaitLoad()

	// Cookie banner
	page.MustElement(`button[data-test="allow-all"]`).MustClick()

	// Search
	searchInputElement := page.MustElement("#desktopSearch")
	searchInputElement.MustInput("xbox series x console")

	searchButton := searchInputElement.MustNext()
	searchButton.MustClick()

	// Search results page (https://www.johnlewis.com/search)
	page.MustWaitLoad()

	productCardTitleElement := page.MustElementR("h2", "Microsoft Xbox Series X Console")

	productCardLinkElement := productCardTitleElement.MustParent().MustParent()

	productCardLinkElement.MustClick()

	// Product details page
	page.MustWaitLoad()

	// Setting Sleeper to nil to not retry
	_, err := page.Sleeper(nil).Element("#button--add-to-basket-out-of-stock")

	const storeName = "John Lewis"

	if err == nil {
		return StockCheckResult{
			storeName: storeName,
			status:    OutOfStock,
		}
	} else if err.Error() == "cannot find element" {
		return StockCheckResult{
			storeName: storeName,
			status:    InStock,
		}
	} else {
		log.Println(err)

		return StockCheckResult{
			storeName: storeName,
			status:    Unknown,
		}
	}
}
