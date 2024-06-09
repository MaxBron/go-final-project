package handlers

import (
	"database/sql"
	"fmt"
	"go-final-project/nextdate"
	"go-final-project/task"
	"net/http"
	"sync"
	"time"
)

func DoneTaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		_, IDOk := r.URL.Query()["id"]
		if !IDOk {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "Не указан идентификатор"}`))
			return
		}

		str := r.URL.Query().Get("id")
		if str == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "Не указан идентификатор"}`))
			return
		}

		db, _ := sql.Open("sqlite", "scheduler.db")
		task := task.Task{}
		defer db.Close()
		var wg sync.WaitGroup
		wg.Add(1)
		var existsOk bool
		go func() {
			exists, _ := db.Query("SELECT EXISTS (SELECT * FROM SCHEDULER WHERE id = :id)", sql.Named("id", str))
			exists.Next()
			exists.Scan(&existsOk)
			defer exists.Close()
			wg.Done()
		}()

		wg.Wait()
		if !existsOk {
			fmt.Println()
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "Задача не найдена"}`))
			return
		}

		wg.Add(1)
		go func() {
			row, err := db.Query("SELECT * FROM SCHEDULER WHERE id = :id", sql.Named("id", str))
			if err != nil {
				w.Header().Set("Content-Type", "applictaion/json")
				w.WriteHeader(http.StatusBadGateway)
				w.Write([]byte(`{"error":"ошибка поиска"}`))
				return
			}

			row.Next()
			err = row.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
			defer row.Close()
			if err != nil {
				w.Header().Set("Content-Type", "applictaion/json")
				w.WriteHeader(http.StatusBadGateway)
				w.Write([]byte(`{"error":"ошибка считывания"}`))
				return
			}

			wg.Done()
		}()

		wg.Wait()
		if task.Repeat == "" {
			_, err := db.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", str))
			if err != nil {
				fmt.Println(err)
				w.Header().Set("Content-Type", "applictaion/json")
				w.WriteHeader(http.StatusBadGateway)
				w.Write([]byte(`{"error":"ошибка повторения"}`))
				return
			}

		} else {
			task.Date, _ = nextdate.NextDate(time.Now(), task.Date, task.Repeat)
			date, _ := time.Parse("20060102", task.Date)
			res, err := db.Exec("UPDATE SCHEDULER SET date = :date WHERE id = :id", sql.Named("date", date.Format("20060102")), sql.Named("id", task.Id))
			if err != nil {
				fmt.Println(err)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"error": "Задача не найдена"}`))
				return
			}

			fmt.Printf("res: %v\n", res)
		}

		w.Header().Set("Content-Type", "applictaion/json")
		w.Write([]byte(`{}`))
	}

}
