package task

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func (task Task) AddTask() int64 {
	db, _ := sql.Open("sqlite", "scheduler.db")
	res, err := db.Exec("INSERT INTO SCHEDULER (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", task.Date), sql.Named("title", task.Title), sql.Named("comment", task.Comment), sql.Named("repeat", task.Repeat))
	if err != nil {
		fmt.Println(err)
		return 1
	}

	defer db.Close()
	id, _ := res.LastInsertId()
	return id
}
