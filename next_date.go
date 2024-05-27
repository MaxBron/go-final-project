package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func NextYear(now time.Time, t time.Time, repeat string) (string, error) {
	t = t.AddDate(1, 0, 0)
	return t.Format("20060102"), nil
}

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

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return time.Now().Format("20060102"), nil
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
