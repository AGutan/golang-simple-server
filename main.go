package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Todo struct {
	Todo string `json:"todo"`
}

var todos []*Todo

func main() {
	router := mux.NewRouter()
	// POST /todo
	router.HandleFunc("/todo", newTodo).Methods(http.MethodPost)
	// GET /todo
	router.HandleFunc("/todo", getTodos).Methods(http.MethodGet)
	// PUT /todo/complete
	router.HandleFunc("/todo/complete", completeTodo).Methods(http.MethodPut)

	// CORS
	headersOK := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	originsOK := handlers.AllowedOrigins([]string{"*"})
	methodsOK := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "DELETE", "PUT"})

	// setup http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", 8080),
		Handler: handlers.CORS(headersOK, originsOK, methodsOK)(router),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalf("Could not start server")
	}
}

func getTodos(w http.ResponseWriter, req *http.Request) {
	// write success response
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(todos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func completeTodo(w http.ResponseWriter, req *http.Request) {
	var todo Todo
	err := json.NewDecoder(req.Body).Decode(&todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	i := 0
	newTodos := make([]*Todo, len(todos))
	for _, t := range todos {
		if t.Todo != todo.Todo {
			newTodos[i] = t
			i++
		}
	}
	todos = newTodos[:i]

	// write success response
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(todos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func newTodo(w http.ResponseWriter, req *http.Request) {
	var todo Todo
	err := json.NewDecoder(req.Body).Decode(&todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	todos = append(todos, &todo)
	log.Printf("todo %v", todo)
	// close body to avoid memory leak
	err = req.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// write success response
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
