package stores

import (
	"log"

	"github.com/go-rod/bypass"
	"github.com/go-rod/rod"
)

func CheckCurrys() StockCheckResult {
	browser := rod.New().MustConnect()
	defer browser.Close()

	page := bypass.MustPage(browser)
	defer page.Close()

	// Product details page
	page.MustNavigate("https://www.currys.co.uk/gbuk/gaming/console-gaming/consoles/microsoft-xbox-series-x-1-tb-10203371-pdt.html")
	page.MustWaitLoad()

	// Setting Sleeper to nil to not retry
	_, err := page.Sleeper(nil).ElementR("li", "Sorry this item is out of stock")

	if err == nil {
		return OutOfStock
	} else if err.Error() == "cannot find element" {
		return InStock
	} else {
		log.Println(err)

		return ErrorOccurred
	}
}
