package time

import "time"

type DateRange struct {
	Start time.Time
	End   time.Time
}

func (r DateRange) IsEmpty() bool {
	return r.Start.IsZero() && r.End.IsZero()
}

type DivideBy string

const (
	DivideBy_NOOP  DivideBy = "NOOP"
	DivideBy_HOUR  DivideBy = "HOUR"
	DivideBy_DAY   DivideBy = "DAY"
	DivideBy_MONTH DivideBy = "MONTH"
)
