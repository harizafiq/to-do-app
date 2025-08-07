package main

import (
	"log"
	"net/http"
	"toDoList/handler"
)

func main() {
	http.HandleFunc("/todos", handler.HandleTodos)
	http.HandleFunc("/todos/", handler.HandleTodos)
	http.HandleFunc("/todos/done/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPatch:
			handler.UpdateDone(w, r)
		}
	})
	log.Println("Server running :8080")
	http.ListenAndServe(":8080", nil)
}
