// Here you have to implement router,
// register the handlers and within middlewares

package router

import (
	"net/http"

	"github.com/Sicarii-NaumaN/poloniex-gateway-api/internal/handler"
	"github.com/Sicarii-NaumaN/poloniex-gateway-api/internal/handler/response"
	"github.com/Sicarii-NaumaN/poloniex-gateway-api/tools/swagger"

	"github.com/gorilla/mux"
)

const AppName = "Gateway"

// CreateRouter router
func CreateRouter(impl handler.API, appPort int, isSwaggerCreated bool) *mux.Router {
	var handlers = []swagger.Handler{
		{
			HandlerFunc:      impl.GetCandleSticks,
			Path:             "/api/v1/candles",
			Method:           http.MethodGet,
			Description:      "Get candles",
			ResponseBody:     response.Response[string]{},
			ResponseMimeType: swagger.MimeJson,
			Opts:             []swagger.Option{},
			Tag:              "Trade",
		},
		{
			HandlerFunc:      impl.GetTrades,
			Path:             "/api/v1/trades",
			Method:           http.MethodGet,
			Description:      "Get trades",
			ResponseBody:     response.Response[string]{},
			ResponseMimeType: swagger.MimeJson,
			Opts:             []swagger.Option{},
			Tag:              "Trade",
		},
	}

	return NewAPI(AppName, appPort, isSwaggerCreated, handlers)
}
