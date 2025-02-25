package poloniex

type Pair string

const (
	Unknown  Pair = "UNKNOWN"
	BtcUsdt       = Pair("BTC_USDT")
	TrxUsdt       = Pair("TRX_USDT")
	EthUsdt       = Pair("ETH_USDT")
	DogeUsdt      = Pair("DOGE_USDT")
	BchUsdt       = Pair("BCH_USDT")
)

type PairType int

const (
	PairTypeUnknown PairType = iota
	PairTypeBtcUsdt
	PairTypeTrxUsdt
	PairTypeEthUsdt
	PairTypeDogeUsdt
	PairTypeBchUsdt
)

var PairsToTypeMap = map[Pair]PairType{
	BtcUsdt:  PairTypeBtcUsdt,
	TrxUsdt:  PairTypeTrxUsdt,
	EthUsdt:  PairTypeEthUsdt,
	DogeUsdt: PairTypeDogeUsdt,
	BchUsdt:  PairTypeBchUsdt,
}

var TypePairsToMap = map[PairType]Pair{
	PairTypeBtcUsdt:  BtcUsdt,
	PairTypeTrxUsdt:  TrxUsdt,
	PairTypeEthUsdt:  EthUsdt,
	PairTypeDogeUsdt: DogeUsdt,
	PairTypeBchUsdt:  BchUsdt,
}

type Side string

const (
	None = Side("UNKNOWN")
	Buy  = Side("BUY")
	Sell = Side("SELL")
)

type SideType int

const (
	SideTypeUnknown SideType = iota
	SideTypeBuy
	SideTypeSell
)

var SideToTypeMap = map[Side]SideType{
	None: SideTypeUnknown,
	Buy:  SideTypeBuy,
	Sell: SideTypeSell,
}

var TypeToSideMap = map[SideType]Side{
	SideTypeBuy:  Buy,
	SideTypeSell: Sell,
}
