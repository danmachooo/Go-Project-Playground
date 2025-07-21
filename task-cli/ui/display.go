package ui

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strings"
	"bufio"
	"task-cli/task"
)

func ClearScreen() {
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

func PrintTaskTable(tasks []task.Task) {
	fmt.Printf("%-10s %-36s %-20s %-40s %-10s\n", "Task No.", "ID", "Name", "Description", "Completed")
	fmt.Println(strings.Repeat("-", 120))
	for idx, t := range tasks {
		fmt.Printf(
			"%-10d %-36s %-20s %-40s %-10v\n",
			idx+1, t.Id, t.Name, t.Description, t.IsCompleted,
		)
	}
}

func GetTasks(tasks []task.Task) {
	if len(tasks) < 1 {
		fmt.Println("Tasks are empty, create a new one...")
		return
	}
	PrintTaskTable(tasks)
}

func FilterTasks(tasks []task.Task, filter string) []task.Task {
	switch filter {
	case "completed":
		var completedTasks []task.Task
		for _, task := range tasks {
			if task.IsCompleted {
				completedTasks = append(completedTasks, task)
			}
		}
		return completedTasks
	case "pending":
		var pendingTasks []task.Task
		for _, task := range tasks {
			if !task.IsCompleted {
				pendingTasks = append(pendingTasks, task)
			}
		}
		return pendingTasks
	}
	return tasks
}

func SortTasks(tasks []task.Task, sortBy string) []task.Task {
	tasksCopy := make([]task.Task, len(tasks))
	copy(tasksCopy, tasks)

	switch sortBy {
	case "completion":
		sort.Slice(tasksCopy, func(i, j int) bool {
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

func GetTasksWithOptions(tasks []task.Task, reader *bufio.Reader) {
	if len(tasks) < 1 {
		fmt.Println("Tasks are empty, create a new one...")
		return
	}

	filter := PromptChoice("Filter tasks (all/completed/pending):", reader)
	filteredTasks := FilterTasks(tasks, filter)

	if len(filteredTasks) < 1 {
		fmt.Printf("No %s tasks found.\n", filter)
		return
	}

	sortBy := PromptChoice("Sort by (none/name/completion):", reader)
	if sortBy != "none" {
		filteredTasks = SortTasks(filteredTasks, sortBy)
	}

	ClearScreen()
	fmt.Printf("\nShowing %s tasks", filter)
	if sortBy != "none" {
		fmt.Printf(" sorted by %s", sortBy)
	}
	fmt.Printf(":\n")

	PrintTaskTable(filteredTasks)
}