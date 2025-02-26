package poloniex

const (
	diffOneMin     = 60 * 1000
	diffFifteenMin = 15 * diffOneMin
	diffOneHour    = 4 * diffFifteenMin
	diffOneDay     = 24 * diffOneHour
)

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

type StartEndInterval struct {
	Start int64
	End   int64
}

func intervalDurationMs(interval Interval) int64 {
	switch IntervalToType[interval] {
	case IntervalTypeOneMin:
		return diffOneMin
	case IntervalTypeFifteenMin:
		return diffFifteenMin
	case IntervalTypeOneHour:
		return diffOneHour
	case IntervalTypeOneDay:
		return diffOneDay
	default:
		return diffOneMin
	}
}

// Возвращает начало и конец свечи для каждого интервала
func GetCandleIntervalsByTime(ts int64) map[Interval]StartEndInterval {
	return map[Interval]StartEndInterval{
		OneMin:     GetCandleIntervalByTime(ts, OneMin),
		FifteenMin: GetCandleIntervalByTime(ts, FifteenMin),
		OneHour:    GetCandleIntervalByTime(ts, OneHour),
		OneDay:     GetCandleIntervalByTime(ts, OneDay),
	}
}

func GetCandleIntervalByTime(ts int64, interval Interval) StartEndInterval {
	diff := intervalDurationMs(interval)
	start := ts - (ts % diff)
	end := start + diff

	return StartEndInterval{
		Start: start,
		End:   end,
	}
}
