package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/jonjam/stock-checker/config"
	"github.com/jonjam/stock-checker/stores"
	"github.com/jonjam/stock-checker/util"
)

// Based off: https://www.twilio.com/blog/2017/09/send-text-messages-golang.html
func Notify(results []stores.StockCheckResult) {
	body := "Stock checker results:\n"

	sort.Slice(results, func(i, j int) bool {
		return results[i].StoreName < results[j].StoreName
	})

	for _, v := range results {
		body = fmt.Sprintf(body+"%v\n", v)
	}

	config := config.GetTwilioConfig()

	requestURL := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", config.AccountSid)

	msgData := url.Values{}
	msgData.Set("To", config.NumberTo)
	msgData.Set("From", config.NumberFrom)
	msgData.Set("Body", body)
	msgDataReader := *strings.NewReader(msgData.Encode())

	client := &http.Client{}
	req, err := http.NewRequest("POST", requestURL, &msgDataReader)

	if err != nil {
		util.Logger.Println(err)
		return
	}

	req.SetBasicAuth(config.AccountSid, config.AuthToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)

	if err != nil {
		util.Logger.Println(err)
		return
	}

	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		util.Logger.Println(fmt.Errorf("unexpected status code from Twilio: %v", resp.StatusCode))
		return
	}

	// Success HTTP Status but need to check response for error
	var data map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&data)

	if err != nil {
		util.Logger.Println(err)
	}
}
