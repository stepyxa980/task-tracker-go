package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	DataFileName = "tasks.json"
	FilePerm     = 0644
)

var (
	// Основные ошибки
	ErrFileIsNil    = errors.New("file is nil")
	ErrInvalidArgs  = errors.New("invalid arguments")
	ErrTaskNotFound = errors.New("task not found")
)

type Task struct {
	// Task Structure
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Text       string `json:"text"`
	Done       bool   `json:"done"`
	InProgress bool   `json:"inProgress"`
}

type TaskManager struct {
	// Task management structure
	tasks    []Task
	nextID   int
	fileName string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: taskmanager <command> [arguments]")
		fmt.Println("Commands: add <name> <text>, read, update <id>, delete <id>, progress <id>, done <id>")
		fmt.Println("Additional commands: list-done, list-not-done, list-in-progress")
		return
	}

	tm := NewTaskManager(DataFileName)

	if err := tm.Load(); err != nil {
		log.Fatalf("Failed to load tasks: %v", err)
	}

	if err := CommandHandler(tm, os.Args[1], os.Args[2:]); err != nil {
		log.Printf("Command failed: %v", err)
	}

	if err := tm.Save(); err != nil {
		log.Printf("Failed to save tasks: %v", err)
	}
}

func NewTaskManager(newFileName string) *TaskManager {
	// Create and transfer a new TaskManager, where the first object is at ID: 1
	return &TaskManager{
		tasks:    make([]Task, 0),
		nextID:   1,
		fileName: newFileName,
	}
}

func CommandHandler(tm *TaskManager, command string, args []string) error {
	switch command {
	case "add":
		if len(args) < 2 {
			return fmt.Errorf("%w: add requires name and text", ErrInvalidArgs)
		}
		err := tm.AddTask(args[0], strings.Join(args[1:], " "))
		return err
	case "read":
		text, err := tm.TaskList()
		if err != nil {
			return err
		}
		fmt.Println(text)
		return nil
	case "update":
		if len(args) < 1 {
			return fmt.Errorf("%w: update requires task ID", ErrInvalidArgs)
		}

		id, err := strconv.Atoi(args[0])
		if err != nil || id <= 0 {
			return fmt.Errorf("%w: task ID must be a positive integer", ErrInvalidArgs)
		}

		return tm.UpdateTaskStatus(id)
	case "delete":
		if len(args) < 1 {
			return fmt.Errorf("%w: delete requires task ID", ErrInvalidArgs)
		}
		id, err := strconv.Atoi(args[0])
		if err != nil || id <= 0 {
			return fmt.Errorf("%w: task ID must be a positive integer", ErrInvalidArgs)
		}
		return tm.DeleteTask(id)
	// Новые команды
	case "progress":
		if len(args) < 1 {
			return fmt.Errorf("%w: progress requires task ID", ErrInvalidArgs)
		}
		id, err := strconv.Atoi(args[0])
		if err != nil || id <= 0 {
			return fmt.Errorf("%w: task ID must be a positive integer", ErrInvalidArgs)
		}
		return tm.MarkTaskInProgress(id)
	case "done":
		if len(args) < 1 {
			return fmt.Errorf("%w: done requires task ID", ErrInvalidArgs)
		}
		id, err := strconv.Atoi(args[0])
		if err != nil || id <= 0 {
			return fmt.Errorf("%w: task ID must be a positive integer", ErrInvalidArgs)
		}
		return tm.MarkTaskDone(id)
	case "list-done":
		text, err := tm.ListDoneTasks()
		if err != nil {
			return err
		}
		fmt.Println(text)
		return nil
	case "list-not-done":
		text, err := tm.ListNotDoneTasks()
		if err != nil {
			return err
		}
		fmt.Println(text)
		return nil
	case "list-in-progress":
		text, err := tm.ListInProgressTasks()
		if err != nil {
			return err
		}
		fmt.Println(text)
		return nil
	default:
		return fmt.Errorf("unknown command: %s", command)
	}
}

func (tm *TaskManager) AddTask(name, text string) error {
	task := Task{
		ID:         tm.nextID,
		Name:       name,
		Text:       text,
		Done:       false,
		InProgress: false,
	}

	tm.tasks = append(tm.tasks, task)
	tm.nextID++

	return nil
}

func (tm *TaskManager) TaskList() (string, error) {
	if len(tm.tasks) == 0 {
		return "No tasks found", nil
	}

	var OutputStr string

	OutputStr += fmt.Sprintf("%-4s %-20s %-30s %-6s %-12s\n", "ID", "Name", "Text", "Done", "In Progress")
	OutputStr += fmt.Sprintln("---------------------------------------------------------------------------")

	for _, task := range tm.tasks {
		doneStatus := "No"
		if task.Done {
			doneStatus = "Yes"
		}
		progressStatus := "No"
		if task.InProgress {
			progressStatus = "Yes"
		}
		OutputStr += fmt.Sprintf("%-4d %-20s %-30s %-6s %-12s\n",
			task.ID, task.Name, task.Text, doneStatus, progressStatus)
	}

	return OutputStr, nil
}

func (tm *TaskManager) UpdateTaskStatus(id int) error {
	for i := range tm.tasks {
		task := &tm.tasks[i]
		if task.ID == id {
			task.Done = !task.Done
			if task.Done {
				task.InProgress = false
			}
			return nil
		}
	}
	return fmt.Errorf("%w: task with ID %d not found", ErrTaskNotFound, id)
}

func (tm *TaskManager) MarkTaskInProgress(id int) error {
	for i := range tm.tasks {
		task := &tm.tasks[i]
		if task.ID == id {
			task.InProgress = true
			task.Done = false
			return nil
		}
	}
	return fmt.Errorf("%w: task with ID %d not found", ErrTaskNotFound, id)
}

func (tm *TaskManager) MarkTaskDone(id int) error {
	for i := range tm.tasks {
		task := &tm.tasks[i]
		if task.ID == id {
			task.Done = true
			task.InProgress = false
			return nil
		}
	}
	return fmt.Errorf("%w: task with ID %d not found", ErrTaskNotFound, id)
}

func (tm *TaskManager) ListDoneTasks() (string, error) {
	var doneTasks []Task
	for _, task := range tm.tasks {
		if task.Done {
			doneTasks = append(doneTasks, task)
		}
	}

	if len(doneTasks) == 0 {
		return "No done tasks found", nil
	}

	return tm.formatTaskList(doneTasks, "Done Tasks")
}

func (tm *TaskManager) ListNotDoneTasks() (string, error) {
	var notDoneTasks []Task
	for _, task := range tm.tasks {
		if !task.Done {
			notDoneTasks = append(notDoneTasks, task)
		}
	}

	if len(notDoneTasks) == 0 {
		return "No not-done tasks found", nil
	}

	return tm.formatTaskList(notDoneTasks, "Not Done Tasks")
}

func (tm *TaskManager) ListInProgressTasks() (string, error) {
	var inProgressTasks []Task
	for _, task := range tm.tasks {
		if task.InProgress {
			inProgressTasks = append(inProgressTasks, task)
		}
	}

	if len(inProgressTasks) == 0 {
		return "No in-progress tasks found", nil
	}

	return tm.formatTaskList(inProgressTasks, "In Progress Tasks")
}

func (tm *TaskManager) formatTaskList(tasks []Task, title string) (string, error) {
	var outputStr string

	outputStr += fmt.Sprintf("%s:\n", title)
	outputStr += fmt.Sprintf("%-4s %-20s %-30s %-6s %-12s\n", "ID", "Name", "Text", "Done", "In Progress")
	outputStr += fmt.Sprintln("---------------------------------------------------------------------------")

	for _, task := range tasks {
		doneStatus := "No"
		if task.Done {
			doneStatus = "Yes"
		}
		progressStatus := "No"
		if task.InProgress {
			progressStatus = "Yes"
		}
		outputStr += fmt.Sprintf("%-4d %-20s %-30s %-6s %-12s\n",
			task.ID, task.Name, task.Text, doneStatus, progressStatus)
	}

	return outputStr, nil
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
