package poloniex

type Interval string

const (
	OneMin     = Interval("MINUTE_1")
	FifteenMin = Interval("MINUTE_15")
	OneHour    = Interval("HOUR_1")
	OneDay     = Interval("DAY_1")
)

type IntervalType int

const (
	IntervalTypeUnknown IntervalType = iota
	IntervalTypeOneMin
	IntervalTypeFifteenMin
	IntervalTypeOneHour
	IntervalTypeOneDay
)

var IntervalToType = map[Interval]IntervalType{
	OneMin:     IntervalTypeOneMin,
	FifteenMin: IntervalTypeFifteenMin,
	OneHour:    IntervalTypeOneHour,
	OneDay:     IntervalTypeOneDay,
}

var TypeToInterval = map[IntervalType]Interval{
	IntervalTypeUnknown:    "",
	IntervalTypeOneMin:     OneMin,
	IntervalTypeFifteenMin: FifteenMin,
	IntervalTypeOneHour:    OneHour,
	IntervalTypeOneDay:     OneDay,
}
