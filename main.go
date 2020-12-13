package main

import (
	"fmt"
	"sort"

	"log"

	"github.com/jonjam/stock-checker/services"
	"github.com/jonjam/stock-checker/stores"
)

func main() {
	stores := []stores.Store{
		stores.Argos{},
		stores.Amazon{},
		stores.Currys{},
		stores.Game{},
		stores.JohnLewis{},
		stores.ShopTo{},
		stores.Smyths{},
	}

	results := services.CheckStores(stores)

	log.Println(results)

	notify(results)
}

// Based off: https://www.twilio.com/blog/2017/09/send-text-messages-golang.html
func notify(results []stores.StockCheckResult) {
	// TODO should be config DO NOT COMMIT
	// accountSid := ""
	// authToken := ""
	// numberTo := ""
	// numberFrom := ""
	// requestURL := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", accountSid)

	body := "Stock checker results:\n"

	sort.Slice(results, func(i, j int) bool {
		return results[i].StoreName < results[j].StoreName
	})

	for _, v := range results {
		body = fmt.Sprintf(body+"%v\n", v)
	}

	log.Println(body)

	// TODO Disabled while testing (in dev mode)
	// msgData := url.Values{}
	// msgData.Set("To", numberTo)
	// msgData.Set("From", numberFrom)
	// msgData.Set("Body", body)
	// msgDataReader := *strings.NewReader(msgData.Encode())

	// client := &http.Client{}
	// req, err := http.NewRequest("POST", requestURL, &msgDataReader)

	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }

	// req.SetBasicAuth(accountSid, authToken)
	// req.Header.Add("Accept", "application/json")
	// req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// resp, err := client.Do(req)

	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }

	// if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
	// 	log.Println(fmt.Errorf("unexpected status code from Twilio: %v", resp.StatusCode))
	// 	return
	// }

	// // Success HTTP Status but need to check response for error
	// var data map[string]interface{}
	// decoder := json.NewDecoder(resp.Body)
	// err = decoder.Decode(&data)

	// if err != nil {
	// 	log.Println(err)
	// }
}
