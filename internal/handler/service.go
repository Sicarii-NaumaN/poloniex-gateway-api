package handler

import (
	"context"
	"github.com/Sicarii-NaumaN/poloniex-gateway-api/models/poloniex"
	"net/http"
)

type API interface {
	GetCandleSticks(rw http.ResponseWriter, r *http.Request)
	GetTrades(rw http.ResponseWriter, r *http.Request)
}

type PoloniexSvc interface {
	GetCandleSticks(ctx context.Context, req poloniex.GetCandleStickReq) (poloniex.GetCandleStickResp, error)
	GetTrades(ctx context.Context, req poloniex.GetTradesReq) (poloniex.GetTradesResp, error)
}

type implementation struct {
	poloniexSvc PoloniexSvc
}

func NewImplementation(poloniexSvc PoloniexSvc) API {
	return &implementation{
		poloniexSvc: poloniexSvc,
	}
}
