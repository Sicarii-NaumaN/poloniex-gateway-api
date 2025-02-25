package ws

import (
	"context"
	"github.com/Sicarii-NaumaN/poloniex-gateway-api/models/poloniex"
	"github.com/Sicarii-NaumaN/poloniex-gateway-api/tools/logger"
)

type Listener interface {
	Subscribe(ctx context.Context, channel string, symbols []poloniex.Pair) error
	Listen(ctx context.Context, tradeChan chan<- poloniex.RecentTrade) error
}

type TradesSaver interface {
	SaveRealTimeTrades(ctx context.Context, tradeChan <-chan poloniex.RecentTrade) error
}

type Service struct {
	tradesSaver TradesSaver
	listener    Listener

	symbols []poloniex.Pair
	channel string
}

func NewService(listener Listener, tradesSaver TradesSaver) *Service {
	return &Service{
		listener:    listener,
		tradesSaver: tradesSaver,
		symbols: []poloniex.Pair{
			poloniex.BtcUsdt,
			poloniex.TrxUsdt,
			poloniex.EthUsdt,
			poloniex.DogeUsdt,
			poloniex.BchUsdt,
		},
		channel: "trades",
	}
}

func (s *Service) RunListening(ctx context.Context) error {
	logger.Info(ctx, "Start listening trades")
	defer logger.Info(ctx, "End listening trades")

	err := s.listener.Subscribe(ctx, s.channel, s.symbols)
	if err != nil {
		logger.Errorf("error while subscribtion %v", err)
		return err
	}

	var tradeChan = make(chan poloniex.RecentTrade)

	go func() {
		err = s.listener.Listen(ctx, tradeChan)
		if err != nil {
			logger.Errorf("error while listening %v", err)
		}
	}()

	err = s.tradesSaver.SaveRealTimeTrades(ctx, tradeChan)
	if err != nil {
		logger.Errorf("error while SaveRealTimeTrades %v", err)
		return err
	}

	return nil
}
