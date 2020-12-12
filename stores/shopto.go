package stores

import (
	"log"

	"github.com/go-rod/rod"
)

type ShopTo struct {
}

func (s ShopTo) Check(pool rod.PagePool, create func() *rod.Page) StockCheckResult {
	page := pool.Get(create)
	defer pool.Put(page)

	// Product details page
	page.MustNavigate("https://www.shopto.net/en/xbxhw01-xbox-series-x-p191471/?utm_source=website&utm_medium=banner&utm_campaign=Xbox%20Series%20X")
	page.MustWaitLoad()

	// Setting Sleeper to nil to not retry
	_, err := page.Sleeper(nil).ElementR("h1", "404")

	const storeName = "ShopTo"

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
