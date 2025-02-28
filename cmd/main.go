// App entrypoint, you to implement it

package main

import (
	"context"
	"fmt"
	"github.com/Sicarii-NaumaN/poloniex-gateway-api/cmd/prepare"
	"github.com/Sicarii-NaumaN/poloniex-gateway-api/config"
	"github.com/Sicarii-NaumaN/poloniex-gateway-api/tools/closer"
	"github.com/Sicarii-NaumaN/poloniex-gateway-api/tools/logger"
	"github.com/Sicarii-NaumaN/poloniex-gateway-api/tools/xcontext"
	_ "github.com/lib/pq"
	"net/http"
)

func main() {
	ctx := context.Background()

	defer func() {
		logger.Info("service has shut down gracefully")
		closer.Wait()
		closer.CloseAll()
		logger.Info("service gracefully down successfully")
	}()

	db := prepare.NewDBConn(ctx)
	poloniexSvc := prepare.InitPoloniexAdapter()

	// Тут очень много способов синхронизации между архивными свечами
	// и текущими можно придумать и нормально реализовать,
	// к сожалению тогда тестовое слишком раздуется
	listener, err := prepare.InitListener(db)
	if err != nil {
		logger.Fatalf("failed to init listener: %v", err)
	}

	go func() {
		dctx := xcontext.NewDetachedContext(ctx)
		if err = listener.RunListening(dctx); err != nil {
			logger.Errorf("fatal error syncer.RunListening: %v", err)
		}
	}()

	syncer := prepare.InitSyncer(db, poloniexSvc)

	// Если мы хотим формировать свечи ОДНОВРЕМЕННО с запуском приложения, то. нам для этого нужно по REST
	// достать предыдущие трейды для каждого размера свечи (GetCandleIntervalByTime) для времени запуска приложения
	go func() {
		dctx := xcontext.NewDetachedContext(ctx)
		err = syncer.RunSyncCandles(dctx)
		if err != nil {
			logger.Errorf("fatal error syncer.RunSyncCandles: %v", err)
		}
	}()
	go func() {
		dctx := xcontext.NewDetachedContext(ctx)
		// Вот тут логика преобразования трейдов свечи
		if err = syncer.RunCandlesBuilder(dctx); err != nil {
			logger.Errorf("fatal error syncer.RunCandlesBuilder: %v", err)
		}
	}()

	// Тут представлен вариант где мы дожидаемся завершение синка и только после этого строим свечи
	//var wg sync.WaitGroup
	//wg.Add(1)
	//go func() {
	//	dctx := xcontext.NewDetachedContext(ctx)
	//	// Тут происходит синк архивных свечей
	//	err = syncer.RunSyncCandles(dctx)
	//	if err != nil {
	//		logger.Errorf("fatal error syncer.RunSyncCandles: %v", err)
	//	}
	//}()
	//wg.Wait()
	//go func() {
	//	dctx := xcontext.NewDetachedContext(ctx)
	//	// Вот тут логика преобразования трейдов свечи
	//	if err = syncer.RunCandlesBuilder(dctx); err != nil {
	//		logger.Errorf("fatal error syncer.RunCandlesBuilder: %v", err)
	//	}
	//}()

	port := config.GetConfigInt(config.Port)
	logger.Info(fmt.Sprintf("Started server at :%d. Swagger docs stated at %d", port, port+1))
	logger.Error(http.ListenAndServe(fmt.Sprintf(":%d", port), prepare.InitImpl(db, poloniexSvc, port)))
}
