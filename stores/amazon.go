package stores

import (
	"log"

	"github.com/go-rod/bypass"
	"github.com/go-rod/rod"
)

func CheckAmazon() StockCheckResult {
	browser := rod.New().MustConnect()
	defer browser.Close()

	page := bypass.MustPage(browser)
	defer page.Close()

	// Product details page
	page.MustNavigate("https://www.amazon.co.uk/Xbox-RRT-00007-Series-X/dp/B08H93GKNJ/ref=sr_1_1")
	page.MustWaitLoad()

	// Setting Sleeper to nil to not retry
	_, err := page.Sleeper(nil).Element("#buybox-see-all-buying-choices-announce")

	if err == nil {
		return OutOfStock
	} else if err.Error() == "cannot find element" {
		return InStock
	} else {
		log.Println(err)

		return ErrorOccurred
	}
}
