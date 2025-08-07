package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

var todos []Todo = []Todo{
	{ID: 1, Text: "Learn Go", Done: false},
	{ID: 2, Text: "Learn Svelte", Done: false},
}

type Todo struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

func HandleTodos(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if path == "/todos" {
		switch r.Method {
		case http.MethodGet:
			GetTodos(w, r)
		case http.MethodPost:
			AddTodos(w, r)
		default:
			// http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Method not allowed",
			})
		}
		return
	}

	//example: /todos/123
	idStr := strings.TrimPrefix(r.URL.Path, "/todos/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Invalid id",
		})
		return
	}

	switch r.Method {
	case http.MethodGet:
		GetTodo(w, r, id)
	case http.MethodPut:
		UpdateTodo(w, r, id)
	case http.MethodDelete:
		DeleteTodo(w, r, id)
	default:
		// http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Method not allowed",
		})
	}
}

// get list
func GetTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Success Get Record",
		"data":    todos,
	})
}

// add list
func AddTodos(w http.ResponseWriter, r *http.Request) {
	var todo Todo

	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		// http.Error(w, "Invalid JSON", http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Invalid JSON",
		})
		return
	}

	//set id
	todo.ID = len(todos) + 1

	//append to list
	todos = append(todos, todo)

	//responsd with JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Created Succesfully",
		"data":    todos,
	})
}

// GET
func GetTodo(w http.ResponseWriter, r *http.Request, id int) {

	var todo Todo
	found := false
	for i := range todos {
		if todos[i].ID == id {
			todo = todos[i]
			found = true
			break
		}
	}

	if !found {
		// http.Error(w, "Record Not Found", http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Record Not found",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Success retrieved",
		"data":    todo,
	})
}

func DeleteTodo(w http.ResponseWriter, r *http.Request, id int) {
	newList := []Todo{}
	found := false

	for _, todo := range todos {
		if todo.ID != id {
			newList = append(newList, todo)
		} else {
			found = true
		}
	}
	todos = newList

	if !found {
		//good for debugging
		// http.Error(w, "Record Not found", http.StatusNotFound)

		//standard API Response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Todo Not Found",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Todo Deleted Successfully",
	})
}

func UpdateTodo(w http.ResponseWriter, r *http.Request, id int) {
	var updateTodo Todo

	// get changes value
	if err := json.NewDecoder(r.Body).Decode(&updateTodo); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Json is invalid",
		})
		return
	}

	found := false
	for i, todo := range todos {
		if todo.ID == id {
			todos[i].Text = updateTodo.Text
			updateTodo = todos[i]
			found = true
			break
		}
	}

	if !found {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Record Not found",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Record Updated",
		"data":    updateTodo,
		// "data":    1,
	})
}

func UpdateDone(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/todos/done/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Invalid Id",
		})
		return
	}

	found := false
	for i := range todos {
		if todos[i].ID == id {
			todos[i].Done = true
			found = true
			break
		}
	}

	if !found {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Record not found",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Record Updated",
		"data":    1,
	})
}
