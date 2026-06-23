package main

import (
	"encoding/json"
	"fmt"
	"os"
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
	// Check if file exists
	if _, err := os.Stat(taskFile); os.IsNotExist(err) {
		// File doesn't exist yet — create it with an empty array
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

	// Handle the edge case of an empty file (0 bytes)
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

func main() {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Printf("Loaded %d task(s)\n", len(tasks))
}