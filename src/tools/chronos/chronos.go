package chronos

import "time"

const (
	// DateTimeMask маска с датой и временем
	DateTimeMask = "02.01.2006 15:04:05"
	// DateTimeWithoutSecondMask маска с датой и временем без секунд
	DateTimeWithoutSecondMask = "02.01.2006 15:04"
	// DateTimeMaskLocale маска с датой и временем с зоной
	DateTimeMaskLocale = "02.01.2006 15:04:05 -0700"
	// DateMask маска с датой
	DateMask = "02.01.2006"
	// DateMaskLocale маска с датой и зоной
	DateMaskLocale = "02.01.2006 -0700"
)

func NowTruncUTC() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
}

func DurationBetween(start, end time.Time) time.Duration {
	return end.Sub(start)
}

func SetTimeZone(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), date.Hour(), date.Minute(), date.Second(), date.Nanosecond(), time.Local)
}

func BeginingOfDate(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
}
