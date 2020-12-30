package services

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jonjam/stock-checker/config"
	configMocks "github.com/jonjam/stock-checker/config/mocks"
	"github.com/stretchr/testify/require"
)

func TestSend_SendsExpectedRequest(t *testing.T) {
	// ARRANGE
	var req *http.Request

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseMultipartForm(1048576)
		req = r
	}))
	defer ts.Close()

	twilioConfig := config.TwilioConfig{
		Url:        ts.URL,
		AccountSid: "SID",
		AuthToken:  "AUTH_TOKEN",
		NumberFrom: "NUMBER_FROM",
		NumberTo:   "NUMBER_TO",
	}
	c := &configMocks.Config{}
	c.On("GetTwilioConfig").Return(twilioConfig)

	body := "SMS content"

	client := NewTwilioClient(c)

	// ACT
	_ = client.Send(body)

	// ASSERT
	require.Equal(t, "POST", req.Method)
	require.Equal(t, fmt.Sprintf("/2010-04-01/Accounts/%s/Messages.json", twilioConfig.AccountSid), req.URL.RequestURI())

	username, password, ok := req.BasicAuth()
	require.True(t, ok)
	require.Equal(t, twilioConfig.AccountSid, username)
	require.Equal(t, twilioConfig.AuthToken, password)

	require.Equal(t, "application/json", req.Header.Get("Accept"))
	require.Equal(t, "application/x-www-form-urlencoded", req.Header.Get("Content-Type"))

	require.Equal(t, body, req.FormValue("Body"))
	require.Equal(t, twilioConfig.NumberFrom, req.FormValue("From"))
	require.Equal(t, twilioConfig.NumberTo, req.FormValue("To"))
}
