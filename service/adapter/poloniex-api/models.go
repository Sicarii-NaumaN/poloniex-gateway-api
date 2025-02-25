package poloniex_api

type tradeData struct {
	Id         string `json:"id"`
	Price      string `json:"price"`
	Quantity   string `json:"quantity"`
	Amount     string `json:"amount"`
	TakerSide  string `json:"takerSide"`
	Ts         int64  `json:"ts"`
	CreateTime int64  `json:"createTime"`
}

// candleData candlestick
type candleData struct {
	Low              string `json:"low"`
	High             string `json:"high"`
	Open             string `json:"open"`
	Close            string `json:"close"`
	Amount           string `json:"amount"`
	Quantity         string `json:"quantity"`
	BuyTakerAmount   string `json:"buyTakerAmount"`
	BuyTakerQuantity string `json:"buyTakerQuantity"`
	TradeCount       int    `json:"tradeCount"`
	Timestamp        int64  `json:"ts"`
	WeightedAverage  string `json:"weightedAverage"`
	Interval         string `json:"interval"`
	StartTime        int64  `json:"startTime"`
	CloseTime        int64  `json:"closeTime"`
}
