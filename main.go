package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

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

	err := notify(results)

	if err != nil {
		log.Println(err)
	}
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
