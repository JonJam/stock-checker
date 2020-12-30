package services

import (
	"testing"

	"github.com/jonjam/stock-checker/config"
	configMocks "github.com/jonjam/stock-checker/config/mocks"
	servicesMocks "github.com/jonjam/stock-checker/services/mocks"
	"github.com/jonjam/stock-checker/stores"

	"go.uber.org/zap"
)

func TestNotify_Disabled_DoesNotSendSMS(t *testing.T) {
	// ARRANGE
	l := zap.NewNop()

	notifierConfig := config.NotifierConfig{
		Enabled: false,
	}
	c := &configMocks.Config{}
	c.On("GetNotifierConfig").Return(notifierConfig)

	client := &servicesMocks.SmsClient{}
	n := NewNotifier(c, l, client)

	results := []stores.StockCheckResult{
		{
			StoreName: "Test",
			Status:    stores.InStock,
		},
	}

	// ACT
	n.Notify(results)

	// ASSERT
	client.AssertNotCalled(t, "Send")
}

func TestNotify_Enabled_SendsSMS(t *testing.T) {
	// ARRANGE
	l := zap.NewNop()

	notifierConfig := config.NotifierConfig{
		Enabled: true,
	}
	c := &configMocks.Config{}
	c.On("GetNotifierConfig").Return(notifierConfig)

	results := []stores.StockCheckResult{
		{
			StoreName: "Test",
			Status:    stores.InStock,
		},
	}
	msg := "Stock checker results:\nTest: In stock\n"

	client := &servicesMocks.SmsClient{}
	client.On("Send", msg).Return(nil)
	n := NewNotifier(c, l, client)

	// ACT
	n.Notify(results)

	// ASSERT
	client.AssertCalled(t, "Send", msg)
}
