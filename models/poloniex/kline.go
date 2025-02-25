package poloniex

type GetCandleStickReq struct {
	Pair      Pair
	Interval  Interval
	StartTime int64
	EndTime   int64
	Limit     int64
}

type GetCandleStickResp struct {
	Candles []Kline
	Total   int64
}

// Структура KL:
type Kline struct {
	Pair      Pair     `json:"pair"`       // название пары в Bitsgap
	TimeFrame Interval `json:"time_frame"` // период формирования свечи (1m, 1h, 1d)
	O         float64  `json:"o"`          // open - цена открытия
	H         float64  `json:"h"`          // high - максимальная цена
	L         float64  `json:"l"`          // low - минимальная цена
	C         float64  `json:"c"`          // close - цена закрытия
	UtcBegin  int64    `json:"utc_begin"`  // время unix начала формирования свечки
	UtcEnd    int64    `json:"utc_end"`    // время unix окончания формирования свечки
	VolumeBS  VBS      `json:"volume_bs"`
}
type VBS struct {
	BuyBase   float64 `json:"buy_base"`   // объём покупок в базовой валюте
	SellBase  float64 `json:"sell_base"`  // объём продаж в базовой валюте
	BuyQuote  float64 `json:"buy_quote"`  // объём покупок в котируемой валюте
	SellQuote float64 `json:"sell_quote"` // объём продаж в котируемой валюте
}
