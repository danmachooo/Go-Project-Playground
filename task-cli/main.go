package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"task-cli/app"
	"task-cli/task"
	"task-cli/ui"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var tasks []task.Task = task.LoadTasks()

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from error:", r)
			task.SaveTasks(tasks)
		}
	}()

	for {
		action := ui.PromptChoice("Select action: add / fetch / list / mark / remove / quit", reader)
		switch action {
		case "add":
			ui.ClearScreen()
			fmt.Printf("\n\nCreate new task...\n\n")
			name := ui.PromptRaw("Enter task name", reader)

			if name == "" {
				fmt.Println("Task name cannot be empty.")
				break
			}
			description := ui.PromptRaw("Enter task description", reader)

			if description == "" {
				fmt.Println("Task description cannot be empty.")
			 break
			}
			newTask := task.NewTask(name, description)
			app.AddTask(newTask, &tasks)

			fmt.Printf("Task \"%v\" has been added.\n\n", newTask.Name)
			task.SaveTasks(tasks)

		case "fetch":
			ui.ClearScreen()
			ui.GetTasks(tasks)

		case "list":
			ui.ClearScreen()
			ui.GetTasksWithOptions(tasks, reader)

		case "mark":
			ui.ClearScreen()
			ui.GetTasks(tasks)
			fmt.Println()
			taskNumber := ui.PromptChoice("Task number: ", reader)
			i, err := strconv.Atoi(taskNumber)
			if err != nil {
				fmt.Println("Invalid number. Please enter a valid task number.")
				break
			}
			markedTask, isMarked := app.MarkTask(i, tasks)

			if isMarked {
				fmt.Printf("Task \"%s\" marked as %v\n\n", markedTask.Name, markedTask.IsCompleted)
			} else {
				fmt.Println("Task number does not exist.")
			}
			task.SaveTasks(tasks)

		case "remove":
			ui.ClearScreen()
			ui.GetTasks(tasks)
			fmt.Println()
			taskNumber := ui.PromptChoice("Task number: ", reader)
			i, err := strconv.Atoi(taskNumber)
			if err != nil {
				fmt.Println("Invalid number. Please enter a valid task number.")
				break
			}

			if i < 1 || i > len(tasks) {
				fmt.Println("Task number does not exist.")
				break
			}

			taskName := tasks[i-1].Name
			if !ui.ConfirmRemoval(taskName, reader) {
				fmt.Println("Task removal cancelled.")
				break
			}

			newTasks, isRemoved := app.RemoveTask(i, tasks)

			if isRemoved {
				tasks = newTasks
				fmt.Printf("Task \"%s\" has been removed\n", taskName)
			} else {
				fmt.Println("Task number does not exist.")
			}
			task.SaveTasks(tasks)

		case "quit":
			confirm := ui.PromptChoice("Are you sure you want to quit? (y/n)", reader)
			if confirm != "y" && confirm != "yes" {
				break
			}
			task.SaveTasks(tasks)
			fmt.Println("Exiting task manager...")
			return
		default:
			fmt.Println("Invalid action. Try again.")
		}
		reader.ReadString('\n')
		ui.ClearScreen()
	}
}