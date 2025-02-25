package prepare

import (
	"github.com/Sicarii-NaumaN/poloniex-gateway-api/config"
	"github.com/Sicarii-NaumaN/poloniex-gateway-api/internal/handler"
	"github.com/Sicarii-NaumaN/poloniex-gateway-api/internal/router"
	syncer "github.com/Sicarii-NaumaN/poloniex-gateway-api/internal/syncer"
	"github.com/Sicarii-NaumaN/poloniex-gateway-api/internal/ws"
	"github.com/Sicarii-NaumaN/poloniex-gateway-api/service/adapter/db"
	poloniex_api "github.com/Sicarii-NaumaN/poloniex-gateway-api/service/adapter/poloniex-api"
	poloniex_ws "github.com/Sicarii-NaumaN/poloniex-gateway-api/service/adapter/poloniex-ws"
	"github.com/Sicarii-NaumaN/poloniex-gateway-api/service/builder"
	"github.com/Sicarii-NaumaN/poloniex-gateway-api/service/builder/repository"
	"github.com/Sicarii-NaumaN/poloniex-gateway-api/service/processor"
	syncerRep "github.com/Sicarii-NaumaN/poloniex-gateway-api/service/processor/repository"
	"github.com/Sicarii-NaumaN/poloniex-gateway-api/service/rt_saver"
	rtSaverRep "github.com/Sicarii-NaumaN/poloniex-gateway-api/service/rt_saver/repository"
	"github.com/Sicarii-NaumaN/poloniex-gateway-api/tools/closer"
	"github.com/Sicarii-NaumaN/poloniex-gateway-api/tools/logger"
	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"net/http"
)

func InitImpl(adapter db.IAdapter, poloniexSvc handler.PoloniexSvc, port int) http.Handler {
	isSwaggerCreated := config.GetConfigBool(config.IsSwaggerCreated)

	api := handler.NewImplementation(poloniexSvc)

	r := router.CreateRouter(api, port, isSwaggerCreated)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"*"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodDelete,
			http.MethodPost,
			http.MethodPut,
		},
		AllowCredentials: true,
	})

	return c.Handler(r)
}

func InitPoloniexAdapter() (poloniexSvc handler.PoloniexSvc) {
	poloniexSvc = poloniex_api.NewAdapter()
	return
}

func InitSyncer(db db.IAdapter, poloniexAdapter handler.PoloniexSvc) *syncer.Service {
	processorSvc := processor.NewService(poloniexAdapter, syncerRep.NewRepository(db))
	builderSvc := builder.NewService(repository.NewRepository(db))
	candleSyncer := syncer.NewService(processorSvc, builderSvc)

	return candleSyncer
}

func InitListener(db db.IAdapter) (*ws.Service, error) {
	// Типо берем из env
	const url = "wss://ws.poloniex.com/ws/public"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		logger.Fatalf("Failed to connect to Poloniex WS: %v", err)
	}
	// Не успел перенести в другой слой
	adapter, err := poloniex_ws.NewAdapter(conn, url)

	if err != nil {
		return nil, err
	}

	saverSvc := rt_saver.NewService(rtSaverRep.NewRepository(db))

	closer.Add(adapter)

	return ws.NewService(adapter, saverSvc), err
}
