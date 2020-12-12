package stores

import (
	"log"

	"github.com/go-rod/rod"
)

type Argos struct {
}

func (a Argos) Check(pool rod.PagePool, create func() *rod.Page) StockCheckResult {
	page := pool.Get(create)
	defer pool.Put(page)

	// Home (https://www.argos.co.uk/)
	page.MustNavigate("https://www.argos.co.uk/")
	page.MustWaitLoad()

	// Cookie banner
	page.MustElement("#consent_prompt_submit").MustClick()

	// Search
	page.MustElement("#searchTerm").MustInput("xbox series x console")
	page.MustElement(`button[data-test="search-button"]`).MustClick()

	// Search page (https://www.argos.co.uk/search/)
	page.MustWaitLoad()
	page.MustElementR("a", "Xbox Series X 1TB Console").MustClick()

	// Out of stock page (https://www.argos.co.uk/vp/oos/xbox.html)
	page.MustWaitLoad()

	// Setting Sleeper to nil to not retry
	_, err := page.Sleeper(nil).ElementR("h1", "Sorry, Xbox is currently unavailable.")

	const storeName = "Argos"

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
