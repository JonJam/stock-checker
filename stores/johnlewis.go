package stores

import (
	"log"

	"github.com/go-rod/bypass"
	"github.com/go-rod/rod"
)

func CheckJohnLewis() StockCheckResult {
	browser := rod.New().MustConnect()
	defer browser.Close()

	page := bypass.MustPage(browser)
	defer page.Close()

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

	if err == nil {
		return OutOfStock
	} else if err.Error() == "cannot find element" {
		return InStock
	} else {
		log.Println(err)

		return ErrorOccurred
	}
}