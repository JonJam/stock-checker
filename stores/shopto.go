package stores

import (
	"log"

	"github.com/go-rod/rod"
)

type ShopTo struct {
}

func (s ShopTo) Check(pool rod.PagePool, create func() *rod.Page) StockCheckResult {
	const storeName = "ShopTo"

	page := pool.Get(create)
	if page == nil {
		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}

	defer pool.Put(page)

	// Product details page
	if err := page.Navigate("https://www.shopto.net/en/xbxhw01-xbox-series-x-p191471/?utm_source=website&utm_medium=banner&utm_campaign=Xbox%20Series%20X"); err != nil {
		log.Println(err)

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}
	if err := page.WaitLoad(); err != nil {
		log.Println(err)

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}

	// Setting Sleeper to nil to not retry
	if _, err := page.Sleeper(nil).ElementR("h1", "404"); err == nil {
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
		log.Println(err)

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}
}
