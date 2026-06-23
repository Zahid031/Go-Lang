package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"` // "todo", "in-progress", "done"
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

const taskFile = "tasks.json"

// loadTasks reads tasks.json and returns the list of tasks.
// If the file doesn't exist, it creates an empty one and returns an empty slice.
func loadTasks() ([]Task, error) {
	if _, err := os.Stat(taskFile); os.IsNotExist(err) {
		empty := []Task{}
		if err := saveTasks(empty); err != nil {
			return nil, fmt.Errorf("failed to create task file: %w", err)
		}
		return empty, nil
	}

	data, err := os.ReadFile(taskFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read task file: %w", err)
	}

	if len(data) == 0 {
		return []Task{}, nil
	}

	var tasks []Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, fmt.Errorf("failed to parse task file (is it corrupted?): %w", err)
	}

	return tasks, nil
}

// saveTasks writes the given task list to tasks.json, pretty-printed.
func saveTasks(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode tasks: %w", err)
	}

	if err := os.WriteFile(taskFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write task file: %w", err)
	}

	return nil
}

// nextID returns the next available task ID (highest existing ID + 1).
func nextID(tasks []Task) int {
	maxID := 0
	for _, t := range tasks {
		if t.ID > maxID {
			maxID = t.ID
		}
	}
	return maxID + 1
}

// findTaskIndex returns the index of the task with the given ID, or -1 if not found.
func findTaskIndex(tasks []Task, id int) int {
	for i, t := range tasks {
		if t.ID == id {
			return i
		}
	}
	return -1
}

func addTask(description string) error {
	tasks, err := loadTasks()
	if err != nil {
		return err
	}

	now := time.Now()
	newTask := Task{
		ID:          nextID(tasks),
		Description: description,
		Status:      "todo",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	tasks = append(tasks, newTask)

	if err := saveTasks(tasks); err != nil {
		return err
	}

	fmt.Printf("Task added successfully (ID: %d)\n", newTask.ID)
	return nil
}

func updateTask(id int, newDescription string) error {
	tasks, err := loadTasks()
	if err != nil {
		return err
	}

	idx := findTaskIndex(tasks, id)
	if idx == -1 {
		return fmt.Errorf("task with ID %d not found", id)
	}

	tasks[idx].Description = newDescription
	tasks[idx].UpdatedAt = time.Now()

	if err := saveTasks(tasks); err != nil {
		return err
	}

	fmt.Printf("Task %d updated successfully\n", id)
	return nil
}

func deleteTask(id int) error {
	tasks, err := loadTasks()
	if err != nil {
		return err
	}

	idx := findTaskIndex(tasks, id)
	if idx == -1 {
		return fmt.Errorf("task with ID %d not found", id)
	}

	// Remove the element at idx by combining the slice before and after it.
	tasks = append(tasks[:idx], tasks[idx+1:]...)

	if err := saveTasks(tasks); err != nil {
		return err
	}

	fmt.Printf("Task %d deleted successfully\n", id)
	return nil
}

func markTaskStatus(id int, status string) error {
	tasks, err := loadTasks()
	if err != nil {
		return err
	}

	idx := findTaskIndex(tasks, id)
	if idx == -1 {
		return fmt.Errorf("task with ID %d not found", id)
	}

	tasks[idx].Status = status
	tasks[idx].UpdatedAt = time.Now()

	if err := saveTasks(tasks); err != nil {
		return err
	}

	fmt.Printf("Task %d marked as %s\n", id, status)
	return nil
}

func printUsage() {
	fmt.Println(`Usage: task-tracker <command> [arguments]

Commands:
  add <description>             Add a new task
  update <id> <description>     Update a task's description
  delete <id>                   Delete a task
  mark-in-progress <id>         Mark a task as in progress
  mark-done <id>                Mark a task as done`)
}

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		printUsage()
		os.Exit(1)
	}

	command := args[0]

	switch command {
	case "add":
		if len(args) < 2 {
			fmt.Println("Error: 'add' requires a task description")
			fmt.Println(`Usage: task-tracker add "task description"`)
			os.Exit(1)
		}
		if err := addTask(args[1]); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

	case "update":
		if len(args) < 3 {
			fmt.Println("Error: 'update' requires an ID and a new description")
			fmt.Println(`Usage: task-tracker update <id> "new description"`)
			os.Exit(1)
		}
		id, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Error: task ID must be a number")
			os.Exit(1)
		}
		if err := updateTask(id, args[2]); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

	case "delete":
		if len(args) < 2 {
			fmt.Println("Error: 'delete' requires an ID")
			fmt.Println("Usage: task-tracker delete <id>")
			os.Exit(1)
		}
		id, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Error: task ID must be a number")
			os.Exit(1)
		}
		if err := deleteTask(id); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

	case "mark-in-progress":
		if len(args) < 2 {
			fmt.Println("Error: 'mark-in-progress' requires an ID")
			os.Exit(1)
		}
		id, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Error: task ID must be a number")
			os.Exit(1)
		}
		if err := markTaskStatus(id, "in-progress"); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

	case "mark-done":
		if len(args) < 2 {
			fmt.Println("Error: 'mark-done' requires an ID")
			os.Exit(1)
		}
		id, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Error: task ID must be a number")
			os.Exit(1)
		}
		if err := markTaskStatus(id, "done"); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

	default:
		fmt.Printf("Error: unknown command '%s'\n", command)
		printUsage()
		os.Exit(1)
	}
}