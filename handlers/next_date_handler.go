package handlers

import (
	"go-final-project/nextdate"
	"net/http"
	"time"
)

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query()["now"][0] == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Пустая строка"))
		return
	}

	if r.URL.Query()["date"][0] == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Пустая строка"))
		return
	}

	now, _ := time.Parse("20060102", r.URL.Query()["now"][0])
	date := r.URL.Query()["date"][0]
	repeat := r.URL.Query()["repeat"][0]
	nextDate, err := nextdate.NextDate(now, date, repeat)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nextDate))
}
