package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"task-cli/task"
)

func main() {
	fmt.Println("Welcome to Task Tracker CLI!")
	fmt.Println("Available commands: add, list, update, delete, exit!")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("\nEnter a command.")
		scanner.Scan()
		command := strings.TrimSpace(scanner.Text())
		switch command {
		case "add":
			fmt.Println("Enter task description")
			scanner.Scan()
			description := strings.TrimSpace(scanner.Text())
			fmt.Print("Enter due date (YYYY-MM-DD) or press Enter to skip: ")
			scanner.Scan()
			dueDate := strings.TrimSpace(scanner.Text())
			fmt.Print("Enter priority (low, medium, high) or press Enter for default (medium): ")
			scanner.Scan()
			priority := strings.TrimSpace(scanner.Text())

			if priority == "" {
				priority = "medium"
			}
			err := task.AddTasks(description, dueDate, priority)

			if err != nil {
				fmt.Println("Error adding task", err)
			} else {
				fmt.Println("Task added succesfully!")
			}
		case "list":
			fmt.Println("Enter task status to filter (todo, in-progress, done) or press Enter for all: ")

			scanner.Scan()

			status := strings.TrimSpace(scanner.Text())

			tasks, err := task.FilterTasks(status)

			if err != nil {
				fmt.Println("Error loading tasks", err)
			} else if len(tasks) == 0 {
				fmt.Println("No tasks found matching this filter.")
			} else {
				fmt.Println("\n---------------------------------------------------------------------")
				fmt.Println("| ID  | Description      | Status     | Due Date     | Priority     |")
				fmt.Println("---------------------------------------------------------------------")
				for _, t := range tasks {
					fmt.Printf("| %-3d | %-16s | %-10s | %-12s | %-12s |\n", t.ID, t.Description, t.Status, t.DueDate, t.Priority)
				}
				fmt.Println("---------------------------------------------------------------------")
			}
		case "search":
			fmt.Println("Enter a keyword to search")

			scanner.Scan()

			keyword := strings.ToLower(strings.TrimSpace(scanner.Text()))

			tasks, err := task.SearchTask(keyword)

			if err != nil {
				fmt.Println("Error loading tasks", err)
			} else if len(tasks) == 0 {
				fmt.Println("No tasks found matching this keyword.")
			} else {
				fmt.Println("\n---------------------------------------------------------------------")
				fmt.Println("| ID  | Description      | Status     | Due Date     | Priority     |")
				fmt.Println("---------------------------------------------------------------------")
				for _, t := range tasks {
					fmt.Printf("| %-3d | %-16s | %-10s | %-12s | %-12s |\n", t.ID, t.Description, t.Status, t.DueDate, t.Priority)
				}
				fmt.Println("---------------------------------------------------------------------")
			}
		case "update":
			fmt.Println("Enter Task ID to update:")
			scanner.Scan()
			idInput := scanner.Text()

			id, err := strconv.Atoi(strings.TrimSpace(idInput))

			if err != nil {
				fmt.Println("Invalid Task ID!")
				break
			}

			fmt.Println("Enter new description:")
			scanner.Scan()
			newDescription := strings.TrimSpace(scanner.Text())

			err = task.UpdateTask(id, newDescription)

			if err != nil {
				fmt.Println("Error updating Task", err)
			} else {
				fmt.Println("Task updated successfully!")
			}
		case "undo":
			err := task.UndoLastAction()
			if err != nil {
				fmt.Println("Error undoing last action", err)
			} else {
				fmt.Println("Last action undone successfully!")
			}
		case "delete":
			fmt.Println("Enter Task ID to delete:")
			scanner.Scan()
			idInput := scanner.Text()

			id, err := strconv.Atoi(strings.TrimSpace(idInput))

			if err != nil {
				fmt.Println("Invalid Task ID!")
				break
			}

			err = task.DeleteTask(id)

			if err != nil {
				fmt.Println("Error deleting Task", err)
			} else {
				fmt.Println("Task deleted successfully!")
			}
		case "mark-in-progress", "mark-done":
			fmt.Println("Enter Task ID:")
			scanner.Scan()
			idInput := scanner.Text()

			id, err := strconv.Atoi(strings.TrimSpace(idInput))

			if err != nil {
				fmt.Println("Invalid Task ID!")
				break
			}
			newStatus := "done"

			if command == "mark-in-progress" {
				newStatus = "in-progress"
			}

			err = task.MarkTask(id, newStatus)

			if err != nil {
				fmt.Println("Error updating Task status.\n", err)
			} else {
				fmt.Printf("Task marked as %s successfully!", newStatus)
			}
		case "exit":
			fmt.Println("Goodbye!")
			return
		case "help":
			fmt.Println("\nAvailable Commands:")
			fmt.Println("  add              Add a new task")
			fmt.Println("  list             Show all tasks")
			fmt.Println("  list done        Show only completed tasks")
			fmt.Println("  list todo        Show only pending tasks")
			fmt.Println("  list in-progress Show only tasks in progress")
			fmt.Println("  search 			Searchs tasks by input keyword")
			fmt.Println("  undo 			Undos last action")
			fmt.Println("  update           Update a task description")
			fmt.Println("  delete           Remove a task")
			fmt.Println("  mark-in-progress Mark a task as 'in-progress'")
			fmt.Println("  mark-done        Mark a task as 'done'")
			fmt.Println("  exit             Exit the application")
		default:
			fmt.Println("Unknown command. Available commands: add, list, list done, list todo, list in-progress, search, undo, update, delete, mark-in-progress, mark-done, exit")
		}
	}

}
