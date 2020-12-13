package stores

import (
	"log"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

type Smyths struct {
}

func (s Smyths) Check(pool rod.PagePool, create func() *rod.Page) StockCheckResult {
	const storeName = "Smyths"

	page := pool.Get(create)
	if page == nil {
		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}

	defer pool.Put(page)

	// Product details page
	if err := page.Navigate("https://www.smythstoys.com/uk/en-gb/video-games-and-tablets/xbox-gaming/xbox-series-x-%7c-s/xbox-series-x-%7c-s-consoles/xbox-series-x-1tb-console/p/192012"); err != nil {
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

	// Cookie banner
	cookieBanner, err := page.ElementR("button", "Yes, I’m happy")
	if err != nil {
		log.Println(err)

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}
	if err = cookieBanner.Click(proto.InputMouseButtonLeft); err != nil {
		log.Println(err)

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}

	addToCartButton, err := page.Element("#addToCartButton")
	if err != nil {
		log.Println(err)

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}

	if value, err := addToCartButton.Attribute("disabled"); err != nil {
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
