package main

//Run with CompileDaemon
//$HOME/go/bin/CompileDaemon -command="./go_api"

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"io/ioutil"
	"strconv"

	"github.com/gorilla/mux"
)

type task struct {
	ID      int    `json:"id"`
	Name    string `json:"Name"`
	Content string `json:"content"`
}

type allTasks []task

var tasks = allTasks{
	{
		ID:      1,
		Name:    "Task 1",
		Content: "This is task 1",
	},
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: getTasks")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func getTaskByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskId, err := strconv.Atoi(vars["id"])
	
	if err != nil {
		fmt.Fprintf(w, "Ivalid ID")
		return
	}
	
	for _, task := range tasks {
		if task.ID == taskId {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(task)
		}
	}
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var newTask task
	reqBody, err := ioutil.ReadAll(r.Body)
	
	if err != nil {
		fmt.Fprintf(w, "Insert a valid task")
	}
	
	json.Unmarshal(reqBody, &newTask)
	
	newTask.ID = len(tasks) +  1
	tasks = append(tasks, newTask)
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskId, err := strconv.Atoi(vars["id"])
	reqBody, err := ioutil.ReadAll(r.Body)
	var updatedTask task
	
	if err != nil {
		fmt.Fprintf(w, "Insert a valid task")
	}
	
	json.Unmarshal(reqBody, &updatedTask)
	
	if err != nil {
		fmt.Fprintf(w, "Ivalid ID")
		return
	}
	
	for i, task := range tasks {
		if task.ID == taskId {
			tasks = append(tasks[:i], tasks[i + 1:]...)
			updatedTask.ID = taskId
			tasks = append(tasks, updatedTask)
			
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(updatedTask)
		}
	}
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskId, err := strconv.Atoi(vars["id"])
	
	if err != nil {
		fmt.Fprintf(w, "Ivalid ID")
		return
	}
	
	for i, task := range tasks {
		if task.ID == taskId {
			tasks = append(tasks[:i], tasks[i + 1:]...)
			fmt.Fprintf(w, "Tasks %v deleted", taskId)
		}
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World! X")
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", index)
	router.HandleFunc("/tasks", getTasks)
	router.HandleFunc("/tasks/{id}", getTaskByID)
	router.HandleFunc("/createTask", createTask).Methods("POST")
	router.HandleFunc("/updateTask/{id}", updateTask).Methods("PUT")
	router.HandleFunc("/deleteTask/{id}", deleteTask).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":3000", router))
}
