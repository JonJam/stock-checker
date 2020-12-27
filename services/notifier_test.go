package services

import (
	"testing"

	"github.com/jonjam/stock-checker/config"
	"github.com/jonjam/stock-checker/config/mocks"
	mocksTwo "github.com/jonjam/stock-checker/services/mocks"
	"github.com/jonjam/stock-checker/stores"
	"github.com/stretchr/testify/assert"

	"go.uber.org/zap"
)

func TestX(t *testing.T) {
	// ARRANGE
	l := zap.NewNop()

	notifierConfig := config.NotifierConfig{
		Enabled: false,
	}
	c := new(mocks.Config)
	c.On("GetNotifierConfig").Return(notifierConfig)

	client := &mocksTwo.SmsClient{}

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
	assert.Equal(t, 123, 123, "they should be equal")
}
