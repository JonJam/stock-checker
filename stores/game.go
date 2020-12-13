package stores

import (
	"github.com/go-rod/rod"
	"github.com/jonjam/stock-checker/util"
)

type Game struct {
}

func (g Game) Check(getPage func() *rod.Page, releasePage func(*rod.Page)) StockCheckResult {
	const storeName = "Game"

	page := getPage()
	if page == nil {
		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}
	defer releasePage(page)

	// Xbox Series X page (https://www.game.co.uk/xbox-series-x)
	if err := page.Navigate("https://www.game.co.uk/xbox-series-x"); err != nil {
		util.Logger.Println(err)

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}
	if err := page.WaitLoad(); err != nil {
		util.Logger.Println(err)

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}

	consolesSection, err := page.Element("#contentPanelsConsoles")
	if err != nil {
		util.Logger.Println(err)

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}

	consoleTitle, err := consolesSection.ElementR("h3", "Series X")
	if err != nil {
		util.Logger.Println(err)

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}

	panelItem, err := consoleTitle.Parent()
	if err != nil {
		util.Logger.Println(err)

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}

	if _, err = panelItem.ElementR("a", "Out of stock"); err == nil {
		return StockCheckResult{
			StoreName: storeName,
			Status:    OutOfStock,
		}
	} else if err.Error() == "cannot find element" {
		return StockCheckResult{
			StoreName: storeName,
			Status:    InStock,
		}
	} else {
		util.Logger.Println(err)

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}
}
