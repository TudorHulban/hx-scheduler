package hxscheduler

type TimeInterval struct {
	TimeStart int64
	TimeEnd   int64

	SecondsOffset int64
}

func (interval *TimeInterval) BreakDown(perDuration int64) []TimeInterval {
	intervals := make([]TimeInterval, 0)

	if perDuration <= 0 {
		return intervals
	}

	currentStart := interval.TimeStart
	remainingDuration := interval.TimeEnd - interval.TimeStart

	for remainingDuration > 0 {
		duration := min(perDuration, remainingDuration)

		intervals = append(
			intervals,
			TimeInterval{
				TimeStart: currentStart,
				TimeEnd:   currentStart + duration,
			},
		)

		currentStart = currentStart + duration
		remainingDuration = remainingDuration - duration
	}

	return intervals
}

func (interval *TimeInterval) GetUTCTimeStart() int64 {
	return interval.TimeStart - interval.SecondsOffset
}

func (interval *TimeInterval) GetUTCTimeEnd() int64 {
	return interval.TimeEnd - interval.SecondsOffset
}
