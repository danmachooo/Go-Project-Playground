package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

const TASK_DIR string = "tasks.json"

type Task struct {
	Id string
	Name string
	Description string
	IsCompleted bool
}

type Data struct {
    Tasks []Task `json:"tasks"`
}

func clearScreen() {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default: // Linux or macOS
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	
	cmd.Run()
}

func promptChoice(prompt string, reader *bufio.Reader) string {
	fmt.Println(prompt)
	fmt.Print(">>")

	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(strings.ToLower(input))
}

func promptRaw(prompt string, reader *bufio.Reader) string {
	fmt.Println(prompt)
	fmt.Print(">>")

	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func addTask(task Task, tasks *[]Task) {
	*tasks = append(*tasks, task)

}

func markTask(taskNumber int, tasks []Task ) (*Task, bool) {
	if taskNumber < 1 ||taskNumber > len(tasks) {
		return nil, false
	}

	idx := taskNumber - 1
	tasks[idx].IsCompleted = !tasks[idx].IsCompleted
	return  &tasks[idx], true
}

func removeTask(taskNumber int, tasks []Task) ([]Task, bool) {
	if taskNumber < 1 || taskNumber > len(tasks) {
		return nil, false
	}

	idx := taskNumber - 1
	newTasks := append(tasks[:idx], tasks[idx+1:]...)
	return newTasks, true
}

func confirmRemoval(taskName string, reader *bufio.Reader) bool {
	confirmation := promptChoice(fmt.Sprintf("Are you sure you want to remove task \"%s\"? (y/n)", taskName), reader)
	return confirmation == "y" || confirmation == "yes"
}

func filterTasks(tasks []Task, filter string) []Task {
	switch filter {
		case "completed":
			var completedTasks []Task
			for _, task := range tasks {
				if task.IsCompleted {
					completedTasks = append(completedTasks, task)
				}
			}
			return completedTasks
		case "pending":
			var pendingTasks []Task
			for _, task := range tasks {
				if !task.IsCompleted {
					pendingTasks = append(pendingTasks, task)
				}
			}
			return pendingTasks
	}
	return tasks
}

func sortTasks(tasks []Task, sortBy string) []Task {
	tasksCopy := make([]Task, len(tasks))
	copy(tasksCopy, tasks)
	
	switch sortBy {
		case "completion":
			sort.Slice(tasksCopy, func(i, j int) bool {
				// Pending tasks first, then completed
				if tasksCopy[i].IsCompleted != tasksCopy[j].IsCompleted {
					return !tasksCopy[i].IsCompleted
				}
				return tasksCopy[i].Name < tasksCopy[j].Name
			})
		case "name":
			sort.Slice(tasksCopy, func(i, j int) bool {
				return tasksCopy[i].Name < tasksCopy[j].Name
			})
	}
	
	return tasksCopy
}

func printTaskTable(tasks []Task) {
    fmt.Printf("%-10s %-36s %-20s %-40s %-10s\n", "Task No.", "ID", "Name", "Description", "Completed")
    fmt.Println(strings.Repeat("-", 120))
    for idx, t := range tasks {
        fmt.Printf(
            "%-10d %-36s %-20s %-40s %-10v\n",
            idx+1, t.Id, t.Name, t.Description, t.IsCompleted,
        )
    }
}


func getTasks(tasks []Task) {
	if len(tasks) < 1 {
		fmt.Println("Tasks are empty, create a new one...")
		return
	}

	printTaskTable(tasks)
}

func getTasksWithOptions(tasks []Task, reader *bufio.Reader) {
	if len(tasks) < 1 {
		fmt.Println("Tasks are empty, create a new one...")
		return
	}

	// Ask for filter
	filter := promptChoice("Filter tasks (all/completed/pending):", reader)
	filteredTasks := filterTasks(tasks, filter)
	
	if len(filteredTasks) < 1 {
		fmt.Printf("No %s tasks found.\n", filter)
		return
	}

	// Ask for sort
	sortBy := promptChoice("Sort by (none/name/completion):", reader)
	if sortBy != "none" {
		filteredTasks = sortTasks(filteredTasks, sortBy)
	}

	clearScreen()
	fmt.Printf("\nShowing %s tasks", filter)
	if sortBy != "none" {
		fmt.Printf(" sorted by %s", sortBy)
	}
	fmt.Printf(":\n")

	printTaskTable(filteredTasks)
}

func saveTasks(tasks []Task) {
	jsonData, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		fmt.Println("Error encoding tasks:", err)
		return
	}

	err = os.WriteFile(TASK_DIR, jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
	fmt.Println("Tasks successfully saved to tasks.json")
}

func loadTasks() []Task {
	data, err := os.ReadFile(TASK_DIR)
	if err != nil {
		return []Task{} // empty if file doesn't exist
	}

	var tasks []Task
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		fmt.Println("Error reading tasks:", err)
		return []Task{}
	}
	return tasks
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var tasks []Task = loadTasks()

    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered from error:", r)
            saveTasks(tasks)
        }
    }()

	for {
		action := promptChoice("Select action: add / fetch / list / mark / remove / quit", reader)
		switch action {
		case "add":
			clearScreen()
			var task Task
			fmt.Printf("\n\nCreate new task...\n\n")
			task.Name = promptRaw("Enter task name", reader )

			if task.Name == "" {
				fmt.Println("Task name cannot be empty.")
				break
			}
			task.Description = promptRaw("Enter task description", reader)

			if task.Description == "" {
				fmt.Println("Task description cannot be empty.")
				break
			}
			task.Id = uuid.New().String()
			addTask(task, &tasks)

			fmt.Printf("Task \"%v\" has been added.\n\n", task.Name)
			
			saveTasks(tasks)


		case "fetch":
			clearScreen()
			getTasks(tasks)
			
		case "list":
			clearScreen()
			getTasksWithOptions(tasks, reader)

		case "mark":
			clearScreen()
			getTasks(tasks) // Show tasks first
			fmt.Println()
			taskNumber := promptChoice("Task number: ", reader)
			i, err := strconv.Atoi(taskNumber)
			if err != nil {
				fmt.Println("Invalid number. Please enter a valid task number.")
				break
			}
			markedTask, isMarked := markTask(i, tasks)

			if isMarked {
				fmt.Printf("Task \"%s\" marked as %v\n\n", markedTask.Name, markedTask.IsCompleted)
			} else {
				fmt.Println("Task number does not exist.")
			}
			saveTasks(tasks)


		case "remove":
			clearScreen()
			getTasks(tasks) // Show tasks first
			fmt.Println()
			taskNumber := promptChoice("Task number: ", reader)
			i, err := strconv.Atoi(taskNumber)
			if err != nil {
				fmt.Println("Invalid number. Please enter a valid task number.")
				break
			}
			
			// Get task name for confirmation before removing
			if i < 1 || i > len(tasks) {
				fmt.Println("Task number does not exist.")
				break
			}
			
			taskName := tasks[i-1].Name
			if !confirmRemoval(taskName, reader) {
				fmt.Println("Task removal cancelled.")
				break
			}
			
			newTasks, isRemoved := removeTask(i, tasks)
			
			if isRemoved {
				tasks = newTasks
				fmt.Printf("Task \"%s\" has been removed\n", taskName)
			} else {
				fmt.Println("Task number does not exist.")
			}
			saveTasks(tasks)

		case "quit":
			confirm := promptChoice("Are you sure you want to quit? (y/n)", reader)
			if confirm != "y" && confirm != "yes" {
				break
			}
			saveTasks(tasks)

			fmt.Println("Exiting task manager...")
			return
		default:
			fmt.Println("Invalid action. Try again.")
		}
		reader.ReadString('\n')
		clearScreen()
	}
}