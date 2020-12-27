package stores

import (
	"github.com/go-rod/rod"
	"go.uber.org/zap"
)

type Game struct {
	logger *zap.Logger
}

func NewGame(l *zap.Logger) Game {
	return Game{
		logger: l,
	}
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
		g.logger.Error("Failed to load to console page.", zap.Error(err))

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}
	if err := page.WaitLoad(); err != nil {
		g.logger.Error("Failed to wait for console page to load.", zap.Error(err))

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}

	consolesSection, err := page.Element("#contentPanelsConsoles")
	if err != nil {
		g.logger.Error("Failed to find consoles section.", zap.Error(err))

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}

	consoleTitle, err := consolesSection.ElementR("h3", "Series X")
	if err != nil {
		g.logger.Error("Failed to find console title.", zap.Error(err))

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}

	panelItem, err := consoleTitle.Parent()
	if err != nil {
		g.logger.Error("Failed to find panel item.", zap.Error(err))

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
		g.logger.Error("Error occurred finding out of stock link.", zap.Error(err))

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}
}
