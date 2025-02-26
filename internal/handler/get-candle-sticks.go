package handler

import (
	"github.com/Sicarii-NaumaN/poloniex-gateway-api/internal/handler/response"
	"github.com/Sicarii-NaumaN/poloniex-gateway-api/models/poloniex"
	"net/http"
)

// Для тестов
func (i *implementation) GetCandleSticks(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	res, err := i.poloniexSvc.GetCandleSticks(ctx, poloniex.GetCandleStickReq{
		Pair:      poloniex.BtcUsdt,
		Interval:  poloniex.OneMin,
		StartTime: 1740528060000,
		EndTime:   1740528119999,
		Limit:     100,
	})
	if err != nil {
		handleErrResponse(rw, err)
		return
	}

	response.OK(rw, res)
}
