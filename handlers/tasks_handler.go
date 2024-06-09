package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go-final-project/task"
	"log"
	"net/http"
)

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := sql.Open("sqlite", "scheduler.db")
	tasks := []task.Task{}

	rows, err := db.Query("SELECT * FROM SCHEDULER ORDER BY DATE LIMIT 10")
	if err != nil {
		w.Header().Set("Content-Type", "applictaion/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"текст ошибки"}`))
		return
	}

	defer rows.Close()
	for rows.Next() {
		task := task.Task{}
		err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			log.Println(err)
			return
		}

		tasks = append(tasks, task)
	}

	w.Header().Set("Content-Type", "application/json")
	m := make(map[string][]task.Task)
	m["tasks"] = tasks
	tasksJson, err := json.Marshal(m)
	fmt.Println(m)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "текст ошибки"}`))
		return
	}

	defer db.Close()

	w.Write(tasksJson)
}
