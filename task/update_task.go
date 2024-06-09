package task

import (
	"database/sql"
	"fmt"
)

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
