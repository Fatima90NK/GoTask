//add main package which is the entry point of the application
package main

//import necessary packages for the application
// fmt for printing error messages to the console
// encoding/json for encoding and decoding JSON data (convert go structs to JSON and vice versa)
// net/http for handling HTTP requests and responses (webserver)
import (
    "fmt"
	"encoding/json"
	"net/http"
)
//define a Task struct to represent a task with an ID, name, and done status
//task is a single to-do item (data-model)
//json tags make it readable in API responses
type Task struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Done       bool   `json:"done"`
}
//tasks slice to store all tasks in memory (not persistent, will be lost when server restarts)
//nextID variable to assign unique IDs to new tasks
var tasks []Task
var nextID int = 1
func main() {
    http.HandleFunc("/tasks", handleTasks)
	http.HandleFunc("/tasks/done", markDone)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Failed to start server:", err)
	}	
}
//handleTasks function to handle requests to the /tasks endpoint
//GET method returns the list of tasks as JSON
//POST method allows adding a new task by sending a JSON body with the task name
func handleTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		json.NewEncoder(w).Encode(tasks)
	case "POST":
		defer r.Body.Close()
		var task Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		task.ID = nextID
		nextID++
		tasks = append(tasks, task)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(task)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
//markDone marks a task as done based on the ID provided in the request body
func markDone(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	var input struct {ID int `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	//loops over all tasks
	//finds the task with the matching ID and marks it as done
	for i, task := range tasks {
		if task.ID == input.ID {
			tasks[i].Done = true
			json.NewEncoder(w).Encode(tasks[i])
			return
		}
	}
	http.Error(w, "Task not found", http.StatusNotFound)
}	
		DELETE /tasks?ID=1 - deletes the task with the specified ID (not implemented in this code)
		