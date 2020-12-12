package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/jonjam/stock-checker/stores"
)

func main() {
	// TODO .rod config is for dev only

	results := map[string]stores.StockCheckResult{
		"Argos": stores.CheckArgos(),
		// "Game": stores.CheckGame(),
		// "John Lewis": stores.CheckJohnLewis(),
		// "Amazon": stores.CheckAmazon(),
		// "Smyths": stores.CheckSmyths(),
		// "Currys": stores.CheckCurrys(),
		// "ShopTo": stores.CheckShopTo(),
	}

	log.Println(results)

	// TODO Disabled while testing
	// err := notify(results)

	// if err != nil {
	// 	log.Println(err)
	// }
}

// Based off: https://www.twilio.com/blog/2017/09/send-text-messages-golang.html
func notify(results map[string]stores.StockCheckResult) error {
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
