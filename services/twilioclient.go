package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/jonjam/stock-checker/config"
)

type TwilioClient struct {
	config config.TwilioConfig
}

func NewTwilioClient(c config.Config) TwilioClient {
	return TwilioClient{
		config: c.GetTwilioConfig(),
	}
}

func (c TwilioClient) Send(body string) error {
	// Based off: https://www.twilio.com/blog/2017/09/send-text-messages-golang.html
	requestURL := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", c.config.AccountSid)

	msgData := url.Values{}
	msgData.Set("To", c.config.NumberTo)
	msgData.Set("From", c.config.NumberFrom)
	msgData.Set("Body", body)
	msgDataReader := *strings.NewReader(msgData.Encode())

	client := &http.Client{}
	req, err := http.NewRequest("POST", requestURL, &msgDataReader)

	if err != nil {
		return err
	}

	req.SetBasicAuth(c.config.AccountSid, c.config.AuthToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		return fmt.Errorf("Unexpected status code from Twilio: %v", resp.StatusCode)
	}

	// Success HTTP Status but need to check response for error
	var data map[string]interface{}
	decoder := json.NewDecoder(resp.Body)

	return decoder.Decode(&data)
}
