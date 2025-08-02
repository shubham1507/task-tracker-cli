// task-cli.go
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

const taskFile = "tasks.json"

func loadTasks() ([]Task, error) {
	if _, err := os.Stat(taskFile); os.IsNotExist(err) {
		return []Task{}, nil
	}
	data, err := ioutil.ReadFile(taskFile)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	err = json.Unmarshal(data, &tasks)
	return tasks, err
}

func saveTasks(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(taskFile, data, 0644)
}

func nextID(tasks []Task) int {
	maxID := 0
	for _, task := range tasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}
	return maxID + 1
}

func findTask(tasks []Task, id int) (int, *Task) {
	for i, task := range tasks {
		if task.ID == id {
			return i, &tasks[i]
		}
	}
	return -1, nil
}

func addTask(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: task-cli add \"task description\"")
		return
	}
	description := args[1]
	tasks, _ := loadTasks()
	task := Task{
		ID:          nextID(tasks),
		Description: description,
		Status:      "todo",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	tasks = append(tasks, task)
	saveTasks(tasks)
	fmt.Printf("Task added successfully (ID: %d)\n", task.ID)
}

func updateTask(args []string) {
	if len(args) < 3 {
		fmt.Println("Usage: task-cli update <id> \"new description\"")
		return
	}
	id, _ := strconv.Atoi(args[1])
	tasks, _ := loadTasks()
	index, task := findTask(tasks, id)
	if task == nil {
		fmt.Println("Task not found")
		return
	}
	task.Description = args[2]
	task.UpdatedAt = time.Now()
	tasks[index] = *task
	saveTasks(tasks)
	fmt.Println("Task updated successfully")
}

func deleteTask(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: task-cli delete <id>")
		return
	}
	id, _ := strconv.Atoi(args[1])
	tasks, _ := loadTasks()
	index, _ := findTask(tasks, id)
	if index == -1 {
		fmt.Println("Task not found")
		return
	}
	tasks = append(tasks[:index], tasks[index+1:]...)
	saveTasks(tasks)
	fmt.Println("Task deleted successfully")
}

func markTask(args []string, status string) {
	if len(args) < 2 {
		fmt.Printf("Usage: task-cli mark-%s <id>\n", status)
		return
	}
	id, _ := strconv.Atoi(args[1])
	tasks, _ := loadTasks()
	index, task := findTask(tasks, id)
	if task == nil {
		fmt.Println("Task not found")
		return
	}
	task.Status = status
	task.UpdatedAt = time.Now()
	tasks[index] = *task
	saveTasks(tasks)
	fmt.Printf("Task marked as %s\n", status)
}

func listTasks(args []string) {
	tasks, _ := loadTasks()
	filter := ""
	if len(args) > 1 {
		filter = args[1]
	}
	for _, task := range tasks {
		if filter == "" || task.Status == filter {
			fmt.Printf("[%d] %s (%s)\nCreated: %s\nUpdated: %s\n\n", task.ID, task.Description, task.Status, task.CreatedAt.Format(time.RFC822), task.UpdatedAt.Format(time.RFC822))
		}
	}
}

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Usage: task-cli <command> [arguments]")
		return
	}
	switch args[1] {
	case "add":
		addTask(args[1:])
	case "update":
		updateTask(args[1:])
	case "delete":
		deleteTask(args[1:])
	case "mark-in-progress":
		markTask(args[1:], "in-progress")
	case "mark-done":
		markTask(args[1:], "done")
	case "list":
		listTasks(args[1:])
	default:
		fmt.Println("Unknown command")
	}
}
