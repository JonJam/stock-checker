package stores

import (
	"log"

	"github.com/go-rod/bypass"
	"github.com/go-rod/rod"
)

func CheckShopTo() StockCheckResult {
	browser := rod.New().MustConnect()
	defer browser.Close()

	page := bypass.MustPage(browser)
	defer page.Close()

	// Product details page
	page.MustNavigate("https://www.shopto.net/en/xbxhw01-xbox-series-x-p191471/?utm_source=website&utm_medium=banner&utm_campaign=Xbox%20Series%20X")
	page.MustWaitLoad()

	// Setting Sleeper to nil to not retry
	_, err := page.Sleeper(nil).ElementR("h1", "404")

	if err == nil {
		return OutOfStock
	} else if err.Error() == "cannot find element" {
		return InStock
	} else {
		log.Println(err)

		return ErrorOccurred
	}
}
