package syncer

import (
	"context"
	"github.com/Sicarii-NaumaN/poloniex-gateway-api/tools/logger"
	"sync"
	"time"
)

type Processor interface {
	Process(context.Context, chan struct{}, time.Time) error
}

type Builder interface {
	Build(context.Context, int64) error
}

type Service struct {
	proc    Processor
	builder Builder

	mu                   *sync.RWMutex
	oneMinLatestCandleTs int64 // Будем завязываться на минутную свечу которая последней засинкалась в процессоре
}

func NewService(proc Processor, builder Builder) *Service {
	return &Service{proc: proc, builder: builder, mu: &sync.RWMutex{}}
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

	timeBeforeAll := time.Now()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		// Когда время открытия свечи последний свечи станет больше времени запуска приложения,
		//значит можем выключаться, тк остальное подберет WS
		case <-ctx.Done():
			logger.Info(ctx, "Successfully synced")
			return
		case <-ticker.C:
			err = s.proc.Process(ctx, finishChan, timeBeforeAll)
			if err != nil {
				cancel()
				return
			}
		}

	}
}

// Хотел изначально для разнообразия в кронджобе реализовать, но времени не было
func (s *Service) RunCandlesBuilder(ctx context.Context) (err error) {
	var from int64
	for {
		s.mu.RLock()
		defer s.mu.RUnlock()
		if s.oneMinLatestCandleTs != 0 { // Ждём, когда синкер архивных свечей закончит работу
			from = s.oneMinLatestCandleTs
			break
		}
	}

	logger.Info(ctx, "Start candles builder")
	defer logger.Info(ctx, "End candles builder")

	ticker := time.NewTicker(45 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err = s.builder.Build(ctx, from)
			if err != nil {
				return
			}
		}

	}
}
