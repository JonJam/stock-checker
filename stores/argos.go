package stores

import (
	"log"

	"github.com/go-rod/bypass"
	"github.com/go-rod/rod"
)

func CheckArgos() StockCheckResult {
	browser := rod.New().MustConnect()
	defer browser.Close()

	page := bypass.MustPage(browser)

	defer page.Close()

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

	if err == nil {
		return OutOfStock
	} else if err.Error() == "cannot find element" {
		return InStock
	} else {
		log.Println(err)

		return ErrorOccurred
	}
}