package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-rod/bypass"
	"github.com/go-rod/rod"
)

// Based off: https://programming.guide/go/define-enumeration-string.html
type stockCheckResult int

const (
	InStock stockCheckResult = iota
	OutOfStock
	ErrorOccurred
)

func (s stockCheckResult) String() string {
	return [...]string{
		"In stock",
		"Out of stock",
		"Error occurred"}[s]
}

func main() {
	// TODO .rod config is for dev only

	results := map[string]stockCheckResult{
		// "Argos": checkArgos(),
		// "Game": checkGame(),
		// "John Lewis": checkJohnLewis(),
		// "Amazon": checkAmazon(),
		// "Smyths": checkSmyths(),
		"Currys": checkCurrys(),
	}

	log.Println(results)

	// TODO Disabled while testing
	// err := notify(results)

	// if err != nil {
	// 	log.Println(err)
	// }
}

func checkArgos() stockCheckResult {
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

func checkGame() stockCheckResult {
	browser := rod.New().MustConnect()
	defer browser.Close()

	page := bypass.MustPage(browser)
	defer page.Close()

	// Xbox Series X page (https://www.game.co.uk/xbox-series-x)
	page.MustNavigate("https://www.game.co.uk/xbox-series-x")
	page.MustWaitLoad()

	consolesSection := page.MustElement("#contentPanelsConsoles")

	consoleTitleElement := consolesSection.MustElementR("h3", "Series X")

	panelItemElement := consoleTitleElement.MustParent()

	_, err := panelItemElement.ElementR("a", "Out of stock")

	if err == nil {
		return OutOfStock
	} else if err.Error() == "cannot find element" {
		return InStock
	} else {
		log.Println(err)

		return ErrorOccurred
	}
}

func checkJohnLewis() stockCheckResult {
	browser := rod.New().MustConnect()
	defer browser.Close()

	page := bypass.MustPage(browser)
	defer page.Close()

	// Home page
	page.MustNavigate("https://www.johnlewis.com/")
	page.MustWaitLoad()

	// Cookie banner
	page.MustElement(`button[data-test="allow-all"]`).MustClick()

	// Search
	searchInputElement := page.MustElement("#desktopSearch")
	searchInputElement.MustInput("xbox series x console")

	searchButton := searchInputElement.MustNext()
	searchButton.MustClick()

	// Search results page (https://www.johnlewis.com/search)
	page.MustWaitLoad()

	productCardTitleElement := page.MustElementR("h2", "Microsoft Xbox Series X Console")

	productCardLinkElement := productCardTitleElement.MustParent().MustParent()

	productCardLinkElement.MustClick()

	// Product details page
	page.MustWaitLoad()

	// Setting Sleeper to nil to not retry
	_, err := page.Sleeper(nil).Element("#button--add-to-basket-out-of-stock")

	if err == nil {
		return OutOfStock
	} else if err.Error() == "cannot find element" {
		return InStock
	} else {
		log.Println(err)

		return ErrorOccurred
	}
}

func checkAmazon() stockCheckResult {
	browser := rod.New().MustConnect()
	defer browser.Close()

	page := bypass.MustPage(browser)
	defer page.Close()

	// Product details page
	page.MustNavigate("https://www.amazon.co.uk/Xbox-RRT-00007-Series-X/dp/B08H93GKNJ/ref=sr_1_1")
	page.MustWaitLoad()

	// Setting Sleeper to nil to not retry
	_, err := page.Sleeper(nil).Element("#buybox-see-all-buying-choices-announce")

	if err == nil {
		return OutOfStock
	} else if err.Error() == "cannot find element" {
		return InStock
	} else {
		log.Println(err)

		return ErrorOccurred
	}
}

func checkSmyths() stockCheckResult {
	browser := rod.New().MustConnect()
	defer browser.Close()

	page := bypass.MustPage(browser)
	defer page.Close()

	// Product details page
	page.MustNavigate("https://www.smythstoys.com/uk/en-gb/video-games-and-tablets/xbox-gaming/xbox-series-x-%7c-s/xbox-series-x-%7c-s-consoles/xbox-series-x-1tb-console/p/192012")
	page.MustWaitLoad()

	// Cookie banner
	page.MustElementR("button", "Yes, I’m happy").MustClick()

	addToCartButton := page.MustElement("#addToCartButton")

	value, err := addToCartButton.Attribute("disabled")

	if err != nil {
		log.Println(err)

		return ErrorOccurred
	} else if value != nil {
		return OutOfStock
	} else {
		return InStock
	}
}

func checkCurrys() stockCheckResult {
	browser := rod.New().MustConnect()
	defer browser.Close()

	page := bypass.MustPage(browser)
	defer page.Close()

	// Product details page
	page.MustNavigate("https://www.currys.co.uk/gbuk/gaming/console-gaming/consoles/microsoft-xbox-series-x-1-tb-10203371-pdt.html")
	page.MustWaitLoad()

	// Setting Sleeper to nil to not retry
	_, err := page.Sleeper(nil).ElementR("li", "Sorry this item is out of stock")

	if err == nil {
		return OutOfStock
	} else if err.Error() == "cannot find element" {
		return InStock
	} else {
		log.Println(err)

		return ErrorOccurred
	}
}

// Based off: https://www.twilio.com/blog/2017/09/send-text-messages-golang.html
func notify(results map[string]stockCheckResult) error {
	// TODO should be config DO NOT COMMIT
	accountSid := ""
	authToken := ""
	numberTo := ""
	numberFrom := ""
	requestURL := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", accountSid)

	body := "Stock checker results:\n"

	for k, v := range results {
		body = fmt.Sprintf(body+"%s: %s\n", k, v.String())
	}

	msgData := url.Values{}
	msgData.Set("To", numberTo)
	msgData.Set("From", numberFrom)
	msgData.Set("Body", body)
	msgDataReader := *strings.NewReader(msgData.Encode())

	client := &http.Client{}
	req, err := http.NewRequest("POST", requestURL, &msgDataReader)

	if err != nil {
		return err
	}

	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		return fmt.Errorf("unexpected status code from Twilio: %v", resp.StatusCode)
	}

	// Success HTTP Status but need to check response for error
	var data map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&data)

	return err
}
