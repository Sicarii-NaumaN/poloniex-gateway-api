package poloniex_ws

import (
	"github.com/Sicarii-NaumaN/poloniex-gateway-api/models/poloniex"
	"strconv"
	"strings"
)

func repackRT(t tradeData) poloniex.RecentTrade {
	return poloniex.RecentTrade{
		Tid:       strconv.Itoa(int(t.ID)), // По условию оно string
		Pair:      poloniex.Pair(t.Symbol),
		Price:     t.Price,
		Amount:    t.Amount,
		Side:      poloniex.Side(strings.ToUpper(t.TakerSide)),
		Timestamp: t.TS,
	}
}
