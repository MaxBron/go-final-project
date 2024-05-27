package main

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type Task struct {
	Id      string `json:"id,omitempty"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

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

func (task Task) UpdateTask() bool {
	db, _ := sql.Open("sqlite", "scheduler.db")
	exists, _ := db.Query("SELECT EXISTS (SELECT * FROM SCHEDULER WHERE id = :id)", sql.Named("id", task.Id))
	exists.Next()
	var existsOk bool
	exists.Scan(&existsOk)
	defer exists.Close()
	if !existsOk {
		return false
	}

	res, err := db.Exec("UPDATE SCHEDULER SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id", sql.Named("date", task.Date),
		sql.Named("title", task.Title), sql.Named("comment", task.Comment), sql.Named("repeat", task.Repeat), sql.Named("id", task.Id))
	defer db.Close()
	fmt.Printf("res: %v\n", res)
	return err != nil
}
