package handler

import (
	"github.com/Sicarii-NaumaN/poloniex-gateway-api/internal/handler/response"
	"github.com/Sicarii-NaumaN/poloniex-gateway-api/models/poloniex"
	"net/http"
)

func (i *implementation) GetTrades(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	res, err := i.poloniexSvc.GetTrades(ctx, poloniex.GetTradesReq{
		Pair:  poloniex.BtcUsdt,
		Limit: 1000,
	})

	if err != nil {
		handleErrResponse(rw, err)
		return
	}

	response.OK(rw, res)
}
