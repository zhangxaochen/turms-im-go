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

func DivideDuration(start, end time.Time, divideBy DivideBy) []DateRange {
	if start.IsZero() || end.IsZero() || !start.Before(end) {
		return nil
	}
	var ranges []DateRange
	current := start
	for current.Before(end) {
		var next time.Time
		switch divideBy {
		case DivideBy_HOUR:
			next = current.Add(time.Hour)
		case DivideBy_DAY:
			next = current.AddDate(0, 0, 1)
		case DivideBy_MONTH:
			next = current.AddDate(0, 1, 0)
		case DivideBy_NOOP:
			next = end
		default:
			next = end
		}
		if next.After(end) {
			next = end
		}
		ranges = append(ranges, DateRange{Start: current, End: next})
		current = next
	}
	return ranges
}
