package services

import (
	"fmt"
	"sort"

	"github.com/jonjam/stock-checker/config"
	"github.com/jonjam/stock-checker/stores"
	"go.uber.org/zap"
)

type Notifier struct {
	config    config.NotifierConfig
	logger    *zap.Logger
	smsClient SmsClient
}

func NewNotifier(c config.Config, l *zap.Logger, s SmsClient) Notifier {
	return Notifier{
		config:    c.GetNotifierConfig(),
		logger:    l,
		smsClient: s,
	}
}

func (n Notifier) Notify(results []stores.StockCheckResult) {
	if !n.config.Enabled {
		return
	}

	body := "Stock checker results:\n"

	sort.Slice(results, func(i, j int) bool {
		return results[i].StoreName < results[j].StoreName
	})

	for _, v := range results {
		body = fmt.Sprintf(body+"%v\n", v)
	}

	if err := n.smsClient.Send(body); err != nil {
		n.logger.Error("Failed to send SMS.", zap.Error(err))
	}
}
