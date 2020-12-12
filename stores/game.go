package stores

import (
	"log"

	"github.com/go-rod/bypass"
	"github.com/go-rod/rod"
)

func checkGame() StockCheckResult {
	browser := rod.New().MustConnect()
	defer browser.Close()

	page := bypass.MustPage(browser)
	defer page.Close()

	// Xbox Series X page (https://www.game.co.uk/xbox-series-x)
	page.MustNavigate("https://www.game.co.uk/xbox-series-x")
	page.MustWaitLoad()

	consolesSection := page.MustElement("#contentPanelsConsoles")

	consoleTitleElement := consolesSection.MustElementR("h3", "Series X")

	panelItemElement := consoleTitleElement.MustParent()

	_, err := panelItemElement.ElementR("a", "Out of stock")

	if err == nil {
		return OutOfStock
	} else if err.Error() == "cannot find element" {
		return InStock
	} else {
		log.Println(err)

		return ErrorOccurred
	}
}
