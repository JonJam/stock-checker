package stores

import (
	"log"

	"github.com/go-rod/rod"
)

type Amazon struct {
}

func (a Amazon) Check(pool rod.PagePool, create func() *rod.Page) StockCheckResult {
	page := pool.Get(create)
	defer pool.Put(page)

	// Product details page
	page.MustNavigate("https://www.amazon.co.uk/Xbox-RRT-00007-Series-X/dp/B08H93GKNJ/ref=sr_1_1")
	page.MustWaitLoad()

	// Setting Sleeper to nil to not retry
	_, err := page.Sleeper(nil).Element("#buybox-see-all-buying-choices-announce")

	const storeName = "Amazon"

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
