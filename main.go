package main

import (
	"fmt"
	"log"

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
	r := checkArgos()

	results := map[string]stockCheckResult{
		"Argos": r,
	}

	notify(results)
}

func checkArgos() stockCheckResult {
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
		return OutOfStock
	} else if err.Error() == "cannot find element" {
		return InStock
	} else {
		log.Println(err)

		return ErrorOccurred
	}
}

func notify(results map[string]stockCheckResult) {
	// TODO should be config DO NOT COMMIT
	// accountSid := ""
	// authToken := ""
	// numberTo := ""
	// numberFrom := ""

	body := "Stock checker results:\n"

	for k, v := range results {
		body = fmt.Sprintf(body+"&v: %v\n", k, v.String())
	}

	// TODO Follow https://www.twilio.com/blog/2017/09/send-text-messages-golang.html

	log.Println(body)
}
