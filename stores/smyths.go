package stores

import (
	"log"

	"github.com/go-rod/bypass"
	"github.com/go-rod/rod"
)

func CheckSmyths() StockCheckResult {
	browser := rod.New().MustConnect()
	defer browser.Close()

	page := bypass.MustPage(browser)
	defer page.Close()

	// Product details page
	page.MustNavigate("https://www.smythstoys.com/uk/en-gb/video-games-and-tablets/xbox-gaming/xbox-series-x-%7c-s/xbox-series-x-%7c-s-consoles/xbox-series-x-1tb-console/p/192012")
	page.MustWaitLoad()

	// Cookie banner
	page.MustElementR("button", "Yes, I’m happy").MustClick()

	addToCartButton := page.MustElement("#addToCartButton")

	value, err := addToCartButton.Attribute("disabled")

	if err != nil {
		log.Println(err)

		return ErrorOccurred
	} else if value != nil {
		return OutOfStock
	} else {
		return InStock
	}
}
