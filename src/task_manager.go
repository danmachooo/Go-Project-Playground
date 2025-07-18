package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
)

type Task struct {
	Id string
	Name string
	Description string
	IsCompleted bool
}

func promptUser(prompt string, reader *bufio.Reader) string {
	fmt.Println(prompt)
	fmt.Print(">>")

	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
func addTask(task Task, tasks []Task) []Task {
	tasks = append(tasks, task)
	return tasks
}

func getTasks(tasks []Task) {
	for _, t := range tasks {
		fmt.Printf("ID: %s\nName: %s\nDescription: %s\nCompleted: %v\n\n", t.Id,  t.Name, t.Description, t.IsCompleted)
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var tasks []Task

	for {
		action := promptUser("Select action: add / getAll / quit", reader)
		switch action {
		case "add":
			var task Task
			task.Name = promptUser("Name", reader )
			task.Description = promptUser("Description", reader)
			task.Id = string(uuid.New().String())
			tasks = addTask(task, tasks)

			if tasks != nil {
				fmt.Printf("Task \"%v\" has been added.\n\n", task.Name)
			}
		case "getAll":
			getTasks(tasks)
		case "quit":
			fmt.Println("Exiting task manager...")
			return
		default:
			fmt.Println("Invalid action. Try again.")
		}
	}
}