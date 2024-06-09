package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"

	"go-final-project/nextdate"
	"go-final-project/task"
	"net/http"
	"strconv"
	"time"
)

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var task task.Task
		var buf bytes.Buffer
		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"ошибка десериализации JSON"}`))
			return
		}

		if task.Title == "" {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"не указан заголовок задачи"}`))
			return
		}

		if task.Date == "" {
			task.Date = time.Now().Format("20060102")
		}

		if dateFormat, _ := time.Parse("20060102", task.Date); task.Date != dateFormat.Format("20060102") {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"дата представлена в формате, отличном от 20060102"}`))
			return
		}

		if _, err := nextdate.NextDate(time.Now(), task.Date, task.Repeat); err != nil {
			w.Header().Set("Content-Type", "applictaion/json; charset=UTF-8")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"правило повторения указано в неправильном формате"}`))
			return
		}

		if task.Date != time.Now().Format("20060102") {
			task.Date, _ = nextdate.NextDate(time.Now(), task.Date, task.Repeat)
		}

		if task.Date == "" {
			task.Date = time.Now().Format("20060102")
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write([]byte(fmt.Sprintf(`{"id":"%d"}`, task.AddTask())))
	case "GET":
		db, _ := sql.Open("sqlite", "scheduler.db")
		defer db.Close()
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

		id, err := strconv.Atoi(str)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "Задача не найдена"}`))
			return
		}

		exists, _ := db.Query("SELECT EXISTS (SELECT * FROM SCHEDULER WHERE id = :id)", sql.Named("id", id))
		exists.Next()
		var existsOk bool
		exists.Scan(&existsOk)
		defer exists.Close()
		if !existsOk {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "Задача не найдена"}`))
			return
		}

		task := task.Task{}
		row, _ := db.Query("SELECT * FROM SCHEDULER WHERE id = :id", sql.Named("id", id))
		defer row.Close()
		row.Next()
		row.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		taskJson, _ := json.Marshal(task)
		w.Header().Set("Content-Type", "application/json")
		w.Write(taskJson)
		return
	case "PUT":
		var task task.Task
		var buf bytes.Buffer
		buf.ReadFrom(r.Body)
		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"ошибка десериализации JSON"}`))
			return
		}

		if task.Title == "" {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"не указан заголовок задачи"}`))
			return
		}

		if task.Date == "" {
			task.Date = time.Now().Format("20060102")
		}

		if dateFormat, _ := time.Parse("20060102", task.Date); task.Date != dateFormat.Format("20060102") {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"дата представлена в формате, отличном от 20060102"}`))
			return
		}

		if _, err := nextdate.NextDate(time.Now(), task.Date, task.Repeat); err != nil {
			w.Header().Set("Content-Type", "applictaion/json; charset=UTF-8")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"правило повторения указано в неправильном формате"}`))
			return
		}

		if task.Date != time.Now().Format("20060102") {
			task.Date, _ = nextdate.NextDate(time.Now(), task.Date, task.Repeat)
		}

		db, _ := sql.Open("sqlite", "scheduler.db")
		res, err := db.Exec("UPDATE SCHEDULER SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id", sql.Named("date", task.Date),
			sql.Named("title", task.Title), sql.Named("comment", task.Comment), sql.Named("repeat", task.Repeat), sql.Named("id", task.Id))
		if err != nil {
			fmt.Println(err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "Задача не найдена"}`))
			return
		}

		exists, _ := db.Query("SELECT EXISTS (SELECT * FROM SCHEDULER WHERE id = :id)", sql.Named("id", task.Id))
		exists.Next()
		var existsOk bool
		exists.Scan(&existsOk)
		defer exists.Close()
		if !existsOk {
			fmt.Println()
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "Задача не найдена"}`))
			return
		}

		fmt.Printf("res: %v\n", res)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{}`))
		defer db.Close()
	case "DELETE":
		_, IDOk := r.URL.Query()["id"]
		if !IDOk {
			w.Header().Set("Content-Type", "application/json")

			w.Write([]byte(`{"error": "Не указан идентификатор"}`))
			return
		}

		str := r.URL.Query().Get("id")
		if str == "" {
			w.Header().Set("Content-Type", "application/json")

			w.Write([]byte(`{"error": "Не указан идентификатор"}`))
			return
		}

		db, _ := sql.Open("sqlite", "scheduler.db")
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
			fmt.Println("ffff")
			w.Header().Set("Content-Type", "application/json")

			w.Write([]byte(`{"error": "Задача не найдена"}`))
			return
		}

		_, err := db.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", str))
		if err != nil {
			fmt.Println(err)
			w.Header().Set("Content-Type", "applictaion/json")

			w.Write([]byte(`{"error":"ошибка"}`))
			return
		}

		w.Header().Set("Content-Type", "applictaion/json")
		fmt.Println("fffa")
		w.Write([]byte(`{}`))
		return
	}

}
