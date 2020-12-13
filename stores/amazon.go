package stores

import (
	"log"

	"github.com/go-rod/rod"
)

type Amazon struct {
}

func (a Amazon) Check(pool rod.PagePool, create func() *rod.Page) StockCheckResult {
	const storeName = "Amazon"

	page := pool.Get(create)

	if page == nil {
		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}

	defer pool.Put(page)

	// Product details page
	if err := page.Navigate("https://www.amazon.co.uk/Xbox-RRT-00007-Series-X/dp/B08H93GKNJ/ref=sr_1_1"); err != nil {
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
	if _, err := page.Sleeper(nil).Element("#add-to-cart-button"); err == nil {
		return StockCheckResult{
			StoreName: storeName,
			Status:    InStock,
		}
	} else if err.Error() == "cannot find element" {
		return StockCheckResult{
			StoreName: storeName,
			Status:    OutOfStock,
		}
	} else {
		log.Println(err)

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}
}
