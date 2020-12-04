package main

import (
	"log"
	"time"

	"github.com/go-rod/rod"
)

func main() {
	// TODO .rod config is for dev only
	// Argos

	// Home (https://www.argos.co.uk/)
	page := rod.New().MustConnect().MustPage("https://www.argos.co.uk/")

	page.MustWaitLoad().MustScreenshot("temp/a.png")

	// Cookie banner
	page.MustElement("#consent_prompt_submit").MustClick()

	// Search
	page.MustElement("#searchTerm").MustInput("xbox series x console")
	page.MustElement(`button[data-test="search-button"]`).MustClick()

	// Search page (https://www.argos.co.uk/search/)
	productLinkElement := page.MustElementR("a", "Xbox Series X 1TB Console")

	productCardTextContainerElement := productLinkElement.MustParent()

	_, err := productCardTextContainerElement.Element(`img[alt="out of stock"]`)

	if err == nil {
		log.Println("Out of stock")

		// TODO return result with out of stock
	} else if err.Error() == "cannot find element" {
		log.Println("Have stock")

		// TODO maybe remove this
		page.MustWaitLoad().MustScreenshot("tmp/argos.png")

		// TODO notify things since have stock
	} else {
		// TODO return err
	}

	// TODO This is for dev only
	time.Sleep(time.Hour)
}
