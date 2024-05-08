package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tabintel/local-first-todo/pkg/database"
)

func main() {
	// Initialize database connection
	err := database.Init("tasks.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	// Initialize database table
	err = database.CreateTaskTable()
	if err != nil {
		log.Fatal(err)
	}

	// Create router
	router := mux.NewRouter()
	router.HandleFunc("/api/tasks", getTasksHandler).Methods("GET")
	router.HandleFunc("/api/tasks", addTaskHandler).Methods("POST")

	// Start HTTP server
	log.Println("Server listening on port 8000... Built with 💜, EkeminiOS")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := database.GetTasks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task database.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.AddTask(task.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
