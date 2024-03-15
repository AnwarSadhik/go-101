package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Todo struct {
	Id        string `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var todos []Todo
var indexPage = `
  <html>
  <head>
    <title>index</title>
  </head>
  <body>
    <h1>Hello World!</h1>
  </body>
  </html>
`

func getTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func addTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var todo Todo
	_ = json.NewDecoder(r.Body).Decode(&todo)
	todo.Id = fmt.Sprintf("%d", len(todos)+1)
	todos = append(todos, todo)
	json.NewEncoder(w).Encode(todo)

}
func updateTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for index, todo := range todos {
		if todo.Id == params["id"] {
			todos = append(todos[:index], todos[index+1:]...)
			var todo Todo
			todo.Id = params["id"]
			_ = json.NewDecoder(r.Body).Decode(&todo)
			todos = append(todos, todo)
			json.NewEncoder(w).Encode("todo updated successfully")
			return
		}
	}

	json.NewEncoder(w).Encode("Todo not found")

}
func deleteTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for index, todo := range todos {
		if todo.Id == params["id"] {
			todos = append(todos[:index], todos[index+1:]...)
			json.NewEncoder(w).Encode("Todo deleted successfully")
		}
	}

	json.NewEncoder(w).Encode("Todo not found")
}

func serveIndexPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, indexPage)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", serveIndexPage)
	r.HandleFunc("/todos", getTodos).Methods("GET")
	r.HandleFunc("/todos", addTodo).Methods("POST")
	r.HandleFunc("/todos/{id}", updateTodo).Methods("PUT")
	r.HandleFunc("/todos/{id}", deleteTodo).Methods("DELETE")

	fmt.Printf("Server started on port 8080\n")
	log.Fatal(http.ListenAndServe(":8080", r))
}
