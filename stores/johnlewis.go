package stores

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/jonjam/stock-checker/util"
)

type JohnLewis struct {
}

func (j JohnLewis) Check(getPage func() *rod.Page, releasePage func(*rod.Page)) StockCheckResult {
	const storeName = "John Lewis"

	page := getPage()
	if page == nil {
		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}
	defer releasePage(page)

	// Home page
	if err := page.Navigate("https://www.johnlewis.com/"); err != nil {
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
	cookieBannerSubmit, err := page.Element(`button[data-test="allow-all"]`)
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
	searchInput, err := page.Element("#desktopSearch")
	if err != nil {
		util.Logger.Println(err)

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

	searchButton, err := searchInput.Next()
	if err != nil {
		util.Logger.Println(err)

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

	// Search results page (https://www.johnlewis.com/search)
	if err = page.WaitLoad(); err != nil {
		util.Logger.Println(err)

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}

	productCardTitle, err := page.ElementR("h2", "Microsoft Xbox Series X Console")
	if err != nil {
		util.Logger.Println(err)

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}

	productCardLinkDiv, err := productCardTitle.Parent()
	if err != nil {
		util.Logger.Println(err)

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}

	productCardLink, err := productCardLinkDiv.Parent()
	if err != nil {
		util.Logger.Println(err)

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}

	if err := productCardLink.Click(proto.InputMouseButtonLeft); err != nil {
		util.Logger.Println(err)

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}

	// Product details page
	if err = page.WaitLoad(); err != nil {
		util.Logger.Println(err)

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}

	addToBasketButton, err := page.Element("div.add-to-basket-summary-and-cta > button")
	if err != nil {
		util.Logger.Println(err)

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}

	addToBasketButtonText, err := addToBasketButton.Text()

	if err != nil {
		util.Logger.Println(err)

		return StockCheckResult{
			StoreName: storeName,
			Status:    Unknown,
		}
	}

	if addToBasketButtonText == "Out of stock" {
		return StockCheckResult{
			StoreName: storeName,
			Status:    OutOfStock,
		}
	}

	return StockCheckResult{
		StoreName: storeName,
		Status:    InStock,
	}
}
