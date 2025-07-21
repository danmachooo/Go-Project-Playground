package task

import (
	"encoding/json"
	"fmt"
	"os"
)

const TASK_DIR string = "tasks.json"

func SaveTasks(tasks []Task) {
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

func LoadTasks() []Task {
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