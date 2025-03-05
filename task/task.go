package task

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"time"
)

type Task struct {
	ID          int
	Description string
	Status      string
	CreatedAt   string
	UpdatedAt   string
	DueDate     string
	Priority    string
}

const fileName = "task/tasks.json"
const backup = "task/tasks_backup.json"

func LoadTasks() ([]Task, error) {
	// Check if the file exists
	if _, err := os.Stat("task"); os.IsNotExist(err) {
		err := os.Mkdir("task", os.ModePerm) // Create the directory
		if err != nil {
			return nil, err
		}
	}
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		// Create an empty tasks.json file
		err := os.WriteFile(fileName, []byte("[]"), 0644)
		if err != nil {
			return nil, err
		}
	}

	// Read the file
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	// Parse JSON into a slice of Task
	var tasks []Task
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func SaveTasks(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", " ")
	if err != nil {
		return err
	}
	err = CreateBackup()
	if err != nil {
		return err
	}

	return os.WriteFile(fileName, data, 0644)
}

func CreateBackup() error {
	data, err := os.ReadFile(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return os.WriteFile(backup, data, 0644)
}

func UndoLastAction() error {
	if _, err := os.Stat(backup); os.IsNotExist(err) {
		return errors.New("No backup found, cannot undo last action.")
	}

	data, err := os.ReadFile(backup)

	if err != nil {
		return err
	}
	return os.WriteFile(fileName, data, 0644)
}

func AddTasks(description, dueDate, priority string) error {
	tasks, err := LoadTasks()
	if err != nil {
		return err
	}

	newID := 1
	if len(tasks) > 0 {
		newID = tasks[len(tasks)-1].ID + 1
	}
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	newTask := Task{
		ID:          newID,
		Description: description,
		Status:      "todo",
		CreatedAt:   currentTime,
		UpdatedAt:   currentTime,
		DueDate:     dueDate,
		Priority:    priority,
	}
	tasks = append(tasks, newTask)

	err = SaveTasks(tasks)
	if err != nil {
		return err
	}
	return nil
}

func UpdateTask(id int, newDescription string) error {
	tasks, err := LoadTasks()

	if err != nil {
		return err
	}

	found := false
	for i, t := range tasks {
		if t.ID == id {
			tasks[i].Description = newDescription
			tasks[i].UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
			found = true
			break
		}

	}

	if !found {
		return errors.New("Task ID not found. Use 'list' to check available tasks.")
	}

	return SaveTasks(tasks)
}

func DeleteTask(id int) error {
	tasks, err := LoadTasks()

	if err != nil {
		return err
	}

	var UpdatedTasks []Task
	found := false

	for _, t := range tasks {
		if t.ID == id {
			found = true
			continue
		}
		UpdatedTasks = append(UpdatedTasks, t)
	}

	if !found {
		return errors.New("Task ID not found. Use 'list' to check available tasks.")
	}

	return SaveTasks(UpdatedTasks)
}

func MarkTask(id int, newStatus string) error {
	tasks, err := LoadTasks()

	if err != nil {
		return err
	}

	found := false

	for i, t := range tasks {
		if t.ID == id {
			if t.Status == newStatus {
				return errors.New("Task is already markes as'" + newStatus + "'.")
			}
			tasks[i].Status = newStatus
			tasks[i].UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
			found = true
			break
		}
	}

	if !found {
		return errors.New("Task ID not found. Use 'list' to check available tasks.")
	}

	return SaveTasks(tasks)
}

func FilterTasks(status string) ([]Task, error) {
	tasks, err := LoadTasks()

	if err != nil {
		return nil, err
	}

	if status == "" {
		return tasks, nil
	}

	var filteredTaks []Task

	for _, t := range tasks {
		if t.Status == status {
			filteredTaks = append(filteredTaks, t)
		}
	}

	return filteredTaks, nil
}

func SearchTask(keyword string) ([]Task, error) {
	tasks, err := LoadTasks()

	if err != nil {
		return nil, err
	}

	var results []Task

	for _, t := range tasks {
		if strings.Contains(strings.ToLower(t.Description), keyword) {
			results = append(results, t)
		}
	}

	return results, nil
}
