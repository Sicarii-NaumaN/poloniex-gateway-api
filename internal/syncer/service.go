package syncer

import (
	"context"
	"github.com/Sicarii-NaumaN/poloniex-gateway-api/models/poloniex"
	"github.com/Sicarii-NaumaN/poloniex-gateway-api/tools/logger"
	"sync"
	"time"
)

type Processor interface {
	Process(context.Context, chan struct{}, time.Time) error
}

type Builder interface {
	Build(context.Context, map[poloniex.Interval]poloniex.StartEndInterval) error
}

type Service struct {
	proc    Processor
	builder Builder

	mu            *sync.RWMutex
	timeBeforeAll time.Time
}

func NewService(proc Processor, builder Builder) *Service {
	return &Service{proc: proc, builder: builder, mu: &sync.RWMutex{}, timeBeforeAll: time.Now().UTC()}
}

func (s *Service) RunSyncCandles(ctx context.Context) (err error) {
	logger.Info(ctx, "Start syncing candles")
	defer logger.Info(ctx, "End syncing candles")

	ctx, cancel := context.WithCancel(ctx)

	finishChan := make(chan struct{})
	go func() {
		<-finishChan
		cancel()
	}()

	ticker := time.NewTicker(4 * time.Second)
	defer ticker.Stop()

	for {
		select {
		// Когда время открытия свечи последний свечи станет больше времени запуска приложения,
		//значит можем выключаться, тк остальное подберет WS
		case <-ctx.Done():
			logger.Info(ctx, "Successfully synced")
			return
		case <-ticker.C:
			err = s.proc.Process(ctx, finishChan, s.timeBeforeAll)
			if err != nil {
				cancel()
				return
			}
		}

	}
}

const runCandlesBuilderTimer = 30 * time.Second

// Хотел изначально для разнообразия в кронджобе реализовать, но времени не было
func (s *Service) RunCandlesBuilder(ctx context.Context) (err error) {
	logger.Info(ctx, "Start candles builder")
	defer logger.Info(ctx, "End candles builder")

	ticker := time.NewTicker(runCandlesBuilderTimer)
	defer ticker.Stop()

	candleTime := s.timeBeforeAll.UnixMilli()
	// Запрашиваем актуальные свечи для каждого интервала, отталкиваясь от веремени старта приложение
	candleIntervalsByTime := poloniex.GetCandleIntervalsByTime(candleTime)

	// Инициализация
	logger.Info(ctx, "Start init builder call")
	err = s.builder.Build(ctx, candleIntervalsByTime)
	logger.Info(ctx, "End init builder call")
	if err != nil {
		logger.Errorf("Failed to build candles: %v", err)
		return
	}
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Условно даем Build два шанса/итерации закончить все сделки и сформировать итоговые свечи
			// Возможный костыль
			updatedTime := time.Now().Add(-runCandlesBuilderTimer).UTC().UnixMilli()

			// Проверяем не начались ли новые свечи для каждого интервала
			candleIntervalsByTime = s.resolveNewIntervals(updatedTime, candleIntervalsByTime)

			logger.Info(ctx, "Start new builder call....")
			err = s.builder.Build(ctx, candleIntervalsByTime)
			if err != nil {
				logger.Errorf("Failed to build candles: %v", err)
				return
			}
			logger.Info(ctx, "End new builder call")
		}

	}
}

func (s *Service) resolveNewIntervals(updatedTime int64,
	currCandleIntervalsByTime map[poloniex.Interval]poloniex.StartEndInterval) map[poloniex.Interval]poloniex.StartEndInterval {

	var resp = make(map[poloniex.Interval]poloniex.StartEndInterval, len(currCandleIntervalsByTime))
	for i, candleInterval := range currCandleIntervalsByTime {

		// Смотрим что актуальное время не превысило конец свечи заданного интервала,
		// иначе обновляем интервалы
		if updatedTime > candleInterval.End {
			resp[i] = poloniex.GetCandleIntervalByTime(updatedTime, i)
			logger.Infof("Updated new candle in Builder for interval %s with start: %d, end: %d",
				i, candleInterval.Start, candleInterval.End)
		} else {
			resp[i] = candleInterval

		}
	}
	return resp
}
