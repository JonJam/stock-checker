package stores

import (
	"log"

	"github.com/go-rod/rod"
)

type Game struct {
}

func (g Game) Check(pool rod.PagePool, create func() *rod.Page) StockCheckResult {
	page := pool.Get(create)
	defer pool.Put(page)

	// Xbox Series X page (https://www.game.co.uk/xbox-series-x)
	page.MustNavigate("https://www.game.co.uk/xbox-series-x")
	page.MustWaitLoad()

	consolesSection := page.MustElement("#contentPanelsConsoles")

	consoleTitleElement := consolesSection.MustElementR("h3", "Series X")

	panelItemElement := consoleTitleElement.MustParent()

	_, err := panelItemElement.ElementR("a", "Out of stock")

	const storeName = "Game"

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
