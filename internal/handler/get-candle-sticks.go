package handler

import (
	"net/http"
	"time"

	"github.com/Sicarii-NaumaN/poloniex-gateway-api/internal/handler/response"
	"github.com/Sicarii-NaumaN/poloniex-gateway-api/models/poloniex"
)

func (i *implementation) GetCandleSticks(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	res, err := i.poloniexSvc.GetCandleSticks(ctx, poloniex.GetCandleStickReq{
		Pair:      poloniex.BtcUsdt,
		Interval:  poloniex.OneHour,
		StartTime: time.Now().Add(-time.Hour * 24).UTC().UnixMilli(),
		EndTime:   time.Now().Add(-time.Hour * 16).UTC().UnixMilli(),
		Limit:     100,
	})
	if err != nil {
		handleErrResponse(rw, err)
		return
	}

	response.OK(rw, res)
}
