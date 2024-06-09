package nextdate

import "time"

func NextYear(now time.Time, t time.Time, repeat string) (string, error) {
	t = t.AddDate(1, 0, 0)
	return t.Format("20060102"), nil
}
