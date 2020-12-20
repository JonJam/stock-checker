package stores

import (
	"github.com/go-rod/rod"
	"github.com/jonjam/stock-checker/util"
)

type ShopTo struct {
}

func (s ShopTo) Check(getPage func() *rod.Page, releasePage func(*rod.Page)) StockCheckResult {
	const storeName = "ShopTo"

	page := getPage()
	if page == nil {
		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}
	defer releasePage(page)

	// Product details page
	if err := page.Navigate("https://www.shopto.net/en/swhw33-nintendo-switch-1-1-neon-red-neon-blue-p113520/"); err != nil {
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

	// Setting Sleeper to nil to not retry
	if _, err := page.Sleeper(nil).ElementR("button", "REGISTER NOW"); err == nil {
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
