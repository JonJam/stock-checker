package stores

import (
	"log"

	"github.com/go-rod/rod"
)

type Smyths struct {
}

func (s Smyths) Check(pool rod.PagePool, create func() *rod.Page) StockCheckResult {
	page := pool.Get(create)
	defer pool.Put(page)

	// Product details page
	page.MustNavigate("https://www.smythstoys.com/uk/en-gb/video-games-and-tablets/xbox-gaming/xbox-series-x-%7c-s/xbox-series-x-%7c-s-consoles/xbox-series-x-1tb-console/p/192012")
	page.MustWaitLoad()

	// Cookie banner
	page.MustElementR("button", "Yes, I’m happy").MustClick()

	addToCartButton := page.MustElement("#addToCartButton")

	value, err := addToCartButton.Attribute("disabled")

	const storeName = "Smyths"

	if err != nil {
		log.Println(err)

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	} else if value != nil {
		return StockCheckResult{
			StoreName: storeName,
			Status:    OutOfStock,
		}
	} else {
		return StockCheckResult{
			StoreName: storeName,
			Status:    InStock,
		}
	}
}
