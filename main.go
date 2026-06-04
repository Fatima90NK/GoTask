// add main package which is the entry point of the application
package main

//import necessary packages for the application
import (
	"context"       // context package for managing timeouts and cancellation of operations (MongoDB connection)
	"encoding/json" // encoding/json package for encoding and decoding JSON data (convert Go structs to JSON and vice versa)
	"fmt"           // fmt package for formatted I/O operations (printing error messages to the console)
	"net/http"      // net/http package for creating an HTTP server and handling HTTP requests
	"os"            // os package for accessing environment variables (getting MongoDB connection string)
	"time"          // time package for setting timeouts (MongoDB connection)

	// strconv package for converting strings to integers (parsing task ID from query parameters)
	"go.mongodb.org/mongo-driver/bson"          // bson package for working with BSON data (MongoDB documents)
	"go.mongodb.org/mongo-driver/mongo"         // mongo package for interacting with MongoDB (connecting to the database and performing operations)
	"go.mongodb.org/mongo-driver/mongo/options" // options package for configuring MongoDB client options (setting the connection URI)
)

// define a Task struct to represent a task with an ID, name, and done status
// task is a single to-do item (data-model)
// json tags make it readable in API responses
type Task struct {
	ID   int    `bson:"id"`
	Name string `bson:"name"`
	Done bool   `bson:"done"`
}

// client and collection variables for MongoDB connection
var client *mongo.Client
var collection *mongo.Collection
var nextID int = 1 // variable to keep track of the next task ID

func init() {
	//get connection string from environment variable
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		fmt.Println("Error: MONGO_URI environment variable not set")
		os.Exit(1)
	}

	// Initialize MongoDB client
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//connect to MongoDB using the connection string
	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		fmt.Println("Error connecting to MongoDB:", err)
		os.Exit(1)
	}

	//Verify the connection to MongoDB
	if err := client.Ping(ctx, nil); err != nil {
		fmt.Println("Error pinging MongoDB:", err)
		os.Exit(1)
	}

	fmt.Println("Connected to MongoDB successfully!")

	//get a reference to the "tasks" collection in the "gotask" database
	collection = client.Database("gotask").Collection("tasks")
}

func main() {
	defer client.Disconnect(context.Background())

	http.HandleFunc("/tasks", handleTasks)   // Handle requests to /tasks endpoint
	http.HandleFunc("/tasks/done", markDone) // Handle requests to /tasks/done endpoint

	fmt.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func handleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getTasks(w, r) // Handle GET requests to retrieve all tasks
	case "POST":
		createTask(w, r) // Handle POST requests to create a new task
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed) // Return 405 for unsupported methods
	}
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{}) // Find all tasks in the collection
	if err != nil {
		http.Error(w, "Error fetching tasks", http.StatusInternalServerError)
		fmt.Println("Error fetching tasks from MongoDB:", err)
		return
	}
	defer cursor.Close(ctx)

	var tasks []Task
	if err = cursor.All(ctx, &tasks); err != nil { // Decode the cursor into a slice of Task structs
		http.Error(w, "Error decoding tasks", http.StatusInternalServerError)
		fmt.Println("Error decoding tasks:", err)
		return
	}

	//if no tasks are found, return an empty array instead of null
	if tasks == nil {
		tasks = []Task{}
	}

	json.NewEncoder(w).Encode(tasks) // Encode the tasks slice to JSON and write it to the response
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil { // Decode the request body into a Task struct
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	task.ID = nextID // Assign a unique ID to the new task
	nextID++

	//insert the new task into the MongoDB collection
	result, err := collection.InsertOne(ctx, task) // Insert the new task into the MongoDB collection
	if err != nil {
		http.Error(w, "Error creating task", http.StatusInternalServerError)
		fmt.Println("Error creating task:", err)
		return
	}

	fmt.Println("Task created with ID:", result.InsertedID)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func markDone(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	var input struct {
		ID int `json:"id"` // Define a struct to parse the task ID from the request body
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil { // Decode the request body into the input struct
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Update the task's done status to true in the MongoDB collection
	result := collection.FindOneAndUpdate(
		ctx,
		bson.M{"id": input.ID},
		bson.M{"$set": bson.M{"done": true}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	var updatedTask Task
	if err := result.Decode(&updatedTask); err != nil { // Decode the updated task from the result
		http.Error(w, "Error updating task", http.StatusInternalServerError)
		fmt.Println("Error updating task:", err)
		return
	}

	json.NewEncoder(w).Encode(updatedTask) // Encode the updated task to JSON and write it to the response
}
