package nextdate

import (
	"fmt"
	"time"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", nil
	}

	if repeat[0] == 'y' {
		t, err := time.Parse("20060102", date)
		if err != nil {
			return "", fmt.Errorf(`время в переменной date не может быть преобразовано в корректную дату — ошибка выполнения time.Parse("20060102", d)`)
		}

		oldT := t
		for (now.After(t) || oldT.After(t)) || now == t || oldT == t {
			date, err = NextYear(now, t, repeat)
			if err != nil {
				return "", err
			}

			t, _ = time.Parse("20060102", date)
		}

		return t.Format("20060102"), err
	}

	if repeat[0] == 'd' {
		t, err := time.Parse("20060102", date)
		if err != nil {
			return "", fmt.Errorf(`время в переменной date не может быть преобразовано в корректную дату — ошибка выполнения time.Parse("20060102", d)`)
		}

		oldT := t
		for (now.After(t) || oldT.After(t)) || now == t || oldT == t {
			date, err = NextDay(now, t, repeat)
			if err != nil {
				return "", err
			}

			t, _ = time.Parse("20060102", date)
		}

		return t.Format("20060102"), err
	}

	return "", fmt.Errorf("недопустимый символ")
}
