package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
)

var (
	ErrFileIsNil    = errors.New("file is nil")
	ErrInvalidArgs  = errors.New("invalid arguments")
	ErrTaskNotFound = errors.New("task not found")
)

const (
	DataFileName = "tasks.json"
	FilePerm     = 0644
)

type Task struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

type TaskManager struct {
	tasks    []Task
	nextID   int
	fileName string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: taskmanager <command> [arguments]")
		fmt.Println("Commands: add <name> <text>, read, update <id>, delete <id>")
		return
	}

	tm := NewTaskManager(DataFileName)

	if err := tm.Load(); err != nil {
		log.Fatalf("Failed to load tasks: %v", err)
	}

	if err := handleCommand(tm, os.Args[1], os.Args[2:]); err != nil {
		log.Printf("Command failed: %v", err)
	}

	if err := tm.Save(); err != nil {
		log.Printf("Failed to save tasks: %v", err)
	}
}

func NewTaskManager(fileName string) *TaskManager {
	return &TaskManager{
		tasks:    make([]Task, 0),
		nextID:   1,
		fileName: fileName,
	}
}

func handleCommand(tm *TaskManager, command string, args []string) error {
	switch command {
	case "add":
		if len(args) < 2 {
			return fmt.Errorf("%w: add requires name and text", ErrInvalidArgs)
		}
		return tm.AddTask(args[0], args[1])
	case "read":
		return tm.ListTasks()
	case "update":
		if len(args) < 1 {
			return fmt.Errorf("%w: update requires task ID", ErrInvalidArgs)
		}
		var id int
		if _, err := fmt.Sscanf(args[0], "%d", &id); err != nil {
			return fmt.Errorf("invalid task ID: %v", err)
		}
		return tm.UpdateTaskStatus(id)
	case "delete":
		if len(args) < 1 {
			return fmt.Errorf("%w: delete requires task ID", ErrInvalidArgs)
		}
		var id int
		if _, err := fmt.Sscanf(args[0], "%d", &id); err != nil {
			return fmt.Errorf("invalid task ID: %v", err)
		}
		return tm.DeleteTask(id)
	default:
		return fmt.Errorf("unknown command: %s", command)
	}
}

func (tm *TaskManager) AddTask(name, text string) error {
	task := Task{
		ID:   tm.nextID,
		Name: name,
		Text: text,
		Done: false,
	}

	tm.tasks = append(tm.tasks, task)
	tm.nextID++

	fmt.Printf("Task added successfully (ID: %d)\n", task.ID)
	return nil
}

func (tm *TaskManager) ListTasks() error {
	if len(tm.tasks) == 0 {
		fmt.Println("No tasks found")
		return nil
	}

	fmt.Printf("%-4s %-20s %-30s %-6s\n", "ID", "Name", "Text", "Done")
	fmt.Println("------------------------------------------------------------")
	for _, task := range tm.tasks {
		status := "No"
		if task.Done {
			status = "Yes"
		}
		fmt.Printf("%-4d %-20s %-30s %-6s\n", task.ID, task.Name, task.Text, status)
	}
	return nil
}

func (tm *TaskManager) UpdateTaskStatus(id int) error {
	for i := range tm.tasks {
		if tm.tasks[i].ID == id {
			tm.tasks[i].Done = !tm.tasks[i].Done
			status := "completed"
			if !tm.tasks[i].Done {
				status = "pending"
			}
			fmt.Printf("Task %d marked as %s\n", id, status)
			return nil
		}
	}
	return fmt.Errorf("%w: task with ID %d not found", ErrTaskNotFound, id)
}

func (tm *TaskManager) DeleteTask(id int) error {
	for i, task := range tm.tasks {
		if task.ID == id {
			tm.tasks = append(tm.tasks[:i], tm.tasks[i+1:]...)
			fmt.Printf("Task %d deleted successfully\n", id)
			return nil
		}
	}
	return fmt.Errorf("%w: task with ID %d not found", ErrTaskNotFound, id)
}

func (tm *TaskManager) Save() error {
	if len(tm.tasks) == 0 {
		tm.tasks = []Task{}
	}

	jsonData, err := json.MarshalIndent(tm.tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal tasks to JSON: %w", err)
	}

	if err := os.WriteFile(tm.fileName, jsonData, FilePerm); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Printf("Tasks saved successfully (%d tasks)\n", len(tm.tasks))
	return nil
}

func (tm *TaskManager) Load() error {
	data, err := os.ReadFile(tm.fileName)
	if err != nil {
		if os.IsNotExist(err) {
			tm.tasks = []Task{}
			tm.nextID = 1
			fmt.Println("New task file created")
			return nil
		}
		return fmt.Errorf("failed to read file: %w", err)
	}

	if len(data) == 0 {
		tm.tasks = []Task{}
		tm.nextID = 1
		return nil
	}

	var tasks []Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	tm.tasks = tasks

	maxID := 0
	for _, task := range tasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}
	tm.nextID = maxID + 1

	fmt.Printf("Loaded %d tasks\n", len(tasks))
	return nil
}
