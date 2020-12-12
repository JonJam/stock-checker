package stores

import (
	"log"

	"github.com/go-rod/rod"
)

type Currys struct {
}

func (c Currys) Check(pool rod.PagePool, create func() *rod.Page) StockCheckResult {
	page := pool.Get(create)
	defer pool.Put(page)

	// Product details page
	page.MustNavigate("https://www.currys.co.uk/gbuk/gaming/console-gaming/consoles/microsoft-xbox-series-x-1-tb-10203371-pdt.html")
	page.MustWaitLoad()

	// Setting Sleeper to nil to not retry
	_, err := page.Sleeper(nil).ElementR("li", "Sorry this item is out of stock")

	const storeName = "Currys"

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
