package main

import (
	"go-final-project/handlers"
	"net/http"
)

func main() {
	http.Get("/api/nextdate")
	http.Handle("/api/task/done", http.HandlerFunc(handlers.DoneTaskHandler))
	http.Handle("/api/tasks", http.HandlerFunc(handlers.TasksHandler))
	http.Handle("/api/task", http.HandlerFunc(handlers.TaskHandler))
	http.Handle("/api/nextdate", http.HandlerFunc(handlers.NextDateHandler))
	webDir := "web"
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	http.ListenAndServe(":7540", nil)
}
