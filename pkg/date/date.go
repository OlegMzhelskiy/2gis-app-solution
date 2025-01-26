package date

import (
	"fmt"
	"time"
)

type CustomDate struct {
	time.Time
}

func (cd *CustomDate) UnmarshalJSON(data []byte) error {
	str := string(data)
	str = str[1 : len(str)-1]

	parsedTime, err := time.Parse("2006-01-02", str)
	if err != nil {
		return fmt.Errorf("invalid date format, use YYYY-MM-DD: %w", err)
	}

	cd.Time = parsedTime
	return nil
}

func Date(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}
