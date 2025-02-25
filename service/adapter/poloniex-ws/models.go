package poloniex_ws

import "github.com/Sicarii-NaumaN/poloniex-gateway-api/models/poloniex"

type subscriptionMessage struct {
	Event   string          `json:"event"`
	Channel []string        `json:"channel"`
	Symbols []poloniex.Pair `json:"symbols"`
}

type tradeData struct {
	Symbol     string `json:"symbol"`
	Amount     string `json:"amount"`
	TakerSide  string `json:"takerSide"`
	Quantity   string `json:"quantity"`
	CreateTime int64  `json:"createTime"`
	Price      string `json:"price"`
	ID         int64  `json:"id,string"`
	TS         int64  `json:"ts"`
}

type incomingMessage struct {
	Channel string      `json:"channel"`
	Data    []tradeData `json:"data"`
}
