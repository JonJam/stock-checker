package stores

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/jonjam/stock-checker/util"
)

type Argos struct {
}

func (a Argos) Check(getPage func() *rod.Page, releasePage func(*rod.Page)) StockCheckResult {
	const storeName = "Argos"

	page := getPage()
	if page == nil {
		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}
	defer releasePage(page)

	// Home (https://www.argos.co.uk/)
	if err := page.Navigate("https://www.argos.co.uk/"); err != nil {
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

	// Cookie banner
	cookieBannerSubmit, err := page.Element("#consent_prompt_submit")
	if err != nil {
		util.Logger.Println(err)

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}
	if err = cookieBannerSubmit.Click(proto.InputMouseButtonLeft); err != nil {
		util.Logger.Println(err)

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}

	// Search
	searchInput, err := page.Element("#searchTerm")
	if err != nil {
		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}
	if err = searchInput.Input("xbox series x console"); err != nil {
		util.Logger.Println(err)

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}

	searchButton, err := page.Element(`button[data-test="search-button"]`)
	if err != nil {
		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}
	if err = searchButton.Click(proto.InputMouseButtonLeft); err != nil {
		util.Logger.Println(err)

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}

	// Search page (https://www.argos.co.uk/search/)
	if err := page.WaitLoad(); err != nil {
		util.Logger.Println(err)

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}

	consoleLink, err := page.ElementR("a", "Xbox Series X 1TB Console")
	if err != nil {
		util.Logger.Println(err)

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}
	if err = consoleLink.Click(proto.InputMouseButtonLeft); err != nil {
		util.Logger.Println(err)

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}

	// Out of stock page (https://www.argos.co.uk/vp/oos/xbox.html)
	if err := page.WaitLoad(); err != nil {
		util.Logger.Println(err)

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}

	// Setting Sleeper to nil to not retry
	if _, err = page.Sleeper(nil).ElementR("h1", "Sorry, Xbox is currently unavailable."); err == nil {
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
