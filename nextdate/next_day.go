package nextdate

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func NextDay(now time.Time, t time.Time, repeat string) (string, error) {
	repeatRule := strings.Split(repeat, " ")
	if len(repeatRule) == 1 {
		return "", fmt.Errorf("не указан интервал в днях")
	}

	newDay, _ := strconv.Atoi(repeatRule[1])
	if t.Year()%4 == 0 {
		if newDay > 366 {
			return "", fmt.Errorf("превышен максимально допустимый интервал")
		}

	}

	if newDay > 365 {
		return "", fmt.Errorf("превышен максимально допустимый интервал")
	}

	t = t.AddDate(0, 0, newDay)
	return t.Format("20060102"), nil
}
